// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package analyzer

import (
	"net/url"
	"reflect"
	"sort"
	"testing"
	"time"

	"golang.org/x/net/context"

	"infra/appengine/sheriff-o-matic/som/analyzer/regrange"
	analyzertest "infra/appengine/sheriff-o-matic/som/analyzer/test"
	clientTest "infra/appengine/sheriff-o-matic/som/client/test"
	"infra/appengine/test-results/model"
	"infra/libs/testing/ansidiff"
	"infra/monitoring/messages"

	. "github.com/smartystreets/goconvey/convey"
)

func fakeNow(t time.Time) func() time.Time {
	return func() time.Time {
		return t
	}
}

func urlParse(s string, t *testing.T) *url.URL {
	p, err := url.Parse(s)
	if err != nil {
		t.Errorf("failed to parse %s: %s", s, err)
	}
	return p
}

type fakeReasonRaw struct {
	signature string
}

func (f *fakeReasonRaw) Signature() string {
	if f.signature != "" {
		return f.signature
	}

	return "fakeSignature"
}

func (f *fakeReasonRaw) Kind() string {
	return "fakeKind"
}

func (f *fakeReasonRaw) Title([]*messages.BuildStep) string {
	return "fakeTitle"
}

func (f *fakeReasonRaw) Severity() messages.Severity {
	return messages.NewFailure
}

type fakeAnalyzer struct {
}

func (f *fakeAnalyzer) Analyze(ctx context.Context, failures []*messages.BuildStep, tree string) ([]messages.ReasonRaw, []error) {
	return fakeFinder(ctx, failures, tree), nil
}

func fakeFinder(ctx context.Context, failures []*messages.BuildStep, tree string) []messages.ReasonRaw {
	raws := make([]messages.ReasonRaw, len(failures))
	for i := range failures {
		raws[i] = &fakeReasonRaw{}
	}
	return raws
}

func newTestAnalyzer(minBuilds, maxBuilds int) *Analyzer {
	a := New(minBuilds, maxBuilds)
	a.regrangeFinder = regrange.Default
	a.StepAnalyzers = append(a.StepAnalyzers, &fakeAnalyzer{})
	a.Milo = nil
	a.CrBug = nil
	a.FindIt = nil

	return a
}

func TestMasterAlerts(t *testing.T) {
	tests := []struct {
		name   string
		master string
		be     messages.BuildExtract
		t      time.Time
		want   []messages.Alert
	}{
		{
			name:   "empty",
			master: "fake-empty",
			want:   []messages.Alert{},
		},
		{
			name:   "Not stale master",
			master: "fake-not-stale",
			be: messages.BuildExtract{
				CreatedTimestamp: messages.EpochTime(100),
			},
			t:    time.Unix(100, 0),
			want: []messages.Alert{},
		},
		{
			name:   "Stale master",
			master: "fake.master",
			be: messages.BuildExtract{
				CreatedTimestamp: messages.EpochTime(100),
			},
			t: time.Unix(100, 0).Add(20 * time.Minute),
			want: []messages.Alert{
				{
					Key:       "stale master: https://ci.chromium.org/p/fake.master",
					Title:     "Stale https://ci.chromium.org/p/fake.master master data",
					Type:      messages.AlertStaleMaster,
					Severity:  messages.StaleMaster,
					Body:      "0h 20m elapsed since last update.",
					Time:      messages.TimeToEpochTime(time.Unix(100, 0).Add(20 * time.Minute)),
					Links:     []messages.Link{{Title: "Master", Href: urlParse("https://ci.chromium.org/p/fake.master", t).String()}},
					StartTime: messages.EpochTime(100),
				},
			},
		},
		{
			name:   "Future master",
			master: "fake.master",
			be: messages.BuildExtract{
				CreatedTimestamp: messages.EpochTime(110),
			},
			t:    time.Unix(100, 0),
			want: []messages.Alert{},
		},
	}

	mc := &clientTest.MockReader{}
	a := newTestAnalyzer(0, 10)
	a.Milo = mc
	ctx := context.Background()

	for _, test := range tests {
		a.Now = fakeNow(test.t)
		got := a.MasterAlerts(ctx, &messages.MasterLocation{URL: *urlParse("https://ci.chromium.org/p/"+test.master, t)}, &test.be)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s failed. Got %+v, want: %+v\nDiff: %v", test.name, got, test.want,
				ansidiff.Diff(got, test.want))
		}
	}
}

func TestBuilderAlerts(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		be           messages.BuildExtract
		filter       string
		t            time.Time
		wantBuilders []messages.Alert
		wantMasters  []messages.Alert
	}{
		{
			name:         "Empty",
			url:          "https://ci.chromium.org/p/fake.master/json",
			wantBuilders: []messages.Alert{},
			wantMasters:  []messages.Alert{},
		},
		{
			name: "No Builders",
			url:  "https://ci.chromium.org/p/fake.master/json",
			be: messages.BuildExtract{
				CreatedTimestamp: messages.EpochTime(100),
			},
			t:            time.Unix(100, 0),
			wantBuilders: []messages.Alert{},
			wantMasters:  []messages.Alert{},
		},
		{
			name: "No Alerts",
			url:  "https://ci.chromium.org/p/fake.master/json",
			be: messages.BuildExtract{
				CreatedTimestamp: messages.EpochTime(100),
				Builders: map[string]messages.Builder{
					"fake.builder": {},
				},
			},
			t:            time.Unix(100, 0),
			wantBuilders: []messages.Alert{},
			wantMasters:  []messages.Alert{},
		},
	}

	mc := &clientTest.MockReader{}
	a := newTestAnalyzer(0, 10)
	a.Milo = mc
	ctx := context.Background()

	for _, test := range tests {
		a.Now = fakeNow(test.t)
		got := a.BuilderAlerts(ctx, "tree", &messages.MasterLocation{URL: *urlParse(test.url, t)}, &test.be)
		if !reflect.DeepEqual(got, test.wantBuilders) {
			t.Errorf("%s failed. Got %+v, want: %+v", test.name, got, test.wantBuilders)
		}
	}
}

