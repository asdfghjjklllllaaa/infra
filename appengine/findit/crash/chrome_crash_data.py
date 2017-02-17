# Copyright 2017 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import logging
from collections import namedtuple

from crash import detect_regression_range
from crash.crash_data import CrashData
from crash.chromecrash_parser import ChromeCrashParser
from crash.stacktrace import Stacktrace


class ChromeCrashData(CrashData):
  """Chrome crash report from Cracas/Fracas.

  Properties:
    identifiers (dict): The key value pairs to uniquely identify a
      ``CrashData``.
    crashed_version (str): The version of project in which the crash occurred.
    signature (str): The signature of the crash.
    platform (str): The platform name; e.g., 'win', 'mac', 'linux', 'android',
      'ios', etc.
    stacktrace (Stacktrace): The stacktrace of the crash. N.B., this is
      an object generated by parsing the string containing the stack trace;
      we do not store the string itself.
    regression_range (pair or None): a pair of the last-good and first-bad
      versions. N.B., because this is an input, it is up to clients
      to call ``DetectRegressionRange`` (or whatever else) in order to
      provide this information. In addition, while this class does
      support storing ``None`` to indicate a missing regression range
      (because the ClusterFuzz client wants that feature), the
      CL-classifier doesn't actually support that so you won't get a
      very good Culprit. The Component- and project-classifiers do still
      return some results at least.
    dependencies (dict): A dict from dependency paths to
      ``Dependency`` objects. The keys are all those deps which are
      used by both the ``crashed_version`` of the code, and at least
      one frame in the ``stacktrace.crash_stack``.
    dependency_rolls (dict) A dict from dependency
      paths to ``DependencyRoll`` objects. The keys are all those
      dependencies which (1) occur in the regression range for the
      ``platform`` where the crash occurred, (2) neither add nor delete
      a dependency, and (3) are also keys of ``dependencies``.
  """

  def __init__(self, crash_data, dep_fetcher, top_n_frames=None):
    """
    Args:
      crash_data (dict): Dicts sent through Pub/Sub by Cracas/Fracas. Example:
      {
          'stack_trace': 'CRASHED [0x43507378...',
          # The Chrome version that produced the stack trace above.
          'chrome_version': '52.0.2743.41',
          # Client could provide customized data.
          'customized_data': {
              'trend_type': 'd',  # see supported types below
              'channel': 'beta',
              # Historical data about crash per million pageload by Chrome
              # version. (Right now last 20 versions)
              'historical_metadata': [
                  {
                      'report_number': 0,
                      'cpm': 0.0,
                      'client_number': 0,
                      'chrome_version': '51.0.2704.103'
                  },
                  ...
                  {
                      'report_number': 10,
                      'cpm': 2.1,
                      'client_number': 8,
                      'chrome_version': '53.0.2768.0'
                  },
              ]
          },
          'platform': 'mac',    # On which platform the crash occurs.
          'client_id': 'fracas',   # Identify which client this request is from.
          'signature': '[ThreadWatcher UI hang] base::RunLoopBase::Run',
          'crash_identifiers': {    # A list of key-value to identify a crash.
              'platform': 'mac',
              'version': '52.0.2743.41',
              'process_type': 'browser',
              'channel': 'beta',
              # Signature for the stack trace.
              'signature': '[ThreadWatcher UI hang] base::RunLoopBase::Run'
          }
      }
      dep_fetcher (ChromeDependencyFetcher): Dependency fetcher that can fetch
        all dependencies related to crashed version.
      top_n_frames (int): number of the frames in stacktrace we should parse.
    """
    super(ChromeCrashData, self).__init__(crash_data)
    self._channel = crash_data['customized_data']['channel']
    self._historical_metadata = crash_data['customized_data'][
        'historical_metadata']

    # Delay the stacktrace parsing to the first time when stacktrace property
    # gets called.
    self._stacktrace_str = crash_data['stack_trace'] or ''
    self._top_n_frames = top_n_frames
    self._stacktrace = None

    self._dep_fetcher = dep_fetcher
    self._crashed_version_deps = None

    self._regression_range = None

    self._dependencies = {}
    self._dependency_rolls = {}

  def _CrashedVersionDeps(self):
    """Gets all dependencies related to crashed_version.

    N.B. All dependencies will be returned, no matter whether they appeared in
    stacktrace or are related to the crash or not.
    """
    if self._crashed_version_deps:
      return self._crashed_version_deps

    self._crashed_version_deps = self._dep_fetcher.GetDependency(
        self.crashed_version, self.platform) if self._dep_fetcher else {}

    return self._crashed_version_deps

  @property
  def channel(self):
    return self._channel

  @property
  def historical_metadata(self):
    return self._historical_metadata

  @property
  def stacktrace(self):
    """Parses stacktrace and returns parsed ``Stacktrace`` object."""
    if self._stacktrace:
      return self._stacktrace

    self._stacktrace = ChromeCrashParser().Parse(
        self._stacktrace_str, self._CrashedVersionDeps(),
        signature=self.signature, top_n_frames=self._top_n_frames)
    if not self._stacktrace:
      logging.warning('Failed to parse the stacktrace %s',
                      self._stacktrace_str)
    return self._stacktrace

  @property
  def regression_range(self):
    """Detects regression range from ``historical_metadata`` and returns it."""
    if self._regression_range:
      return self._regression_range

    regression_range = detect_regression_range.DetectRegressionRange(
        self.historical_metadata)
    if regression_range is None: # pragma: no cover
      logging.warning('Got ``None`` for the regression range.')
    else:
      self._regression_range = tuple(regression_range)

    return self._regression_range

  @property
  def dependencies(self):
    """Get all dependencies that are in the crash stack of stacktrace."""
    if self._dependencies:
      return self._dependencies

    if not self.stacktrace:
      logging.warning('Cannot get depenencies without stacktrace.')
      return {}

    self._dependencies = {
        frame.dep_path: self._CrashedVersionDeps()[frame.dep_path]
        for frame in self.stacktrace.crash_stack.frames
        if frame.dep_path and frame.dep_path in self._CrashedVersionDeps()
    }
    return self._dependencies

  @property
  def dependency_rolls(self):
    """Gets all dependency rolls of ``dependencies`` in regression range."""
    if self._dependency_rolls:
      return self._dependency_rolls

    # Short-circuit when we know the deprolls must be empty.
    if not self.regression_range or not self.stacktrace:
      logging.warning('Cannot get deps and dep rolls for report without '
                      'regression range or stacktrace.')
      return {}

    # Get ``DependencyRoll` objects for all dependencies in the regression
    # range (for the particular platform that crashed).
    regression_range_dep_rolls = self._dep_fetcher.GetDependencyRollsDict(
        self.regression_range[0], self.regression_range[1], self.platform)
    # Filter out the ones which add or delete a dependency, because we
    # can't really be sure whether to blame them or not. This rarely
    # happens, so our inability to decide shouldn't be too much of a problem.
    def HasBothRevisions(dep_path, dep_roll):
      has_both_revisions = bool(dep_roll.old_revision) and bool(
          dep_roll.new_revision)
      if not has_both_revisions:
        logging.info(
            'Skip %s dependency %s',
            'added' if dep_roll.new_revision else 'deleted',
            dep_path)
      return has_both_revisions

    # Apply the above filter, and also filter to only retain those
    # which occur in ``crashed_stack_deps``.
    self._dependency_rolls = {
        dep_path: dep_roll
        for dep_path, dep_roll in regression_range_dep_rolls.iteritems()
        if HasBothRevisions(dep_path, dep_roll) and dep_path
        in self.dependencies
    }

    return self._dependency_rolls
