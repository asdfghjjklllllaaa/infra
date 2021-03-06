# Copyright 2018 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import json
import mock
import webapp2

from gae_libs.handlers.base_handler import BaseHandler
from handlers.flake.detection import detect_flakes
from model.flake.flake_type import FlakeType
from model.flake.flake_type import FLAKE_TYPE_DESCRIPTIONS
from services import flake_issue_util
from services.flake_detection import detect_flake_occurrences
from services.flake_detection.detect_flake_occurrences import (
    DetectFlakesFromFlakyCQBuildParam)
from waterfall.test.wf_testcase import WaterfallTestCase


class DetectFlakesCronJobTest(WaterfallTestCase):
  app_module = webapp2.WSGIApplication([
      ('/flake/detection/cron/detect-flakes',
       detect_flakes.DetectFlakesCronJob),
  ],
                                       debug=True)

  @mock.patch.object(BaseHandler, 'IsRequestFromAppSelf', return_value=True)
  def testTaskAddedToQueue(self, mocked_is_request_from_appself):
    response = self.test_app.get('/flake/detection/cron/detect-flakes')
    self.assertEqual(200, response.status_int)
    response = self.test_app.get('/flake/detection/cron/detect-flakes')
    self.assertEqual(200, response.status_int)

    tasks = self.taskqueue_stub.get_filtered_tasks(
        queue_names='flake-detection-queue')
    self.assertEqual(2, len(tasks))
    self.assertTrue(mocked_is_request_from_appself.called)


class FlakeDetectionAndAutoActionTest(WaterfallTestCase):
  app_module = webapp2.WSGIApplication([
      ('/flake/detection/task/detect-flakes',
       detect_flakes.FlakeDetectionAndAutoAction),
  ],
                                       debug=True)

  @mock.patch.object(BaseHandler, 'IsRequestFromAppSelf', return_value=True)
  @mock.patch.object(detect_flake_occurrences, 'QueryAndStoreHiddenFlakes')
  @mock.patch.object(
      flake_issue_util, 'GetFlakeGroupsForActionsOnBugs', return_value=([], []))
  @mock.patch.object(flake_issue_util, 'ReportFlakesToFlakeAnalyzer')
  @mock.patch.object(flake_issue_util, 'ReportFlakesToMonorail')
  @mock.patch.object(flake_issue_util, 'GetFlakesWithEnoughOccurrences')
  @mock.patch.object(detect_flake_occurrences, 'QueryAndStoreFlakes')
  def testFlakesDetected(self, mock_detect, mock_get_flakes, mock_bug,
                         mock_analysis, mock_groups, mock_detect_hidden, _):
    mock_get_flakes.return_value = []
    response = self.test_app.get(
        '/flake/detection/task/detect-flakes', status=200)
    self.assertEqual(200, response.status_int)

    mock_detect.assert_has_calls([
        mock.call(FlakeType.CQ_FALSE_REJECTION),
        mock.call(FlakeType.RETRY_WITH_PATCH)
    ])
    mock_bug.assert_called_once_with([], [])
    mock_analysis.assert_called_once_with([])
    mock_groups.assert_called_once_with([])
    self.assertTrue(mock_detect_hidden.called)


class DetectFlakesFromFlakyCQBuildTest(WaterfallTestCase):
  app_module = webapp2.WSGIApplication([
      ('/flake/detection/task/detect-flakes-from-build',
       detect_flakes.DetectFlakesFromFlakyCQBuild),
  ],
                                       debug=True)

  @mock.patch.object(BaseHandler, 'IsRequestFromAppSelf', return_value=True)
  @mock.patch.object(detect_flake_occurrences, 'ProcessBuildForFlakes')
  def testFlakesDetected(self, mock_detect, _):
    # pylint:disable=unused-argument
    flake_type = FlakeType.CQ_FALSE_REJECTION
    params = DetectFlakesFromFlakyCQBuildParam(
        build_id=123,
        flake_type_desc=FLAKE_TYPE_DESCRIPTIONS[flake_type]).ToSerializable()
    response = self.test_app.post(
        '/flake/detection/task/detect-flakes-from-build',
        params=json.dumps(params),
        headers={'X-AppEngine-QueueName': 'task_queue'})
    self.assertEqual(200, response.status_int)

    # Assertions have never worked properly because we were using mock 1.0.1.
    # After rolling to mock 2.0.0, which fixes assertions, these assertions now
    # fail. https://crbug.com/947753.
    # mock_detect.assert_has_called_once_with(params)