func TestLittleBBuilderAlerts(t *testing.T) {
	Convey("builderAlerts", t, func() {
		tests := []struct {
			name       string
			master     string
			builder    string
			b          messages.Builder
			builds     map[string]*messages.Build
			time       time.Time
			wantAlerts []messages.Alert
			wantErrs   []error
		}{
			{
				name:     "empty",
				wantErrs: []error{errNoRecentBuilds},
			},
			{
				name:    "builders ok",
				master:  "fake.master",
				builder: "fake.builder",
				builds: analyzertest.NewBuilderFaker("fake.master", "fake.builder").
					Build(0).Times(10, 100).
					Step("fake_step").Times(10, 100).BuilderFaker.Builds,
				b: messages.Builder{
					State:        messages.StateBuilding,
					BuilderName:  "fake.builder",
					CachedBuilds: []int64{0},
				},
				wantAlerts: []messages.Alert{},
				wantErrs:   []error{},
			},
			{
				name:    "builder building for too long",
				master:  "fake.master",
				builder: "hung.builder",
				builds: analyzertest.NewBuilderFaker("fake.master", "hung.builder").
					Build(0).Times(10, 100).
					Step("fake_step").Times(10, 100).BuilderFaker.
					Build(1).Times(100, 0).
					Step("fake_step").Times(100, 0).BuilderFaker.Builds,
				b: messages.Builder{
					State:        messages.StateBuilding,
					BuilderName:  "fake.builder",
					CachedBuilds: []int64{0, 1},
				},
				time: time.Unix(0, 0).Add(4 * time.Hour),
				wantAlerts: []messages.Alert{
					{
						Key:       "fake.master.hung.builder.hung",
						Title:     "fake.master.hung.builder is hung in step fake_step.",
						Type:      messages.AlertHungBuilder,
						StartTime: 100,
						Time:      messages.TimeToEpochTime(time.Unix(0, 0).Add(4 * time.Hour)),
						Severity:  messages.HungBuilder,
						Links: []messages.Link{
							{Title: "Builder", Href: urlParse("https://ci.chromium.org/p/fake.master/builders/hung.builder", t).String()},
							{Title: "Last build", Href: urlParse("https://ci.chromium.org/p/fake.master/builders/hung.builder/1", t).String()},
							{Title: "Last build step", Href: urlParse("https://ci.chromium.org/p/fake.master/builders/hung.builder/builds/1/steps/fake_step", t).String()},
						},
					},
				},
				wantErrs: []error{},
			},
			{
				name:    "builder offline for not long enough",
				master:  "fake.master",
				builder: "offline.builder",
				builds: analyzertest.NewBuilderFaker("fake.master", "offline.builder").
					Build(0).Times(10, 100).
					Step("fake_step").Times(0, 60*60).BuilderFaker.
					Build(1).Times(100, 0).
					Step("fake_step").Times(60*60, 0).BuilderFaker.Builds,
				b: messages.Builder{
					State:        messages.StateOffline,
					BuilderName:  "offline.builder",
					CachedBuilds: []int64{0, 1},
				},
				// Last step is at an hour, 1.5 hours is the timeout
				time:       time.Unix(0, 0).Add(2 * time.Hour).Add(30 * time.Minute).Add(-time.Second),
				wantAlerts: []messages.Alert{},
				wantErrs:   []error{},
			},
			{
				name:    "builder offline for too long",
				master:  "fake.master",
				builder: "offline.builder",
				builds: analyzertest.NewBuilderFaker("fake.master", "offline.builder").
					Build(0).Times(10, 100).
					Step("fake_step").Times(0, 2.5*60*60).BuilderFaker.
					Build(1).Times(100, 0).
					Step("fake_step").Times(2.5*60*60, 0).BuilderFaker.Builds,
				b: messages.Builder{
					State:        messages.StateOffline,
					BuilderName:  "offline.builder",
					CachedBuilds: []int64{0, 1},
				},
				// Last step is at an hour, 1.5 hours is the timeout
				time: time.Unix(0, 0).Add(4 * time.Hour).Add(time.Second),
				wantAlerts: []messages.Alert{
					{
						Key:       "fake.master.offline.builder.offline",
						Title:     "fake.master.offline.builder is offline.",
						Type:      messages.AlertOfflineBuilder,
						StartTime: 2.5 * 60 * 60,
						Time:      messages.TimeToEpochTime(time.Unix(0, 0).Add(4 * time.Hour).Add(time.Second)),
						Severity:  messages.OfflineBuilder,
						Links: []messages.Link{
							{Title: "Builder", Href: urlParse("https://ci.chromium.org/p/fake.master/builders/offline.builder", t).String()},
							{Title: "Last build", Href: urlParse("https://ci.chromium.org/p/fake.master/builders/offline.builder/1", t).String()},
							{Title: "Last build step", Href: urlParse("https://ci.chromium.org/p/fake.master/builders/offline.builder/builds/1/steps/fake_step", t).String()},
						},
					},
				},
				wantErrs: []error{},
			},
			{
				name:    "builder idle, not enough pending builds",
				master:  "fake.master",
				builder: "idle.builder",
				builds: analyzertest.NewBuilderFaker("fake.master", "idle.builder").
					Build(0).Times(10, 100).
					Step("fake_step").Times(10, 100).BuilderFaker.
					Build(1).Times(100, 0).
					Step("fake_step").Times(100, 0).BuilderFaker.Builds,
				b: messages.Builder{
					State:         messages.StateIdle,
					BuilderName:   "idle.builder",
					CachedBuilds:  []int64{0, 1},
					PendingBuilds: 49,
				},
				time:       time.Unix(0, 0).Add(4 * time.Hour),
				wantAlerts: []messages.Alert{},
				wantErrs:   []error{},
			},
			{
				name:    "builder idle, too many pending builds",
				master:  "fake.master",
				builder: "idle.builder",
				builds: analyzertest.NewBuilderFaker("fake.master", "idle.builder").
					Build(0).Times(10, 100).
					Step("fake_step").Times(10, 100).BuilderFaker.
					Build(1).Times(100, 0).
					Step("fake_step").Times(100, 0).BuilderFaker.Builds,
				b: messages.Builder{
					State:         messages.StateIdle,
					BuilderName:   "idle.builder",
					CachedBuilds:  []int64{0, 1},
					PendingBuilds: 51,
				},
				time: time.Unix(0, 0).Add(4 * time.Hour),
				wantAlerts: []messages.Alert{
					{
						Key:       "fake.master.idle.builder.idle",
						Title:     "fake.master.idle.builder is idle with 51 pending builds.",
						Type:      messages.AlertIdleBuilder,
						StartTime: 100,
						Time:      messages.TimeToEpochTime(time.Unix(0, 0).Add(4 * time.Hour)),
						Severity:  messages.IdleBuilder,
						Links: []messages.Link{
							{Title: "Builder", Href: urlParse("https://ci.chromium.org/p/fake.master/builders/idle.builder", t).String()},
							{Title: "Last build", Href: urlParse("https://ci.chromium.org/p/fake.master/builders/idle.builder/1", t).String()},
							{Title: "Last build step", Href: urlParse("https://ci.chromium.org/p/fake.master/builders/idle.builder/builds/1/steps/fake_step", t).String()},
						},
					},
				},
				wantErrs: []error{},
			},
		}

		ctx := context.Background()
		a := newTestAnalyzer(0, 10)

		for _, test := range tests {
			test := test
			Convey(test.name, func() {
				a.Now = fakeNow(test.time)
				a.Milo = clientTest.MockReader{
					Builds: test.builds,
				}
				gotAlerts, gotErrs := a.builderAlerts(ctx, "tree", &messages.MasterLocation{URL: *urlParse("https://ci.chromium.org/p/"+test.master, t)}, test.builder, &test.b)
				So(gotAlerts, ShouldResemble, test.wantAlerts)
				So(gotErrs, ShouldResemble, test.wantErrs)
			})
		}
	})
}

