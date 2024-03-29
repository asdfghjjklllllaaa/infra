# Copyright 2019 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

from collections import defaultdict
import logging

from buildbucket_proto import common_pb2

from findit_v2.model.compile_failure import CompileFailure
from findit_v2.model.compile_failure import CompileFailureAnalysis
from findit_v2.model.compile_failure import CompileRerunBuild
from findit_v2.model.messages.findit_result import (
    BuildFailureAnalysisResponse)
from findit_v2.services import build_util
from findit_v2.services import projects
from findit_v2.services.analysis import analysis_util
from findit_v2.services.analysis.compile_failure import (
    compile_failure_rerun_analysis)
from findit_v2.services.analysis.compile_failure import pre_compile_analysis
from findit_v2.services.failure_type import StepTypeEnum


def AnalyzeCompileFailure(context, build, compile_steps):
  """Analyzes compile failure from a failed ci/postsubmit build.

  Args:
    context (findit_v2.services.context.Context): Scope of the analysis.
    build (buildbucket build.proto): ALL info about the build.
    compile_steps (list of buildbucket step.proto): The failed compile steps.

  Returns:
    (bool): Returns True if a new analysis starts, otherwise False.
  """
  luci_project = context.luci_project_name
  if luci_project == 'chromium':
    logging.warning('Findit does not support chromium project in v2.')
    return False

  project_api = projects.GetProjectAPI(luci_project)
  if not project_api:
    logging.debug('Unsupported project %s', luci_project)
    return False

  # Project config for if failures should be grouped to reduce duplicated
  # analyses.
  should_group_failures = projects.PROJECT_CFG.get(
      luci_project, {}).get('should_group_failures')

  detailed_compile_failures = project_api.GetCompileFailures(
      build, compile_steps)
  # Checks previous builds to look for first time failures for all the failures
  # in current failed build.
  pre_compile_analysis.UpdateCompileFailuresWithFirstFailureInfo(
      context, build, detailed_compile_failures)
  pre_compile_analysis.SaveCompileFailures(context, build,
                                           detailed_compile_failures)

  # Looks for the failures that started to fail in the current build.
  first_failures_in_current_build = (
      pre_compile_analysis.GetFirstFailuresInCurrentBuild(
          context, build, detailed_compile_failures))
  if not first_failures_in_current_build.get('failures'):
    logging.info(
        'No new analysis for build %d because all failures have '
        'happened in previous builds.', build.id)
    return False

  # Filters out the first failures with existing failure group.
  if should_group_failures:
    failures_without_existing_group = (
        pre_compile_analysis.GetFirstFailuresInCurrentBuildWithoutGroup(
            context, build, first_failures_in_current_build))
  else:
    failures_without_existing_group = first_failures_in_current_build

  if not failures_without_existing_group.get('failures'):
    logging.info(
        'All failures have matching failure groups in build %s,'
        ' no need to start a new analysis.', build.id)
    return False

  # Start a new analysis to analyze the first time failures.
  pre_compile_analysis.SaveCompileAnalysis(
      context, build, failures_without_existing_group, should_group_failures)
  compile_failure_rerun_analysis.RerunBasedAnalysis(context, build.id, build)
  return True


def _ProcessAndSaveRerunBuildResult(context, analyzed_build_id, rerun_build):
  """Gets results of a completed rerun build and save it in datastore.

  Args:
     context (findit_v2.services.context.Context): Scope of the analysis.
     analyzed_build_id (int): Id of the failed ci/post_submit build that's being
       analyzed.
     rerun_build (buildbucket build.proto): ALL info about the rerun build.

   Returns:
     True if the rerun build entity is updated, otherwise False.
  """
  project_api = projects.GetProjectAPI(context.luci_project_name)
  assert project_api, 'Unsupported project %s.' % context.luci_project_name

  analysis = CompileFailureAnalysis.GetVersion(analyzed_build_id)
  if not analysis:
    logging.error('CompileFailureAnalysis missing for %d.', analyzed_build_id)
    return False

  rerun_build_entity = CompileRerunBuild.get_by_id(
      rerun_build.id, parent=analysis.key)
  if not rerun_build_entity:
    logging.error('CompileRerunBuild entity for build %d missing.',
                  rerun_build.id)
    return False

  detailed_compile_failures = {}
  if rerun_build.status == common_pb2.FAILURE:
    failed_steps = build_util.GetFailedStepsInBuild(context, rerun_build)
    compile_steps = [
        fs[0] for fs in failed_steps if fs[1] == StepTypeEnum.COMPILE
    ]
    detailed_compile_failures = project_api.GetCompileFailures(
        rerun_build, compile_steps) if compile_steps else {}

  rerun_build_entity.SaveRerunBuildResults(rerun_build.status,
                                           detailed_compile_failures)
  return True


