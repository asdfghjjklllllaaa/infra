# Copyright 2018 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import datetime
import json
import mock
import webapp2

from handlers.flake.detection import rank_flakes
from model.flake.flake import Flake
from model.flake.flake_issue import FlakeIssue
from waterfall.test.wf_testcase import WaterfallTestCase


class RankFlakesTest(WaterfallTestCase):
  app_module = webapp2.WSGIApplication(
      [
          ('/ranked-flakes', rank_flakes.RankFlakes),
      ], debug=True)

  def setUp(self):
    super(RankFlakesTest, self).setUp()
    self.flake_issue = FlakeIssue.Create(
        monorail_project='chromium', issue_id=900)
    self.flake_issue.last_updated_time = datetime.datetime(2018, 1, 1)
    self.flake_issue.put()

    luci_project = 'chromium'
    normalized_step_name = 'normalized_step_name'
    normalized_test_name = 'normalized_test_name'
    self.flake1 = Flake.Create(
        luci_project=luci_project,
        normalized_step_name=normalized_step_name,
        normalized_test_name=normalized_test_name,
        test_label_name=normalized_test_name)

    self.flake1.flake_issue_key = self.flake_issue.key
    self.flake1.false_rejection_count_last_week = 2
    self.flake1.impacted_cl_count_last_week = 2
    self.flake1.put()

    self.flake2 = Flake.Create(
        luci_project=luci_project,
        normalized_step_name=normalized_step_name,
        normalized_test_name='another_test',
        test_label_name='another_test')
    self.flake2.put()

    self.flake3 = Flake.Create(
        luci_project=luci_project,
        normalized_step_name=normalized_step_name,
        normalized_test_name='another_test',
        test_label_name='another_test')
    self.flake3.false_rejection_count_last_week = 5
    self.flake3.impacted_cl_count_last_week = 2
    self.flake3.put()

    self.flake1_dict = self.flake1.to_dict()
    self.flake1_dict['flake_urlsafe_key'] = self.flake1.key.urlsafe()
    self.flake1_dict['flake_issue'] = self.flake_issue.to_dict()
    self.flake1_dict['flake_issue']['issue_link'] = FlakeIssue.GetLinkForIssue(
        self.flake_issue.monorail_project, self.flake_issue.issue_id)

    self.flake3_dict = self.flake3.to_dict()
    self.flake3_dict['flake_urlsafe_key'] = self.flake3.key.urlsafe()

  def testRankFlakes(self):
    response = self.test_app.get(
        '/ranked-flakes', params={
            'format': 'json',
        }, status=200)

    self.assertEqual(
        json.dumps(
            {
                'flakes_data': [self.flake3_dict, self.flake1_dict],
                'prev_cursor': '',
                'cursor': '',
                'n': '',
                'luci_project': '',
                'test_name_filter': ''
            },
            default=str), response.body)

  @mock.patch.object(
      Flake, 'NormalizeTestName', return_value='normalized_test_name')
  def testSearchRedirect(self, _):
    response = self.test_app.get(
        '/ranked-flakes?test_name=test_name',
        params={
            'format': 'json',
        },
        status=302)

    expected_url_suffix = (
        '/flake/occurrences?key=%s' % self.flake1.key.urlsafe())

    self.assertTrue(
        response.headers.get('Location', '').endswith(expected_url_suffix))