func TestBuilderStepAlerts(t *testing.T) {

	Convey("test BuilderStepAlerts", t, func() {
		regrange.URLToNameMapping = map[string]string{
			"http://test": "test",
		}
		tests := []struct {
			name          string
			master        string
			builder       string
			tree          string
			recentBuilds  []int64
			testData      *analyzertest.BuilderFaker
			finditData    []*messages.FinditResult
			finditDataV2  []*messages.FinditResultV2
			time          time.Time
			buildsAtFault []int
			stepsAtFault  []string
			wantAlerts    []messages.Alert
			wantErrs      []error
		}{
			{
				name: "empty",
			},
			{
				name:         "builders ok",
				master:       "fake.master",
				builder:      "fake.builder",
				tree:         "tree",
				recentBuilds: []int64{0},
				testData: analyzertest.NewBuilderFaker("fake.master", "fake.builder").
					Build(0).Times(0, 1).Step("fake_step").Results(0).BuilderFaker,
			},
			{
				name:         "one build failure",
				master:       "fake.master",
				builder:      "fake.builder",
				tree:         "tree",
				recentBuilds: []int64{0},
				testData: analyzertest.NewBuilderFaker("fake.master", "fake.builder").
					Build(0).Times(0, 1).IncludeChanges("http://test", "refs/heads/master@{#291569}").
					Step("fake_step").Results(2).BuilderFaker,
				buildsAtFault: []int{0},
				stepsAtFault:  []string{"fake_step"},
				wantAlerts: []messages.Alert{
					{
						Key:      "fake.master.fake.builder.fake_step.",
						Title:    "fakeTitle",
						Type:     messages.AlertBuildFailure,
						Body:     "",
						Severity: messages.NewFailure,
						Extension: messages.BuildFailure{
							Builders: []messages.AlertedBuilder{
								{
									Name:   "fake.builder",
									Master: "fake.master",
									URL:    urlParse("https://ci.chromium.org/p/fake.master/builders/fake.builder", t).String(),
									Count:  1,
								},
							},
							RegressionRanges: []*messages.RegressionRange{
								{
									Repo:      "test",
									URL:       "http://test",
									Revisions: []string{"291569"},
									Positions: []string{"refs/heads/master@{#291569}"},
								},
							},
							Reason: &messages.Reason{
								Raw: &fakeReasonRaw{},
							},
						},
					},
				},
			},
			{
				name:         "one build failure with findit running",
				master:       "fake.master",
				builder:      "fake.builder",
				tree:         "chromium",
				recentBuilds: []int64{123},
				testData: analyzertest.NewBuilderFaker("fake.master", "fake.builder").
					Build(123).Times(0, 1).IncludeChanges("http://test", "refs/heads/master@{#291569}").
					Step("fake_step").Results(2).BuilderFaker,
				finditData: []*messages.FinditResult{
					{
						MasterURL:   "https://ci.chromium.org/p/fake.master",
						BuilderName: "fake.builder",
						BuildNumber: 123,
						StepName:    "fake_step",
						HasFindings: false,
						IsFinished:  false,
						IsSupported: true,
					},
				},
				buildsAtFault: []int{123},
				stepsAtFault:  []string{"fake_step"},
				wantAlerts: []messages.Alert{
					{
						Key:      "fake.master.fake.builder.fake_step.",
						Title:    "fakeTitle",
						Type:     messages.AlertBuildFailure,
						Body:     "",
						Severity: messages.NewFailure,
						Extension: messages.BuildFailure{
							Builders: []messages.AlertedBuilder{
								{
									Name:          "fake.builder",
									URL:           urlParse("https://ci.chromium.org/p/fake.master/builders/fake.builder", t).String(),
									Master:        "fake.master",
									FirstFailure:  123,
									LatestFailure: 123,
									Count:         1,
								},
							},
							Reason: &messages.Reason{
								Raw: &fakeReasonRaw{},
							},
							RegressionRanges: []*messages.RegressionRange{
								{
									Repo:      "test",
									URL:       "http://test",
									Revisions: []string{"291569"},
									Positions: []string{"refs/heads/master@{#291569}"},
								},
							},
							FinditURL:   "https://findit-for-me.appspot.com/waterfall/failure?url=https://ci.chromium.org/p/fake.master/builders/fake.builder/builds/123",
							HasFindings: false,
							IsFinished:  false,
							IsSupported: true,
						},
					},
				},
			},
			{
				name:         "one build failure with findit",
				master:       "fake.master",
				builder:      "fake.builder",
				tree:         "chromium",
				recentBuilds: []int64{0},
				testData: analyzertest.NewBuilderFaker("fake.master", "fake.builder").
					Build(0).Times(0, 1).IncludeChanges("http://test", "refs/heads/master@{#291569}").
					Step("fake_step").Results(2).BuilderFaker,
				finditData: []*messages.FinditResult{
					{
						MasterURL:                   "https://ci.chromium.org/p/fake.master",
						BuilderName:                 "fake.builder",
						BuildNumber:                 0,
						FirstKnownFailedBuildNumber: 0,
						TryJobStatus:                "FINISHED",
						StepName:                    "fake_step",
						SuspectedCLs: []messages.SuspectCL{
							{
								RepoName:         "test",
								Revision:         "291569",
								Confidence:       90,
								AnalysisApproach: "HEURISTIC",
								RevertCLURL:      "https://codereview/12345",
								RevertCommitted:  true,
							},
						},
						HasFindings: true,
						IsFinished:  true,
						IsSupported: true,
					},
				},
				buildsAtFault: []int{0},
				stepsAtFault:  []string{"fake_step"},
				wantAlerts: []messages.Alert{
					{
						Key:      "fake.master.fake.builder.fake_step.",
						Title:    "fakeTitle",
						Type:     messages.AlertBuildFailure,
						Body:     "",
						Severity: messages.NewFailure,
						Extension: messages.BuildFailure{
							Builders: []messages.AlertedBuilder{
								{
									Name:   "fake.builder",
									URL:    urlParse("https://ci.chromium.org/p/fake.master/builders/fake.builder", t).String(),
									Master: "fake.master",
									Count:  1,
								},
							},
							Reason: &messages.Reason{
								Raw: &fakeReasonRaw{},
							},
							RegressionRanges: []*messages.RegressionRange{
								{
									Repo:      "test",
									URL:       "http://test",
									Revisions: []string{"291569"},
									Positions: []string{"refs/heads/master@{#291569}"},
									RevisionsWithResults: []messages.RevisionWithFinditResult{
										{
											Revision:         "291569",
											IsSuspect:        true,
											Confidence:       90,
											AnalysisApproach: "HEURISTIC",
										},
									},
								},
							},
							SuspectedCLs: []messages.SuspectCL{
								{
									RepoName:         "test",
									Revision:         "291569",
									Confidence:       90,
									AnalysisApproach: "HEURISTIC",
									RevertCLURL:      "https://codereview/12345",
									RevertCommitted:  true,
								},
							},
							FinditStatus: "FINISHED",
							FinditURL:    "https://findit-for-me.appspot.com/waterfall/failure?url=https://ci.chromium.org/p/fake.master/builders/fake.builder/builds/0",
							HasFindings:  true,
							IsFinished:   true,
							IsSupported:  true,
						},
					},
				},
			},
			{
				name:         "one build failure on non-chromium tree",
				master:       "fake.master",
				builder:      "fake.builder",
				tree:         "tree",
				recentBuilds: []int64{0},
				testData: analyzertest.NewBuilderFaker("fake.master", "fake.builder").
					Build(0).Times(0, 1).IncludeChanges("http://test", "refs/heads/master@{#291569}").
					Step("fake_step").Results(2).BuilderFaker,
				buildsAtFault: []int{0},
				stepsAtFault:  []string{"fake_step"},
				wantAlerts: []messages.Alert{
					{
						Key:      "fake.master.fake.builder.fake_step.",
						Title:    "fakeTitle",
						Type:     messages.AlertBuildFailure,
						Body:     "",
						Severity: messages.NewFailure,
						Extension: messages.BuildFailure{
							Builders: []messages.AlertedBuilder{
								{
									Name:   "fake.builder",
									Master: "fake.master",
									URL:    urlParse("https://ci.chromium.org/p/fake.master/builders/fake.builder", t).String(),
									Count:  1,
								},
							},
							Reason: &messages.Reason{
								Raw: &fakeReasonRaw{},
							},
							RegressionRanges: []*messages.RegressionRange{
								{
									Repo:      "test",
									URL:       "http://test",
									Revisions: []string{"291569"},
									Positions: []string{"refs/heads/master@{#291569}"},
								},
							},
						},
					},
				},
			},
			{
				name:         "two build failures with findit",
				master:       "fake.master",
				builder:      "fake.builder",
				tree:         "chromium",
				recentBuilds: []int64{0},
				testData: analyzertest.NewBuilderFaker("fake.master", "fake.builder").
					Build(0).Times(0, 1).IncludeChanges("http://test", "refs/heads/master@{#291569}").
					Step("fake_step").Results(2).BuildFaker.ViewPath("/p/fake.master/builders/fake.builder/1").
					Step("other_step").Results(2).BuilderFaker,
				finditData: []*messages.FinditResult{
					{
						MasterURL:                   "https://ci.chromium.org/p/fake.master",
						BuilderName:                 "fake.builder",
						BuildNumber:                 0,
						FirstKnownFailedBuildNumber: 0,
						TryJobStatus:                "FINISHED",
						StepName:                    "fake_step",
						SuspectedCLs: []messages.SuspectCL{
							{
								RepoName:         "test",
								Revision:         "291569",
								Confidence:       90,
								AnalysisApproach: "HEURISTIC",
							},
						},
						HasFindings: true,
						IsFinished:  true,
						IsSupported: true,
					},
				},
				finditDataV2:  []*messages.FinditResultV2{},
				buildsAtFault: []int{0, 0},
				stepsAtFault:  []string{"fake_step", "other_step"},
				wantAlerts: []messages.Alert{
					{
						Key:      "fake.master.fake.builder.fake_step.",
						Title:    "fakeTitle",
						Type:     messages.AlertBuildFailure,
						Body:     "",
						Severity: messages.NewFailure,
						Extension: messages.BuildFailure{
							Builders: []messages.AlertedBuilder{
								{
									Name:   "fake.builder",
									Master: "fake.master",
									URL:    urlParse("https://ci.chromium.org/p/fake.master/builders/fake.builder", t).String(),
									Count:  1,
								},
							},
							Reason: &messages.Reason{
								Raw: &fakeReasonRaw{},
							},
							RegressionRanges: []*messages.RegressionRange{
								{
									Repo:      "test",
									URL:       "http://test",
									Revisions: []string{"291569"},
									Positions: []string{"refs/heads/master@{#291569}"},
									RevisionsWithResults: []messages.RevisionWithFinditResult{
										{
											Revision:         "291569",
											IsSuspect:        true,
											Confidence:       90,
											AnalysisApproach: "HEURISTIC",
										},
									},
								},
							},
							SuspectedCLs: []messages.SuspectCL{
								{
									RepoName:         "test",
									Revision:         "291569",
									Confidence:       90,
									AnalysisApproach: "HEURISTIC",
								},
							},
							FinditStatus: "FINISHED",
							FinditURL:    "https://findit-for-me.appspot.com/waterfall/failure?url=https://ci.chromium.org/p/fake.master/builders/fake.builder/builds/0",
							HasFindings:  true,
							IsFinished:   true,
							IsSupported:  true,
						},
					},
					{
						Key:      "fake.master.fake.builder.other_step.",
						Title:    "fakeTitle",
						Type:     messages.AlertBuildFailure,
						Body:     "",
						Severity: messages.NewFailure,
						Extension: messages.BuildFailure{
							Builders: []messages.AlertedBuilder{
								{
									Name:   "fake.builder",
									Master: "fake.master",
									URL:    urlParse("https://ci.chromium.org/p/fake.master/builders/fake.builder", t).String(),
									Count:  1,
								},
							},
							Reason: &messages.Reason{
								Raw: &fakeReasonRaw{},
							},
							RegressionRanges: []*messages.RegressionRange{
								{
									Repo:      "test",
									URL:       "http://test",
									Revisions: []string{"291569"},
									Positions: []string{"refs/heads/master@{#291569}"},
								},
							},
						},
					},
				},
			},
			{
				name:         "repeated build failure",
				master:       "fake.master",
				builder:      "fake.builder",
				tree:         "tree",
				recentBuilds: []int64{0, 1, 2, 3},
				testData: analyzertest.NewBuilderFaker("fake.master", "fake.builder").
					Build(0).Times(0, 1).IncludeChanges("http://test", "refs/heads/master@{#291569}").
					Step("fake_step").Results(2).BuilderFaker.
					Build(1).Times(2, 3).IncludeChanges("http://test", "refs/heads/master@{#291570}").
					Step("fake_step").Results(2).BuilderFaker.
					Build(2).Times(4, 5).IncludeChanges("http://test", "refs/heads/master@{#291570}").
					Step("fake_step").Results(2).BuilderFaker.
					Build(3).Times(6, 7).IncludeChanges("http://test", "refs/heads/master@{#291570}").
					Step("fake_step").Results(2).BuilderFaker,
				buildsAtFault: []int{3},
				stepsAtFault:  []string{"fake_step"},
				wantAlerts: []messages.Alert{
					{
						Key:      "fake.master.fake.builder.fake_step.",
						Title:    "fakeTitle",
						Type:     messages.AlertBuildFailure,
						Body:     "",
						Severity: messages.ReliableFailure,
						Time:     messages.EpochTime(6),
						Extension: messages.BuildFailure{
							Builders: []messages.AlertedBuilder{
								{
									Name:          "fake.builder",
									Master:        "fake.master",
									URL:           urlParse("https://ci.chromium.org/p/fake.master/builders/fake.builder", t).String(),
									FirstFailure:  0,
									LatestFailure: 3,
									Count:         4,
								},
							},
							Reason: &messages.Reason{
								Raw: &fakeReasonRaw{},
							},
							RegressionRanges: []*messages.RegressionRange{
								{
									Repo: "test",
									URL:  "http://test",
									Revisions: []string{
										"291569",
									},
									Positions: []string{
										"refs/heads/master@{#291569}",
									},
								},
							},
						},
					},
				},
			},
			{
				name:         "new double failures counted",
				master:       "fake.master",
				builder:      "fake.builder",
				tree:         "tree",
				recentBuilds: []int64{0, 1, 2},
				testData: analyzertest.NewBuilderFaker("fake.master", "fake.builder").
					Build(0).Times(0, 1).IncludeChanges("http://test", "refs/heads/master@{#291569}").
					Step("fake_step").Results(2).BuilderFaker.
					Build(1).Times(2, 3).IncludeChanges("http://test", "refs/heads/master@{#291570}").
					Step("fake_step").Results(2).BuilderFaker.
					Build(2).Times(4, 5).IncludeChanges("http://test", "refs/heads/master@{#291570}").
					Step("fake_step").Results(2).BuildFaker.ViewPath("/p/fake.master/builders/fake.builder/1").
					Step("other_step").Results(2).BuilderFaker,
				buildsAtFault: []int{2, 2},
				stepsAtFault:  []string{"other_step", "fake_step"},
				wantAlerts: []messages.Alert{
					{
						Key:       "fake.master.fake.builder.other_step.",
						Title:     "fakeTitle",
						Type:      messages.AlertBuildFailure,
						Body:      "",
						Time:      messages.EpochTime(4),
						StartTime: messages.EpochTime(4),
						Severity:  messages.NewFailure,
						Extension: messages.BuildFailure{
							Builders: []messages.AlertedBuilder{
								{
									Name:          "fake.builder",
									Master:        "fake.master",
									URL:           urlParse("https://ci.chromium.org/p/fake.master/builders/fake.builder", t).String(),
									StartTime:     messages.EpochTime(4),
									FirstFailure:  2,
									LatestFailure: 2,
									Count:         1,
								},
							},
							Reason: &messages.Reason{
								Raw: &fakeReasonRaw{},
							},
							RegressionRanges: []*messages.RegressionRange{
								{
									Repo:      "test",
									URL:       "http://test",
									Revisions: []string{"291570"},
									Positions: []string{"refs/heads/master@{#291570}"},
								},
							},
						},
					},
					{
						Key:      "fake.master.fake.builder.fake_step.",
						Title:    "fakeTitle",
						Type:     messages.AlertBuildFailure,
						Body:     "",
						Severity: messages.ReliableFailure,
						Time:     messages.EpochTime(4),
						Extension: messages.BuildFailure{
							Builders: []messages.AlertedBuilder{
								{
									Name:          "fake.builder",
									Master:        "fake.master",
									URL:           urlParse("https://ci.chromium.org/p/fake.master/builders/fake.builder", t).String(),
									FirstFailure:  0,
									LatestFailure: 2,
									Count:         3,
								},
							},
							Reason: &messages.Reason{
								Raw: &fakeReasonRaw{},
							},
							RegressionRanges: []*messages.RegressionRange{
								{
									Repo: "test",
									URL:  "http://test",
									Revisions: []string{
										"291569",
									},
									Positions: []string{
										"refs/heads/master@{#291569}",
									},
								},
							},
						},
					},
				},
			},
			{
				name:         "old failures not counted",
				master:       "fake.master",
				builder:      "fake.builder",
				tree:         "tree",
				recentBuilds: []int64{0, 1, 2},
				testData: analyzertest.NewBuilderFaker("fake.master", "fake.builder").
					Build(0).Times(0, 1).IncludeChanges("http://test", "refs/heads/master@{#291569}").
					Step("fake_step").Results(2).BuildFaker.ViewPath("/p/fake.master/builders/fake.builder/1").
					Step("other_step").Results(2).BuilderFaker.
					Build(1).Times(2, 3).IncludeChanges("http://test", "refs/heads/master@{#291570}").
					Step("fake_step").Results(2).BuilderFaker.
					Build(2).Times(4, 5).IncludeChanges("http://test", "refs/heads/master@{#291571}").
					Step("fake_step").Results(2).BuilderFaker,
				buildsAtFault: []int{2},
				stepsAtFault:  []string{"fake_step"},
				wantAlerts: []messages.Alert{
					{
						Key:      "fake.master.fake.builder.fake_step.",
						Title:    "fakeTitle",
						Type:     messages.AlertBuildFailure,
						Body:     "",
						Severity: messages.ReliableFailure,
						Time:     messages.EpochTime(4),
						Extension: messages.BuildFailure{
							Builders: []messages.AlertedBuilder{
								{
									Name:          "fake.builder",
									Master:        "fake.master",
									URL:           urlParse("https://ci.chromium.org/p/fake.master/builders/fake.builder", t).String(),
									FirstFailure:  0,
									LatestFailure: 2,
									Count:         3,
								},
							},
							Reason: &messages.Reason{
								Raw: &fakeReasonRaw{},
							},
							RegressionRanges: []*messages.RegressionRange{
								{
									Repo: "test",
									URL:  "http://test",
									Revisions: []string{
										"291569",
									},
									Positions: []string{
										"refs/heads/master@{#291569}",
									},
								},
							},
						},
					},
				},
			},
			{
				name:         "green history, failure recently",
				master:       "fake.master",
				builder:      "fake.builder",
				tree:         "tree",
				recentBuilds: []int64{0, 1, 2, 3, 4, 5},
				testData: analyzertest.NewBuilderFaker("fake.master", "fake.builder").
					Build(0).Times(0, 1).IncludeChanges("http://test", "refs/heads/master@{#291569}").
					Step("fake_step").Results(0).BuilderFaker.
					Build(1).Times(2, 3).IncludeChanges("http://test", "refs/heads/master@{#291570}").
					Step("fake_step").Results(0).BuilderFaker.
					Build(2).Times(4, 5).IncludeChanges("http://test", "refs/heads/master@{#291571}").
					Step("fake_step").Results(0).BuilderFaker.
					Build(3).Times(6, 7).IncludeChanges("http://test", "refs/heads/master@{#291572}").
					Step("fake_step").Results(0).BuilderFaker.
					Build(4).Times(8, 9).IncludeChanges("http://test", "refs/heads/master@{#291573}").
					Step("fake_step").Results(2).BuilderFaker.
					Build(5).Times(10, 11).IncludeChanges("http://test", "refs/heads/master@{#291574}").
					Step("fake_step").Results(2).BuilderFaker,
				buildsAtFault: []int{5},
				stepsAtFault:  []string{"fake_step"},
				wantAlerts: []messages.Alert{
					{
						Key:       "fake.master.fake.builder.fake_step.",
						Title:     "fakeTitle",
						Type:      messages.AlertBuildFailure,
						Body:      "",
						Severity:  messages.ReliableFailure,
						StartTime: messages.EpochTime(8),
						Time:      messages.EpochTime(10),
						Extension: messages.BuildFailure{
							Builders: []messages.AlertedBuilder{
								{
									Name:          "fake.builder",
									Master:        "fake.master",
									URL:           urlParse("https://ci.chromium.org/p/fake.master/builders/fake.builder", t).String(),
									FirstFailure:  4,
									LatestFailure: 5,
									LatestPassing: 3,
									Count:         2,
									StartTime:     messages.EpochTime(8),
								},
							},
							Reason: &messages.Reason{
								Raw: &fakeReasonRaw{},
							},
							RegressionRanges: []*messages.RegressionRange{
								{
									Repo: "test",
									URL:  "http://test",
									Revisions: []string{
										"291573",
									},
									Positions: []string{
										"refs/heads/master@{#291573}",
									},
								},
							},
						},
					},
				},
			},
		}

		a := newTestAnalyzer(0, 10)

		for _, test := range tests {
			test := test
			ctx := context.Background()
			Convey(test.name, func() {
				a.Now = fakeNow(time.Unix(0, 0))
				var builds map[string]*messages.Build
				if test.testData != nil {
					builds = test.testData.Builds
				}
				mc := clientTest.MockReader{
					Builds:          builds,
					FinditResults:   test.finditData,
					FinditResultsV2: test.finditDataV2,
				}
				a.Milo = mc
				a.FindIt = mc
				So(test.buildsAtFault, ShouldHaveLength, len(test.wantAlerts))
				So(test.stepsAtFault, ShouldHaveLength, len(test.wantAlerts))

				newAlerts := []messages.Alert(nil)
				stepsAtFault := analyzertest.StepsAtFault(test.testData, test.buildsAtFault, test.stepsAtFault)
				for i, alr := range test.wantAlerts {
					ext := alr.Extension.(messages.BuildFailure)
					ext.StepAtFault = &stepsAtFault[i]
					alr.Extension = ext
					newAlerts = append(newAlerts, alr)
				}
				test.wantAlerts = newAlerts

				gotAlerts, gotErrs := a.builderStepAlerts(ctx, test.tree, &messages.MasterLocation{URL: *urlParse("https://ci.chromium.org/p/"+test.master, t)}, test.builder, test.recentBuilds)

				sort.Sort(sortAlerts(gotAlerts))
				sort.Sort(sortAlerts(test.wantAlerts))

				So(gotAlerts, ShouldResemble, test.wantAlerts)
				if !reflect.DeepEqual(gotErrs, test.wantErrs) {
					t.Errorf("%s failed. Got %+v, want: %+v", test.name, gotErrs, test.wantErrs)
				}
			})
		}
	})
}

