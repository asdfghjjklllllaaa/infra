package analyzer

import (
	"encoding/json"
	"fmt"
	"infra/monitoring/messages"
	"time"

	bbpb "go.chromium.org/luci/buildbucket/proto"
	"go.chromium.org/luci/common/data/stringset"
	"go.chromium.org/luci/common/logging"
	"golang.org/x/net/context"
)

func stepSet(s []*bbpb.Step) stringset.Set {
	steps := stringset.New(0)
	for _, step := range s {
		steps.Add(step.Name)
	}

	return steps
}

func builderIDToStr(id *bbpb.BuilderID) string {
	byts, err := json.Marshal(id)
	if err != nil {
		panic(fmt.Sprintf("could not marshal builder ID %+v: %+v", id, err))
	}
	return string(byts)
}

func strToBuilderID(s string) *bbpb.BuilderID {
	ret := &bbpb.BuilderID{}
	if err := json.Unmarshal([]byte(s), ret); err != nil {
		panic(fmt.Sprintf("could not unmarshal builderID %q: %+v", s, err))
	}
	return ret
}

// BuildBucketAlerts returns alertable build failures generated from the given buildbucket builder.
// Invariants:
// - If the latest build succeeded, the list returned is empty for that builder.
// - Only step failures that exist in the latest build may be grouped into a range of
//   failures across builds.  Put another way: only step failures that occurred in the
//   latest build will appear anywhere in the alerts returned for that builder.
// - [TODO: diagram of some kind for the above]
func (a *Analyzer) BuildBucketAlerts(ctx context.Context, builderIDs []*bbpb.BuilderID) ([]messages.BuildFailure, error) {
	allRecentBuilds, err := a.BuildBucket.LatestBuilds(ctx, builderIDs)
	if err != nil {
		return nil, err
	}
	if len(allRecentBuilds) == 0 {
		return nil, fmt.Errorf("no recent builds from %+v", builderIDs)
	}

	// First get all of the failures that have been occurring so far for each
	// builder. This may be an empty list if there are no failures in the most
	// recent build.
	buildsByBuilderID := map[string][]*bbpb.Build{}
	for _, build := range allRecentBuilds {
		builderKey := builderIDToStr(build.Builder)
		if _, ok := buildsByBuilderID[builderKey]; !ok {
			buildsByBuilderID[builderKey] = []*bbpb.Build{}
		}
		buildsByBuilderID[builderKey] = append(buildsByBuilderID[builderKey], build)
	}

	ret := []messages.BuildFailure{}

	// Build up a list of build ranges per builder for each alertable step.
	// The AlertedBuilder type represents a build range.
	alertedBuildersByStep := map[string][]messages.AlertedBuilder{}
	for builderKey, recentBuilds := range buildsByBuilderID {
		// Assumed: recentBuilds are sorted most recent first.  TODO: Verify this is true.
		builderID := strToBuilderID(builderKey)
		latestBuild := recentBuilds[0]
		if latestBuild.Status == bbpb.Status_SUCCESS {
			continue
		}
		latestStepFailures := stepSet(alertableStepFailures(latestBuild))

		// buildsByFailingStep contains one key for each failing step in latestBuild.
		// values are slices of Build records representing continuous build runs of that failure,
		// starting at latestBuild and ending at the build where the failure first appeared.
		buildsByFailingStep := map[string][]*bbpb.Build{}

		// Now scan through earilier builds looking for the first instances of each step failure identified.
		// Note that runs of step failures may begin at different times.
		for _, build := range recentBuilds {
			allAttemptedSteps := stepSet(build.Steps)
			stepFailures := stepSet(alertableStepFailures(build))
			// Now do some set calculations:
			// - step failures that exist in latestStepFailures but not in stepFailures:
			//   The previously examined build is where that step failure started. Stop looking
			//   for this step failure in subsequent iterations.
			// - step failures that exist in latestStepFailures and also in stepFailures:
			//   Include this build in the run of builders for each of these failures.
			terminatingFailures := latestStepFailures.Difference(stepFailures)

			// Don't terminate a failure unless it actually executed in this build.
			// For example, if a test fails in build X, and in build X-1 a requisite build step
			// fails and prevents tests from running, then we haven't learned anything about
			// whether the tests are still failing. If in build X-2 the tests are still failing,
			// we want that to be part of the same run of failures that started in build X.
			terminatingFailures = terminatingFailures.Intersect(allAttemptedSteps)

			// Remove terminatingFailures from latestStepFailures.  We don't want to keep
			// looking for them in subsequent iterations of this loop.
			latestStepFailures = latestStepFailures.Difference(terminatingFailures)

			// Any failures in this build that were failing in the last examined build
			// will continue the run.
			continuingFailures := latestStepFailures.Intersect(stepFailures)

			for stepFailure := range continuingFailures {
				// Append this build to the runs of builds failing on stepFailure.
				if _, ok := buildsByFailingStep[stepFailure]; !ok {
					buildsByFailingStep[stepFailure] = []*bbpb.Build{}
				}
				buildsByFailingStep[stepFailure] = append(buildsByFailingStep[stepFailure], build)
			}
		}

		for stepName, builds := range buildsByFailingStep {
			// check ret first to see if there's already a build failure for this step
			// on some other builder. If so, just append this builder to it.
			firstFailure, latestFailure := builds[len(builds)-1], builds[0]
			if _, ok := alertedBuildersByStep[stepName]; !ok {
				alertedBuildersByStep[stepName] = []messages.AlertedBuilder{}
			}
			alertedBuilder := messages.AlertedBuilder{
				Project: builderID.Project,
				Bucket:  builderID.Bucket,
				Name:    builderID.Builder,
				URL:     fmt.Sprintf("https://ci.chromium.org/p/%s/builders/%s/%s", builderID.Project, builderID.Bucket, builderID.Builder),
				// TODO: add more buildbucket specifics to the AlertedBuilder type.
				FirstFailure:  int64(firstFailure.Number),
				LatestFailure: int64(latestFailure.Number),
				StartTime:     messages.TimeToEpochTime(time.Unix(firstFailure.StartTime.GetSeconds(), int64(firstFailure.StartTime.GetNanos()))),
				// Commit positions etc...
			}
			alertedBuildersByStep[stepName] = append(alertedBuildersByStep[stepName], alertedBuilder)
			logging.Debugf(ctx, "should merge %d failures into an alert for step: %q", len(builds), stepName)
		}
	}

	// Now group up the alerted builder ranges into individual alerts, one per step.
	// Each will contain the list of builder ranges where the step has been failing.
	for stepName, alertedBuilders := range alertedBuildersByStep {
		ret = append(ret, messages.BuildFailure{
			StepAtFault: &messages.BuildStep{
				Step: &messages.Step{
					Name: stepName,
				},
			},
			Builders: alertedBuilders,
		})
	}

	return ret, err
}

func alertableStepFailures(build *bbpb.Build) []*bbpb.Step {
	ret := []*bbpb.Step{}
	for _, buildStep := range build.Steps {
		if buildStep.Status != bbpb.Status_SUCCESS {
			ret = append(ret, buildStep)
		}
	}
	return ret
}