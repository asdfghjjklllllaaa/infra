# Copyright 2019 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import logging

from buildbucket_proto import common_pb2

from findit_v2.services import projects
from findit_v2.services.failure_type import StepTypeEnum


def OnBuildFailure(context, build):
  """Processes the failed build within the given context.

  Args:
    context (findit_v2.services.context.Context): Scope of the analysis.
    build (buildbucket build.proto): ALL info about the build.

  Returns:
    True if the failed build is supported and analyzed; otherwise False.
  """
  logging.info('Context of analysis: %r', context)
  logging.info('Failed build is: %r', build.id)

  failed_steps = []
  project_api = projects.GERRIT_PROJECTS[context
                                         .luci_project_name]['project-api']
  for step in build.steps:
    if step.status != common_pb2.FAILURE:
      continue
    failure_type = project_api.ClassifyStepType(step)
    failed_steps.append((step, failure_type))

  if not failed_steps:
    return False

  compile_steps = [
      fs[0] for fs in failed_steps if fs[1] == StepTypeEnum.COMPILE
  ]
  if compile_steps:
    _AnalyzeCompileFailure(context, build, compile_steps)
    return True

  logging.info('Unsupported failure types: %r', [fs[1] for fs in failed_steps])
  return False


def _AnalyzeCompileFailure(context, build, compile_steps):
  """Analyzes the compile failure.

  Args:
    context (findit_v2.services.context.Context): Scope of the analysis.
    build (buildbucket build.proto): ALL info about the build.
    compile_steps (buildbucket step.proto): The failed compile steps.
  """
  #pylint: disable=unused-argument
  # TODO(chanli): hook up regression detection and analysis.
  pass