type sortAlerts []messages.Alert

func (a sortAlerts) Len() int           { return len(a) }
func (a sortAlerts) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortAlerts) Less(i, j int) bool { return a[i].Key > a[j].Key }

func TestStepFailures(t *testing.T) {

	Convey("test StepFailures", t, func() {
		tests := []struct {
			name            string
			master, builder string
			b               *messages.Build
			buildNum        int64
			bCache          map[string]*messages.Build
			want            []*messages.BuildStep
			wantErr         error
		}{
			{
				name:    "empty",
				master:  "fake.master",
				builder: "fake.builder",
			},
			{
				name:     "breaking step",
				master:   "stepCheck.master",
				builder:  "fake.builder",
				buildNum: 0,
				bCache: analyzertest.NewBuilderFaker("stepCheck.master", "fake.builder").
					Build(0).Step("ok_step").Results(0).BuildFaker.ViewPath("/p/stepCheck.master/builders/fake.builder/1").
					Step("broken_step").Results(3).BuilderFaker.Builds,
				want: []*messages.BuildStep{
					{
						Master: &messages.MasterLocation{URL: *urlParse("https://ci.chromium.org/p/stepCheck.master", t)},
						Build: &messages.Build{
							Master:      "stepCheck.master",
							BuilderName: "fake.builder",
							ViewPath:    "/p/stepCheck.master/builders/fake.builder/1",
							Steps: []messages.Step{
								{
									Name:       "ok_step",
									IsFinished: true,
									Results:    []interface{}{float64(0)},
								},
								{
									Name:       "broken_step",
									IsFinished: true,
									Results:    []interface{}{float64(3)},
								},
							},
						},
						Step: &messages.Step{
							Name:       "broken_step",
							IsFinished: true,
							Results:    []interface{}{float64(3)},
						},
					},
				},
			},
		}

		mc := &clientTest.MockReader{}
		ctx := context.Background()
		a := newTestAnalyzer(0, 10)
		a.Milo = mc

		for _, test := range tests {
			test := test
			Convey(test.name, func() {
				mc.BuildValue = test.b
				mc.BCache = test.bCache
				got, err := a.stepFailures(ctx, &messages.MasterLocation{URL: *urlParse("https://ci.chromium.org/p/"+test.master, t)}, test.builder, test.buildNum)
				So(got, ShouldResemble, test.want)
				So(err, ShouldResemble, test.wantErr)
			})
		}
	})
}

