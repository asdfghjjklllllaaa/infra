# Copyright 2018 The Chromium Authors. All rights reserved.
# Use of this source code is govered by a BSD-style
# license that can be found in the LICENSE file or at
# https://developers.google.com/open-source/licenses/bsd

"""Tests for the users servicer."""

import unittest

import mox
from components.prpc import codes
from components.prpc import context
from components.prpc import server

from api import users_servicer
from api import monorailcontext
from api.api_proto import users_pb2
from framework import authdata
from testing import fake
from services import service_manager


class UsersServicerTest(unittest.TestCase):

  def setUp(self):
    self.mox = mox.Mox()
    self.cnxn = fake.MonorailConnection()
    self.services = service_manager.Services(
        config=fake.ConfigService(),
        issue=fake.IssueService(),
        user=fake.UserService(),
        usergroup=fake.UserGroupService(),
        project=fake.ProjectService(),
        features=fake.FeaturesService())
    self.project = self.services.project.TestAddProject('proj', project_id=987)
    self.user = self.services.user.TestAddUser('owner@example.com', 111L)
    self.users_svcr = users_servicer.UsersServicer(
        self.services, make_rate_limiter=False)
    self.prpc_context = context.ServicerContext()
    self.prpc_context.set_code(codes.StatusCode.OK)

  def tearDown(self):
    self.mox.UnsetStubs()
    self.mox.ResetAll()

  def testGetUsers(self):
    """API call to GetUsers can reach the Do* method."""
    self.assertIsNone(self.users_svcr.rate_limiter)
    request = users_pb2.GetUserRequest(email='test@example.com')
    response = self.users_svcr.GetUser(
        request, self.prpc_context, cnxn=self.cnxn,
        auth=authdata.AuthData(user_id=111L, email='owner@example.com'))
    self.assertEqual(codes.StatusCode.OK, self.prpc_context._code)
    self.assertEqual(hash('test@example.com'), response.id)

  def testDoGetUsers_Normal(self):
    """We can get a user by email address."""
    request = users_pb2.GetUserRequest(email='test@example.com')
    mc = monorailcontext.MonorailContext(
        self.services, cnxn=self.cnxn, requester='owner@example.com')
    response = self.users_svcr.DoGetUser(mc, request)
    self.assertEqual(hash('test@example.com'), response.id)
