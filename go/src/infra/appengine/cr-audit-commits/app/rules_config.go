// Copyright 2017 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package crauditcommits

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/context"

	"go.chromium.org/luci/common/proto/git"
)

// DynamicRefFunc is a functype for functions that match a RepoConfig with a
// dynamically determined ref.
//
// It is expected to receive the generic RepoConfig as hardcoded in RulesMap,
// passed by value to prevent the implementation from modifying it.
// It is expected to return a slice of references to RepoConfigs, where each
// matches a ref to audit, and its values BranchName and Metadata have been
// modified accordingly.
//
// Note that for changes to any other field of RepoConfig made by functions
// implementing this interface to persist and apply to audits, the Scheduler
// needs to be modified to save them to the RepoState, and the SetConcreteRef
// function below needs to be modified to set them in the copy of RepoConfig to
// be passed to the scan/audit/notify functions.
type DynamicRefFunc func(context.Context, RepoConfig) ([]*RepoConfig, error)

// RepoConfig represents the hard-coded config for a monitored repo and a
// pointer to the entity representing its datastore-persisted state.
type RepoConfig struct { // These are expected to be hard-coded.
	BaseRepoURL     string
	GerritURL       string
	BranchName      string
	Metadata        string
	StartingCommit  string
	MonorailAPIURL  string // Only intended for release branches
	MonorailProject string
	// Do not use "AuditFailure" as a key in this map, it may cause a clash
	// with the notification state for failed audits.
	Rules              map[string]RuleSet
	NotifierEmail      string
	DynamicRefFunction DynamicRefFunc
}

// BranchInfo represents the main branch information of a specific Chrome release
type BranchInfo struct {
	PdfiumBranch   string `json:"pdfium_branch"`
	SkiaBranch     string `json:"skia_branch"`
	WebrtcBranch   string `json:"webrtc_branch"`
	V8Branch       string `json:"v8_branch"`
	ChromiumBranch string `json:"chromium_branch"`
	Milestone      int    `json:"milestone"`
}

// RepoURL composes the url of the repository by appending the branch.
func (rc *RepoConfig) RepoURL() string {
	return rc.BaseRepoURL + "/+/" + rc.BranchName
}

// LinkToCommit composes a url to a specific commit
func (rc *RepoConfig) LinkToCommit(commit string) string {
	return rc.BaseRepoURL + "/+/" + commit
}

// SetConcreteRef returns a copy of the repoconfig modified to account for
// dynamic refs.
func (rc *RepoConfig) SetConcreteRef(ctx context.Context, rs *RepoState) *RepoConfig {
	// Make a copy.
	result := *rc
	if rs.BranchName != "" {
		result.BranchName = rs.BranchName
	}
	if rs.Metadata != "" {
		result.Metadata = rs.Metadata
	}
	return &result
}