func TestStepFailureAlerts(t *testing.T) {

	Convey("test StepFailureAlerts", t, func() {
		tests := []struct {
			name        string
			failures    []*messages.BuildStep
			testResults *model.FullResult
			alerts      []messages.Alert
			err         error
		}{
			{
				name:   "empty",
				alerts: []messages.Alert{},
			},
			{
				name: "single failure",
				failures: []*messages.BuildStep{
					{
						Master: &messages.MasterLocation{URL: *urlParse("https://ci.chromium.org/p/fake.master", t)},
						Build: &messages.Build{
							Master:      "fake.master",
							BuilderName: "fake.builder",
							Number:      2,
							Times:       []messages.EpochTime{0, 1},
							ViewPath:    "/p/fake.master/builders/fake.builder/1",
						},
						Step: &messages.Step{
							Name: "steps",
						},
					},
					{
						Master: &messages.MasterLocation{URL: *urlParse("https://ci.chromium.org/p/fake.master", t)},
						Build: &messages.Build{
							Master:      "fake.master",
							BuilderName: "fake.builder",
							Number:      42,
							Times:       []messages.EpochTime{0, 1},
							ViewPath:    "/p/fake.master/builders/fake.builder/1",
						},
						Step: &messages.Step{
							Name: "fake_tests",
						},
					},
				},
				testResults: &model.FullResult{},
				alerts: []messages.Alert{
					{
						Key:      "fake.master.fake.builder.fake_tests.",
						Title:    "fakeTitle",
						Body:     "",
						Severity: messages.NewFailure,
						Type:     messages.AlertBuildFailure,
						Extension: messages.BuildFailure{
							Builders: []messages.AlertedBuilder{
								{
									Name:          "fake.builder",
									Master:        "fake.master",
									URL:           urlParse("https://ci.chromium.org/p/fake.master/builders/fake.builder", t).String(),
									FirstFailure:  42,
									LatestFailure: 42,
									Count:         1,
								},
							},
							Reason: &messages.Reason{
								Raw: &fakeReasonRaw{},
							},
							StepAtFault: &messages.BuildStep{
								Master: &messages.MasterLocation{URL: *urlParse("https://ci.chromium.org/p/fake.master", t)},
								Build: &messages.Build{
									Master:      "fake.master",
									BuilderName: "fake.builder",
									Number:      42,
									Times:       []messages.EpochTime{0, 1},
									ViewPath:    "/p/fake.master/builders/fake.builder/1",
								},
								Step: &messages.Step{
									Name: "fake_tests",
								},
							},
						},
					},
				},
			},
			{
				name: "single failure (step links)",
				failures: []*messages.BuildStep{
					{
						Master: &messages.MasterLocation{URL: *urlParse("https://ci.chromium.org/p/fake.master", t)},
						Build: &messages.Build{
							Master:      "fake.master",
							BuilderName: "fake.builder",
							Number:      2,
							Times:       []messages.EpochTime{0, 1},
							ViewPath:    "/p/fake.master/builders/fake.builder/1",
						},
						Step: &messages.Step{
							Name: "steps",
						},
					},
					{
						Master: &messages.MasterLocation{URL: *urlParse("https://ci.chromium.org/p/fake.master", t)},
						Build: &messages.Build{
							Master:      "fake.master",
							BuilderName: "fake.builder",
							Number:      42,
							Times:       []messages.EpochTime{0, 1},
							ViewPath:    "/p/fake.master/builders/fake.builder/1",
						},
						Step: &messages.Step{
							Name: "fake_tests",
							Links: map[string]string{
								"Screenshot": "https://console.developers.google.com/cloud_storage",
							},
						},
					},
				},
				alerts: []messages.Alert{
					{
						Key:      "fake.master.fake.builder.fake_tests.",
						Title:    "fakeTitle",
						Body:     "",
						Severity: messages.NewFailure,
						Type:     messages.AlertBuildFailure,
						Extension: messages.BuildFailure{
							Builders: []messages.AlertedBuilder{
								{
									Name:          "fake.builder",
									Master:        "fake.master",
									URL:           urlParse("https://ci.chromium.org/p/fake.master/builders/fake.builder", t).String(),
									FirstFailure:  42,
									LatestFailure: 42,
									Count:         1,
								},
							},
							Reason: &messages.Reason{
								Raw: &fakeReasonRaw{},
							},
							StepAtFault: &messages.BuildStep{
								Master: &messages.MasterLocation{URL: *urlParse("https://ci.chromium.org/p/fake.master", t)},
								Build: &messages.Build{
									Master:      "fake.master",
									BuilderName: "fake.builder",
									Number:      42,
									Times:       []messages.EpochTime{0, 1},
									ViewPath:    "/p/fake.master/builders/fake.builder/1",
								},
								Step: &messages.Step{
									Name: "fake_tests",
									Links: map[string]string{
										"Screenshot": "https://console.developers.google.com/cloud_storage",
									},
								},
							},
						},
						Links: []messages.Link{
							{
								Title: "Screenshot",
								Href:  "https://console.developers.google.com/cloud_storage",
							},
						},
					},
				},
			},
			{
				name: "single failure (blacklisted step links)",
				failures: []*messages.BuildStep{
					{
						Master: &messages.MasterLocation{URL: *urlParse("https://ci.chromium.org/p/fake.master", t)},
						Build: &messages.Build{
							Master:      "fake.master",
							BuilderName: "fake.builder",
							Number:      2,
							Times:       []messages.EpochTime{0, 1},
							ViewPath:    "/p/fake.master/builders/fake.builder/1",
						},
						Step: &messages.Step{
							Name: "steps",
						},
					},
					{
						Master: &messages.MasterLocation{URL: *urlParse("https://ci.chromium.org/p/fake.master", t)},
						Build: &messages.Build{
							Master:      "fake.master",
							BuilderName: "fake.builder",
							Number:      42,
							Times:       []messages.EpochTime{0, 1},
							ViewPath:    "/p/fake.master/builders/fake.builder/1",
						},
						Step: &messages.Step{
							Name: "fake_tests",
							Links: map[string]string{
								"The best site for free stuff": "https://best.site/all_the_stuff",
							},
						},
					},
				},
				testResults: &model.FullResult{},
				alerts: []messages.Alert{
					{
						Key:      "fake.master.fake.builder.fake_tests.",
						Title:    "fakeTitle",
						Body:     "",
						Severity: messages.NewFailure,
						Type:     messages.AlertBuildFailure,
						Extension: messages.BuildFailure{
							Builders: []messages.AlertedBuilder{
								{
									Name:          "fake.builder",
									Master:        "fake.master",
									URL:           urlParse("https://ci.chromium.org/p/fake.master/builders/fake.builder", t).String(),
									FirstFailure:  42,
									LatestFailure: 42,
									Count:         1,
								},
							},
							Reason: &messages.Reason{
								Raw: &fakeReasonRaw{},
							},
							StepAtFault: &messages.BuildStep{
								Master: &messages.MasterLocation{URL: *urlParse("https://ci.chromium.org/p/fake.master", t)},
								Build: &messages.Build{
									Master:      "fake.master",
									BuilderName: "fake.builder",
									Number:      42,
									Times:       []messages.EpochTime{0, 1},
									ViewPath:    "/p/fake.master/builders/fake.builder/1",
								},
								Step: &messages.Step{
									Name: "fake_tests",
									Links: map[string]string{
										"The best site for free stuff": "https://best.site/all_the_stuff",
									},
								},
							},
						},
					},
				},
			},
			{
				name: "single failure (weird test suite name)",
				failures: []*messages.BuildStep{
					{
						Master: &messages.MasterLocation{URL: *urlParse("https://ci.chromium.org/p/fake.master", t)},
						Build: &messages.Build{
							Master:      "fake.master",
							BuilderName: "fake.builder",
							Number:      2,
							Times:       []messages.EpochTime{0, 1},
							ViewPath:    "/p/fake.master/builders/fake.builder/1",
						},
						Step: &messages.Step{
							Name: "steps",
						},
					},
					{
						Master: &messages.MasterLocation{URL: *urlParse("https://ci.chromium.org/p/fake.master", t)},
						Build: &messages.Build{
							Master:      "fake.master",
							BuilderName: "fake.builder",
							Number:      42,
							Times:       []messages.EpochTime{0, 1},
							ViewPath:    "/p/fake.master/builders/fake.builder/1",
						},
						Step: &messages.Step{
							Name: "fake_tests on Intel GPU on Linux on Ubuntu-12.04",
						},
					},
				},
				testResults: &model.FullResult{},
				alerts: []messages.Alert{
					{
						Key:      "fake.master.fake.builder.fake_tests.",
						Title:    "fakeTitle",
						Body:     "",
						Severity: messages.NewFailure,
						Type:     messages.AlertBuildFailure,
						Extension: messages.BuildFailure{
							Builders: []messages.AlertedBuilder{
								{
									Name:          "fake.builder",
									Master:        "fake.master",
									URL:           urlParse("https://ci.chromium.org/p/fake.master/builders/fake.builder", t).String(),
									FirstFailure:  42,
									LatestFailure: 42,
									Count:         1,
								},
							},
							Reason: &messages.Reason{
								Raw: &fakeReasonRaw{},
							},
							StepAtFault: &messages.BuildStep{
								Master: &messages.MasterLocation{URL: *urlParse("https://ci.chromium.org/p/fake.master", t)},
								Build: &messages.Build{
									Master:      "fake.master",
									BuilderName: "fake.builder",
									Number:      42,
									Times:       []messages.EpochTime{0, 1},
									ViewPath:    "/p/fake.master/builders/fake.builder/1",
								},
								Step: &messages.Step{
									Name: "fake_tests on Intel GPU on Linux on Ubuntu-12.04",
								},
							},
						},
					},
				},
			},
			{
				name: "single infra failure",
				failures: []*messages.BuildStep{
					{
						Master: &messages.MasterLocation{URL: *urlParse("https://ci.chromium.org/p/fake.master", t)},
						Build: &messages.Build{
							Master:      "fake.master",
							BuilderName: "fake.builder",
							Number:      2,
							Times:       []messages.EpochTime{0, 1},
							ViewPath:    "/p/fake.master/builders/fake.builder/1",
						},
						Step: &messages.Step{
							Name: "steps",
						},
					},
					{
						Master: &messages.MasterLocation{URL: *urlParse("https://ci.chromium.org/p/fake.master", t)},
						Build: &messages.Build{
							Master:      "fake.master",
							BuilderName: "fake.builder",
							Number:      42,
							Times:       []messages.EpochTime{0, 1},
							ViewPath:    "/p/fake.master/builders/fake.builder/1",
						},
						Step: &messages.Step{
							Name:    "fake_tests",
							Results: []interface{}{float64(resInfraFailure)},
							Times:   []messages.EpochTime{0, 1},
						},
					},
				},
				testResults: &model.FullResult{},
				alerts: []messages.Alert{
					{
						Key:      "fake.master.fake.builder.fake_tests.4",
						Title:    "fake_tests failing on fake.master/fake.builder",
						Body:     "infrastructure failure",
						Severity: messages.InfraFailure,
						Type:     messages.AlertInfraFailure,
						Extension: messages.BuildFailure{
							Builders: []messages.AlertedBuilder{
								{
									Name:          "fake.builder",
									Master:        "fake.master",
									URL:           urlParse("https://ci.chromium.org/p/fake.master/builders/fake.builder", t).String(),
									FirstFailure:  42,
									LatestFailure: 42,
								},
							},
							Reason: &messages.Reason{
								Raw: &fakeReasonRaw{},
							},
							StepAtFault: &messages.BuildStep{
								Master: &messages.MasterLocation{URL: *urlParse("https://ci.chromium.org/p/fake.master", t)},
								Build: &messages.Build{
									Master:      "fake.master",
									BuilderName: "fake.builder",
									Number:      42,
									Times:       []messages.EpochTime{0, 1},
									ViewPath:    "/p/fake.master/builders/fake.builder/1",
								},
								Step: &messages.Step{
									Name:    "fake_tests",
									Results: []interface{}{float64(resInfraFailure)},
									Times:   []messages.EpochTime{0, 1},
								},
							},
						},
					},
				},
			},
		}

		mc := &clientTest.MockReader{}
		ctx := context.Background()

		a := newTestAnalyzer(0, 10)
		a.Now = fakeNow(time.Unix(0, 0))

		for _, test := range tests {
			test := test
			Convey(test.name, func() {
				mc.TestResultsValue = test.testResults
				alerts, err := a.stepFailureAlerts(ctx, "tree", test.failures, []*messages.FinditResult{})
				So(err, ShouldResemble, test.err)
				So(alerts, ShouldResemble, test.alerts)
			})
		}
	})
}

