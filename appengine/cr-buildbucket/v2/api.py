# Copyright 2018 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import functools

from google.appengine.ext import ndb

from components import auth
from components import prpc

# Some of these imports are required to populate proto symbol db.
from proto import rpc_pb2  # pylint: disable=unused-import
from proto import rpc_prpc_pb2
from proto import step_pb2  # pylint: disable=unused-import

import buildtags
import model
import service
import v2


class StatusCodeError(Exception):

  def __init__(self, code, details):
    super(StatusCodeError, self).__init__('%s: %s' % (code, details))
    self.code = code
    self.details = details


def status_code_error_class(code):
  code_name = code[1]

  class Error(StatusCodeError):

    def __init__(self, details=code_name):
      super(Error, self).__init__(code, details)

  return Error


NotFound = status_code_error_class(prpc.StatusCode.NOT_FOUND)
InvalidArgument = status_code_error_class(prpc.StatusCode.INVALID_ARGUMENT)


def api_method(fn):
  """Decorates an RPC implementation with error handling."""

  @functools.wraps(fn)
  def decorated(self, req, context):
    try:
      try:
        return fn(self, req, context)
      except auth.AuthorizationError:
        raise NotFound()
    except StatusCodeError as ex:
      context.set_code(ex.code)
      context.set_details(ex.details)
      return None

  return decorated


def v1_bucket(builder_id):
  return 'luci.%s.%s' % (builder_id.project, builder_id.bucket)


def to_build_messages(builds, field_mask):
  """Converts model.Build instances to build_pb2.Build messages."""
  builds_msgs = map(v2.build_to_v2_partial, builds)

  # TODO(nodir): load steps conditionally, according to the field_mask.
  annotations = ndb.get_multi(
      [model.BuildAnnotations.key_for(b.key) for b in builds])
  for b, build_ann in zip(builds_msgs, annotations):
    if build_ann:
      b.steps.extend(v2.parse_steps(build_ann))

  # TODO(nodir): apply field_mask.
  return builds_msgs


class BuildsApi(object):
  """Implements buildbucket.v2.Builds proto service."""

  DESCRIPTION = rpc_prpc_pb2.BuildsServiceDescription

  @api_method
  def GetBuild(self, req, ctx):
    """Retrieves a build by id or number."""
    if req.id:
      if req.HasField('builder') or req.build_number:
        raise InvalidArgument(
            'id is mutually exclusive with builder and build_number')
      build = service.get(req.id)
    elif req.HasField('builder') and req.build_number:
      bucket = v1_bucket(req.builder)
      tag = buildtags.build_address_tag(bucket, req.builder.builder,
                                        req.build_number)
      builds, _ = service.search(
          service.SearchQuery(
              buckets=[bucket],
              tags=[tag],
          ))
      build = builds[0] if builds else None
    else:
      raise InvalidArgument('id or (builder and build_number) are required')

    if not build:
      raise NotFound()
    return to_build_messages([build], req.fields)[0]