// RuleMap maps each monitored repository to a list of account/rules structs.
var RuleMap = map[string]*RepoConfig{
	"chromium-src-master": {
		BaseRepoURL: "https://chromium.googlesource.com/chromium/src.git",
		GerritURL:   "https://chromium-review.googlesource.com",
		BranchName:  "master",
		// No special meaning, ToT as of the time this line was added.
		StartingCommit:  "bafa682dc0ce1dde367ba44f31f8ec1ad07e569e",
		MonorailAPIURL:  "https://monorail-prod.appspot.com/_ah/api/monorail/v1",
		MonorailProject: "chromium",
		NotifierEmail:   "notifier@cr-audit-commits.appspotmail.com",
		Rules: map[string]RuleSet{
			"autoroll-rules-afdo":         AutoRollRulesAFDOVersion("afdo-chromium-autoroll@skia-buildbots.google.com.iam.gserviceaccount.com"),
			"autoroll-rules-angle":        AutoRollRulesDEPS("angle-chromium-autoroll@skia-buildbots.google.com.iam.gserviceaccount.com"),
			"autoroll-rules-catapult":     AutoRollRulesDEPS("catapult-chromium-autoroll@skia-buildbots.google.com.iam.gserviceaccount.com"),
			"autoroll-rules-chromite":     AutoRollRulesDEPS("chromite-chromium-autoroll@skia-buildbots.google.com.iam.gserviceaccount.com"),
			"autoroll-rules-depot-tools":  AutoRollRulesDEPS("depot-tools-chromium-autoroll@skia-buildbots.google.com.iam.gserviceaccount.com"),
			"autoroll-rules-fuchsia-sdk":  AutoRollRulesFuchsiaSDKVersion("fuchsia-sdk-chromium-autoroll@skia-buildbots.google.com.iam.gserviceaccount.com"),
			"autoroll-rules-ios-internal": AutoRollRulesDEPS("ios-internal-chromium-autoroll@skia-buildbots.google.com.iam.gserviceaccount.com"),
			"autoroll-rules-nacl":         AutoRollRulesDEPS("nacl-chromium-autoroll@skia-buildbots.google.com.iam.gserviceaccount.com"),
			"autoroll-rules-pdfium":       AutoRollRulesDEPS("pdfium-chromium-autoroll@skia-buildbots.google.com.iam.gserviceaccount.com"),
			"autoroll-rules-perfetto":     AutoRollRulesDEPS("perfetto-chromium-autoroll@skia-buildbots.google.com.iam.gserviceaccount.com"),
			"autoroll-rules-skia":         AutoRollRulesDEPS("skia-chromium-autoroll@skia-buildbots.google.com.iam.gserviceaccount.com"),
			"autoroll-rules-src-internal": AutoRollRulesDEPS("src-internal-chromium-autoroll@skia-buildbots.google.com.iam.gserviceaccount.com"),
			"autoroll-rules-webrtc":       AutoRollRulesDEPS("webrtc-chromium-autoroll@skia-buildbots.google.com.iam.gserviceaccount.com"),
			"findit-rules": AccountRules{
				Account: "findit-for-me@appspot.gserviceaccount.com",
				Funcs: []RuleFunc{
					AutoCommitsPerDay,
					AutoRevertsPerDay,
					CulpritAge,
					CulpritInBuild,
					FailedBuildIsAppropriateFailure,
					RevertOfCulprit,
					OnlyCommitsOwnChange,
				},
				notificationFunction: fileBugForFinditViolation,
			},
			"release-bot-rules": AccountRules{
				Account: "chrome-release-bot@chromium.org",
				Funcs: []RuleFunc{
					OnlyModifiesReleaseFiles,
				},
				notificationFunction: fileBugForReleaseBotViolation,
			},
		},
	},
	"chromium-src-release-branches": {
		BaseRepoURL:     "https://chromium.googlesource.com/chromium/src.git",
		GerritURL:       "https://chromium-review.googlesource.com",
		MonorailAPIURL:  "https://monorail-prod.appspot.com/_ah/api/monorail/v1",
		MonorailProject: "chromium",
		NotifierEmail:   "notifier@cr-audit-commits.appspotmail.com",
		Rules: map[string]RuleSet{
			"merge-approval-rules": AccountRules{
				Account: "*",
				Funcs: []RuleFunc{
					OnlyMergeApprovedChange,
				},
				notificationFunction: fileBugForMergeApprovalViolation,
			},
		},
		DynamicRefFunction: ReleaseConfig,
	},
	"skia-master": {
		BaseRepoURL: "https://skia.googlesource.com/skia.git",
		GerritURL:   "https://skia-review.googlesource.com",
		BranchName:  "refs/heads/master",
		// No special meaning, ToT as of the time this line was added.
		StartingCommit:  "82a33425166aacd0726bdd283c6de749420819a8",
		MonorailAPIURL:  "https://monorail-prod.appspot.com/_ah/api/monorail/v1",
		MonorailProject: "chromium",
		NotifierEmail:   "notifier@cr-audit-commits.appspotmail.com",
		Rules: map[string]RuleSet{
			"autoroll-rules-angle":       AutoRollRulesDEPS("angle-skia-autoroll@skia-buildbots.google.com.iam.gserviceaccount.com"),
			"autoroll-rules-chromium":    AutoRollRulesDEPS("chromium-skia-autoroll@skia-buildbots.google.com.iam.gserviceaccount.com"),
			"autoroll-rules-skcms":       AutoRollRulesSKCMS("skcms-skia-autoroll@skia-buildbots.google.com.iam.gserviceaccount.com"),
			"autoroll-rules-swiftshader": AutoRollRulesDEPS("swiftshader-skia-autoroll@skia-buildbots.google.com.iam.gserviceaccount.com"),
			"bookmaker": AccountRules{
				Account: "skia-bookmaker@skia-swarming-bots.iam.gserviceaccount.com",
				Funcs: []RuleFunc{
					func(ctx context.Context, ap *AuditParams, rc *RelevantCommit, cs *Clients) *RuleResult {
						return OnlyModifiesDirRule(ctx, ap, rc, cs, "OnlyModifiesAPIDocs", "site/user/api")
					},
				},
				notificationFunction: fileBugForAutoRollViolation,
			},
			"recreate-skps": AccountRules{
				Account: "skia-recreate-skps@skia-swarming-bots.iam.gserviceaccount.com",
				Funcs: []RuleFunc{
					func(ctx context.Context, ap *AuditParams, rc *RelevantCommit, cs *Clients) *RuleResult {
						files := []string{
							"infra/bots/assets/skp/VERSION",
							"infra/bots/tasks.json",
						}
						return OnlyModifiesFilesRule(ctx, ap, rc, cs, "OnlyModifiesVersionFile", files)
					},
				},
				notificationFunction: fileBugForAutoRollViolation,
			},
		},
	},
	"skia-lottie-ci": {
		BaseRepoURL: "https://skia.googlesource.com/lottie-ci.git",
		GerritURL:   "https://skia-review.googlesource.com",
		BranchName:  "refs/heads/master",
		// No special meaning, ToT as of the time this line was added.
		StartingCommit:  "6844651ced137fd86d73a11cd0c4d74e71c6fb98",
		MonorailAPIURL:  "https://monorail-prod.appspot.com/_ah/api/monorail/v1",
		MonorailProject: "chromium",
		NotifierEmail:   "notifier@cr-audit-commits.appspotmail.com",
		Rules: map[string]RuleSet{
			"autoroll-rules-lottie-web": AutoRollRulesDEPS("lottie-web-lottie-ci-autoroll@skia-buildbots.google.com.iam.gserviceaccount.com"),
			"autoroll-rules-skia":       AutoRollRulesDEPSAndTasks("skia-lottie-ci-autoroll@skia-buildbots.google.com.iam.gserviceaccount.com"),
		},
	},
	"fuchsia-topaz-master": {
		BaseRepoURL: "https://fuchsia.googlesource.com/topaz.git",
		GerritURL:   "https://fuchsia-review.googlesource.com",
		BranchName:  "refs/heads/master",
		// No special meaning, ToT as of the time this line was added.
		StartingCommit:  "ec7b9088a64bb6a71d8e327a0d04ee9a2f6bb9ec",
		MonorailAPIURL:  "https://monorail-prod.appspot.com/_ah/api/monorail/v1",
		MonorailProject: "chromium",
		NotifierEmail:   "notifier@cr-audit-commits.appspotmail.com",
		Rules: map[string]RuleSet{
			"autoroll-rules-skia": AccountRules{
				Account: "skia-fuchsia-autoroll@skia-buildbots.google.com.iam.gserviceaccount.com",
				Funcs: []RuleFunc{
					func(ctx context.Context, ap *AuditParams, rc *RelevantCommit, cs *Clients) *RuleResult {
						return OnlyModifiesFileRule(ctx, ap, rc, cs, "OnlyModifiesSkiaManifest", "manifest/skia")
					},
				},
				notificationFunction: fileBugForAutoRollViolation,
			},
		},
	},
}