func TestLatestBuildStep(t *testing.T) {
	tests := []struct {
		name       string
		b          messages.Build
		wantStep   string
		wantUpdate messages.EpochTime
		wantErr    error
	}{
		{
			name:    "blank",
			wantErr: errNoBuildSteps,
		},
		{
			name: "done time is latest",
			b: messages.Build{
				Times: []messages.EpochTime{0},
				Steps: []messages.Step{
					{
						Name: "done step",
						Times: []messages.EpochTime{
							messages.TimeToEpochTime(time.Unix(6, 0)),
							messages.TimeToEpochTime(time.Unix(42, 0)),
						},
					},
				},
			},
			wantStep:   "done step",
			wantUpdate: messages.TimeToEpochTime(time.Unix(42, 0)),
		},
		{
			name: "started time is latest",
			b: messages.Build{
				Times: []messages.EpochTime{0},
				Steps: []messages.Step{
					{
						Name: "start step",
						Times: []messages.EpochTime{
							messages.TimeToEpochTime(time.Unix(42, 0)),
						},
					},
				},
			},
			wantStep:   "start step",
			wantUpdate: messages.TimeToEpochTime(time.Unix(42, 0)),
		},
		{
			name: "started time is latest, multiple steps",
			b: messages.Build{
				Times: []messages.EpochTime{0},
				Steps: []messages.Step{
					{
						Name: "start step",
						Times: []messages.EpochTime{
							messages.TimeToEpochTime(time.Unix(6, 0)),
							messages.TimeToEpochTime(time.Unix(7, 0)),
						},
					},
					{
						Name: "second step",
						Times: []messages.EpochTime{
							messages.TimeToEpochTime(time.Unix(42, 0)),
							messages.TimeToEpochTime(time.Unix(0, 0)),
						},
					},
				},
			},
			wantStep:   "second step",
			wantUpdate: messages.TimeToEpochTime(time.Unix(42, 0)),
		},
		{
			name: "done time is latest, multiple steps",
			b: messages.Build{
				Times: []messages.EpochTime{0},
				Steps: []messages.Step{
					{
						Name: "start step",
						Times: []messages.EpochTime{
							messages.TimeToEpochTime(time.Unix(0, 0)),
							messages.TimeToEpochTime(time.Unix(6, 0)),
						},
					},
					{
						Name: "second step",
						Times: []messages.EpochTime{
							messages.TimeToEpochTime(time.Unix(7, 0)),
							messages.TimeToEpochTime(time.Unix(42, 0)),
						},
					},
				},
			},
			wantStep:   "second step",
			wantUpdate: messages.TimeToEpochTime(time.Unix(42, 0)),
		},
		{
			name: "build is done",
			b: messages.Build{
				Times: []messages.EpochTime{
					messages.TimeToEpochTime(time.Unix(0, 0)),
					messages.TimeToEpochTime(time.Unix(42, 0)),
				},
				Steps: []messages.Step{
					{
						Name: "start step",
						Times: []messages.EpochTime{
							messages.TimeToEpochTime(time.Unix(0, 0)),
							messages.TimeToEpochTime(time.Unix(0, 0)),
						},
					},
					{
						Name: "second step",
						Times: []messages.EpochTime{
							messages.TimeToEpochTime(time.Unix(6, 0)),
							messages.TimeToEpochTime(time.Unix(7, 0)),
						},
					},
				},
			},
			wantStep:   StepCompletedRun,
			wantUpdate: messages.TimeToEpochTime(time.Unix(42, 0)),
		},
	}

	a := newTestAnalyzer(0, 10)
	a.Now = fakeNow(time.Unix(0, 0))
	for _, test := range tests {
		gotStep, gotUpdate, gotErr := a.latestBuildStep(&test.b)
		if gotStep != test.wantStep {
			t.Errorf("%s failed. Got %q, want %q.", test.name, gotStep, test.wantStep)
		}
		if gotUpdate != test.wantUpdate {
			t.Errorf("%s failed. Got %v, want %v.", test.name, gotUpdate, test.wantUpdate)
		}
		if gotErr != test.wantErr {
			t.Errorf("%s failed. Got %s, want %s.", test.name, gotErr, test.wantErr)
		}
	}
}