def OnCompileRerunBuildCompletion(context, rerun_build):
  """Processes the completed rerun build then continues the analysis.

  Processes and stores results of the completed rerun build. And then resume the
  analysis to either start the next rerun build or identify the culprit.

  Args:
    context (findit_v2.services.context.Context): Scope of the analysis.
    rerun_build (buildbucket build.proto): ALL info about the rerun build.

  Returns:
    True if the rerun build entity is updated, otherwise False.
  """
  analyzed_build_id = build_util.GetAnalyzedBuildIdFromRerunBuild(rerun_build)
  if not analyzed_build_id:
    logging.error('analyzed_build_id not set in the rerun build %d.',
                  rerun_build.id)
    return False

  build_saved = _ProcessAndSaveRerunBuildResult(context, analyzed_build_id,
                                                rerun_build)
  if not build_saved:
    return False

  compile_failure_rerun_analysis.RerunBasedAnalysis(context, analyzed_build_id,
                                                    rerun_build)
  return True


def _IsAnalysisFinished(failures):
  """ Checks if the analysis has finished.

  Here the 'analysis' refers to the overall analysis for the failures, it's
  possible that there are multiple analyses analyzing different failures, but
  Findit will just summarizes all and return a single flag.

  Returns:
    True if analyses exist and have completed, otherwise False.
  """
  analysis_build_ids = set([failure.build_id for failure in failures])
  for analysis_build_id in analysis_build_ids:
    analysis = CompileFailureAnalysis.GetVersion(analysis_build_id)
    if not analysis:
      logging.warning('Not found compile analysis for build %d',
                      analysis_build_id)
      return False

    if not analysis.completed:
      return False

  return True


def OnCompileFailureAnalysisResultRequested(request, requested_build):
  """Returns the findings for the requested build's compile failure.

  Since SoM doesn't have atomic failure info for compile steps, currently Findit
  will only respond with aggregated step level results.

  Args:
    request(findit_result.BuildFailureAnalysisRequest): request for a build
      failure.
    requested_build(LuciFailedBuild): A LuciFailedBuild entity with COMPILE
      build_failure_type.

  Returns:
    [findit_result.BuildFailureAnalysisResponse]: Analysis results
      for the requested build.
  """
  compile_failures = CompileFailure.query(ancestor=requested_build.key).fetch()
  if not compile_failures:
    return None

  requested_steps = request.failed_steps
  requested_failures = defaultdict(list)

  for failure in compile_failures:
    if requested_steps and failure.step_ui_name not in requested_steps:
      continue
    requested_failures[failure.step_ui_name].append(failure)

  responses = []
  for step_ui_name, requested_failures_in_step in requested_failures.iteritems(
  ):
    merged_failures = []
    for failure in requested_failures_in_step:
      # Merged failures are the failures being actually analyzed and only they
      # have stored culprits info.
      merged_failures.append(failure.GetMergedFailure())

    culprits = analysis_util.GetCulpritsForFailures(merged_failures)
    response = BuildFailureAnalysisResponse(
        build_id=request.build_id,
        build_alternative_id=request.build_alternative_id,
        step_name=step_ui_name,
        test_name=None,
        culprits=culprits,
        is_finished=_IsAnalysisFinished(merged_failures),
        is_supported=True,
    )
    responses.append(response)

  return responses