// RuleSet is a group of rules that can be decided to apply or not to a
// specific commit, as a unit.
//
// Note that the methods in this interface are not rules, but tests that can
// decide whether the rules in the set apply to a given commit.
type RuleSet interface {
	MatchesCommit(*git.Commit) bool
	MatchesRelevantCommit(*RelevantCommit) bool
	NotificationFunction() NotificationFunc
}

// AccountRules is a RuleSet that applies to a commit if the commit has a given
// account as either its author or its committer.
type AccountRules struct {
	Account              string
	Funcs                []RuleFunc
	notificationFunction NotificationFunc
}

// NotificationFunction exposes the NotificationFunc assigned to this struct
// as required by the RuleSet interface.
func (ar AccountRules) NotificationFunction() NotificationFunc {
	return ar.notificationFunction
}

// MatchesCommit determines whether the AccountRules set it's bound to, applies
// to the given commit.
func (ar AccountRules) MatchesCommit(c *git.Commit) bool {
	return ar.Account == "*" || c.GetCommitter().GetEmail() == ar.Account || c.GetAuthor().GetEmail() == ar.Account
}

// MatchesRelevantCommit determines whether the AccountRules set it's bound to,
// applies to the given commit entity.
func (ar AccountRules) MatchesRelevantCommit(c *RelevantCommit) bool {
	return ar.Account == "*" || c.CommitterAccount == ar.Account || c.AuthorAccount == ar.Account
}