func TestWouldCloseTree(t *testing.T) {
	ctx := context.Background()
	Convey("gatekepeer", t, func() {
		gkr := NewGatekeeperRules(ctx, []*messages.GatekeeperConfig{
			{
				Masters: map[string][]messages.MasterConfig{
					"https://ci.chromium.org/p/fake.master": {{
						Builders: map[string]messages.BuilderConfig{
							"fake.builder": {
								ClosingSteps: []string{"*"},
							},
						},
						ExcludedBuilders: []string{"other.builder"},
					}},
				},
			},
		}, map[string][]messages.TreeMasterConfig{
			"test_tree": {
				messages.TreeMasterConfig{
					Masters: map[messages.MasterLocation][]string{
						{URL: *urlParse(
							"https://ci.chromium.org/p/fake.master", t)}: {"fake.builder", "other.builder"},
					},
				},
			},
		})

		loc := &messages.MasterLocation{URL: *urlParse("https://ci.chromium.org/p/fake.master", t)}
		So(gkr.WouldCloseTree(ctx, loc, "fake.builder", "fake.step"), ShouldEqual, true)
		So(gkr.WouldCloseTree(ctx, loc, "other.builder", "fake.step"), ShouldEqual, false)
	})
}

func TestExcludeFailure(t *testing.T) {
	tests := []struct {
		name                        string
		gk                          messages.GatekeeperConfig
		gkt                         map[string][]messages.TreeMasterConfig
		master, builder, step, tree string
		want                        bool
	}{
		{
			name:    "empty config",
			tree:    "test_tree1",
			master:  "fake.master",
			builder: "fake.builder",
			step:    "fake_step",
			want:    false,
		},
		{
			name:    "specifically excluded builder",
			tree:    "test_tree2",
			master:  "fake.master",
			builder: "fake.builder",
			step:    "fake_step",
			gk: messages.GatekeeperConfig{Masters: map[string][]messages.MasterConfig{
				"https://ci.chromium.org/p/fake.master": {{
					ExcludedBuilders: []string{"fake.builder"},
				}},
			}},
			want: true,
		},
		{
			name:    "specifically excluded master step",
			tree:    "test_tree3",
			master:  "fake.master",
			builder: "fake.builder",
			step:    "fake_step",
			gk: messages.GatekeeperConfig{Masters: map[string][]messages.MasterConfig{
				"https://ci.chromium.org/p/fake.master": {{
					ExcludedSteps: []string{"fake_step"},
				}},
			}},
			want: true,
		},
		{
			name:    "specifically excluded builder step",
			tree:    "test_tree4",
			master:  "fake.master",
			builder: "fake.builder",
			step:    "fake_step",
			gk: messages.GatekeeperConfig{Masters: map[string][]messages.MasterConfig{
				"https://ci.chromium.org/p/fake.master": {{
					Builders: map[string]messages.BuilderConfig{
						"fake.builder": {
							ExcludedSteps: []string{"fake_step"},
						},
					}},
				}},
			},
			want: true,
		},
		{
			name:    "wildcard builder excluded",
			tree:    "test_tree5",
			master:  "fake.master",
			builder: "fake.builder",
			step:    "fake_step",
			gk: messages.GatekeeperConfig{Masters: map[string][]messages.MasterConfig{
				"https://ci.chromium.org/p/fake.master": {{
					ExcludedBuilders: []string{"*"},
				}},
			}},
			want: true,
		},
		{
			name:    "config should exclude builder (tree config)",
			tree:    "test_tree6",
			master:  "fake.master",
			builder: "fake.builder",
			step:    "fake_step",
			gk: messages.GatekeeperConfig{Masters: map[string][]messages.MasterConfig{
				"https://ci.chromium.org/p/fake.master": {{
					Builders: map[string]messages.BuilderConfig{
						"*": {},
					},
				}},
			}},
			gkt: map[string][]messages.TreeMasterConfig{
				"test_tree": {
					messages.TreeMasterConfig{
						Masters: map[messages.MasterLocation][]string{
							{URL: *urlParse(
								"https://ci.chromium.org/p/fake.master", t)}: {"other.builder"},
						},
					},
				},
			},
			want: true,
		},
		{
			name:    "config shouldn't exclude builder (tree config)",
			tree:    "test_tree7",
			master:  "fake.master",
			builder: "fake.builder",
			step:    "fake_step",
			gk: messages.GatekeeperConfig{Masters: map[string][]messages.MasterConfig{
				"https://ci.chromium.org/p/fake.master": {{
					Builders: map[string]messages.BuilderConfig{
						"*": {},
					},
				}},
			}},
			gkt: map[string][]messages.TreeMasterConfig{
				"test_tree7": {
					messages.TreeMasterConfig{
						Masters: map[messages.MasterLocation][]string{
							{URL: *urlParse(
								"https://ci.chromium.org/p/fake.master", t)}: {"fake.builder"},
						},
					},
				},
			},
			want: false,
		},
		{
			name:    "config shouldn't exclude builder (tree config glob)",
			tree:    "test_tree8",
			master:  "fake.master",
			builder: "fake.builder",
			step:    "fake_step",
			gk: messages.GatekeeperConfig{Masters: map[string][]messages.MasterConfig{
				"https://ci.chromium.org/p/fake.master": {{
					Builders: map[string]messages.BuilderConfig{
						"*": {},
					},
				}},
			}},
			gkt: map[string][]messages.TreeMasterConfig{
				"test_tree8": {
					messages.TreeMasterConfig{
						Masters: map[messages.MasterLocation][]string{
							{URL: *urlParse(
								"https://ci.chromium.org/p/fake.master", t)}: {"*"},
						},
					},
				},
			},
			want: false,
		},
		{
			name:    "partial glob step excluded",
			tree:    "test_tree9",
			master:  "fake.master",
			builder: "fake.builder",
			step:    "fake_step (experimental)",
			gk: messages.GatekeeperConfig{Masters: map[string][]messages.MasterConfig{
				"https://ci.chromium.org/p/fake.master": {{
					Builders: map[string]messages.BuilderConfig{
						"*": {},
					},
					ExcludedSteps: []string{
						"* (experimental)",
					},
				}},
			}},
			gkt: map[string][]messages.TreeMasterConfig{
				"test_tree9": {
					messages.TreeMasterConfig{
						Masters: map[messages.MasterLocation][]string{
							{URL: *urlParse(
								"https://ci.chromium.org/p/fake.master", t)}: {"*"},
						},
					},
				},
			},
			want: true,
		},
		{
			name:    "partial glob step excluded by builder category",
			tree:    "test_tree10",
			master:  "fake.master",
			builder: "fake.builder",
			step:    "fake_step (experimental)",
			gk: messages.GatekeeperConfig{
				Masters: map[string][]messages.MasterConfig{
					"https://ci.chromium.org/p/fake.master": {{
						Builders: map[string]messages.BuilderConfig{
							"fake.builder": {
								Categories: []string{
									"experimental_tests",
								},
							},
						},
					}},
				},
				Categories: map[string]messages.CategoryConfig{
					"experimental_tests": {
						ExcludedSteps: []string{
							"* (experimental)",
						},
					},
				},
			},
			want: true,
		},
		{
			name:    "partial glob step excluded by master category",
			tree:    "test_tree11",
			master:  "fake.master",
			builder: "fake.builder",
			step:    "fake_step (experimental)",
			gk: messages.GatekeeperConfig{
				Masters: map[string][]messages.MasterConfig{
					"https://ci.chromium.org/p/fake.master": {{
						Builders: map[string]messages.BuilderConfig{
							"*": {},
						},
						Categories: []string{
							"experimental_tests",
						},
					}},
				},
				Categories: map[string]messages.CategoryConfig{
					"experimental_tests": {
						ExcludedSteps: []string{
							"* (experimental)",
						},
					},
				},
			},
			want: true,
		},
	}

	ctx := context.Background()

	a := newTestAnalyzer(0, 10)
	for _, test := range tests {
		a.Gatekeeper = NewGatekeeperRules(ctx, []*messages.GatekeeperConfig{&test.gk}, test.gkt)
		got := a.Gatekeeper.ExcludeFailure(ctx, test.tree, &messages.MasterLocation{URL: *urlParse("https://ci.chromium.org/p/"+test.master, t)}, test.builder, test.step)
		if got != test.want {
			t.Errorf("%s failed. Got: %+v, want: %+v", test.name, got, test.want)
		}
	}
}