// AuditParams exposes object shared by all rules (and the worker goroutines
// they are run on).
type AuditParams struct {
	TriggeringAccount string
	RepoCfg           *RepoConfig
	RepoState         *RepoState
}

// RuleFunc is the function type for audit rules.
//
// They are expected to accept a context, an AuditParams, a Clients struct with
// connections to external services configured and ready, and the datastore
// entity to be audited.
//
// Rules are expected to panic if they cannot determine whether a policy has
// been broken or not.
//
// Rules should return a reference to a RuleResult
type RuleFunc func(context.Context, *AuditParams, *RelevantCommit, *Clients) *RuleResult

// ReleaseConfig is the skeleton of a function to get the ref and milestone
// dynamically.
func ReleaseConfig(ctx context.Context, cfg RepoConfig) ([]*RepoConfig, error) {

	var branchRefsURLContents []string
	// ------------------
	// OMAHAPROXY||CHROMIUMDASH MAGIC
	// ------------------

	// https://chromiumdash.appspot.com/fetch_milestones is a legacy API that needs some clean up. Here,
	// the platform could be any Chrome platform orther than Android and the result will still be the same.
	contents, err := getURLAsString(ctx, "https://chromiumdash.appspot.com/fetch_milestones?platform=Android")
	if err != nil {
		return nil, err
	}
	branchInfos := []BranchInfo{}
	err = json.Unmarshal([]byte(contents), &branchInfos)
	if err != nil {
		return nil, err
	}
	// We only need Stable and Beta branch details.
	for i := range branchInfos[:2] {
		branchRefsURL := fmt.Sprintf("https://chromium.googlesource.com/chromium/src.git/+log/refs/heads/master..refs/branch-heads/%s/?format=json&n=1000", branchInfos[i].ChromiumBranch)
		// When scanning a branch for the first time, it's unlikely that there'll be more than 1000 commits in it. In subsequent scans, the starting commit will be ignored,
		// but instead the last scanned commit will be used. So even if the commits in the branch exceed 1000 there will be no effect in the auditing.
		branchContents, err := getURLAsString(ctx, branchRefsURL)
		if err != nil {
			return nil, err
		}
		branchRefsURLContents = append(branchRefsURLContents, branchContents)
	}

	return GetReleaseConfig(ctx, cfg, branchRefsURLContents, branchInfos)
}

// GetReleaseConfig is a helper function to get the ref and milestone dynamically.
func GetReleaseConfig(ctx context.Context, cfg RepoConfig, branchRefsURLContents []string, branchInfos []BranchInfo) ([]*RepoConfig, error) {
	concreteConfigs := []*RepoConfig{}
	var err error
	r := regexp.MustCompile(`"commit": "(.+)?"`)
	for i := range branchRefsURLContents {
		temp := strings.Split(branchRefsURLContents[i], "\n")
		lastCommit := ""
		for j := 0; j < len(temp); j++ {
			if r.MatchString(temp[j]) {
				lastCommit = temp[j]
			}
		}
		if lastCommit == "" {
			return nil, errors.New("commit not found or invalid")
		}
		concreteConfig := cfg
		concreteConfig.StartingCommit = r.FindStringSubmatch(lastCommit)[1]
		concreteConfig.BranchName = fmt.Sprintf("refs/branch-heads/%s", branchInfos[i].ChromiumBranch)
		concreteConfig.Metadata, err = SetToken(ctx, "MilestoneNumber", strconv.Itoa(branchInfos[i].Milestone), concreteConfig.Metadata)
		concreteConfigs = append(concreteConfigs, &concreteConfig)
	}
	return concreteConfigs, err
}
