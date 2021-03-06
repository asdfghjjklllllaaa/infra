# Copyright 2019 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

from google.appengine.ext import ndb


class GitilesCommit(ndb.Model):
  """A class for a gitiles commit."""

  # Gitiles hostname, e.g. 'chromium.googlesource.com'.
  gitiles_host = ndb.StringProperty(required=True)

  # Project name the commit belongs to, e.g. 'chromium/src'.
  gitiles_project = ndb.StringProperty(required=True)

  # Associated git ref of the commit, e.g. 'refs/heads/master'.
  # NOT a branch name: must start with 'refs/'.
  # If not specified when query build, default to 'refs/heads/master'.
  gitiles_ref = ndb.StringProperty(required=True)

  # SHA1 of the change.
  # May be called as 'revision' or 'git_hash' somewhere else.
  gitiles_id = ndb.StringProperty(required=True)

  # Integer identifier of a commit, required to sort commits.
  commit_position = ndb.IntegerProperty(required=True)


class Culprit(GitilesCommit):
  """Base class for a suspected or culprit commit."""
  # Urlsafe_keys to atom failures this culprit is responsible for.
  # Uses urlsafe_keys so that it can accept both compile and test failures.
  failure_urlsafe_keys = ndb.StringProperty(repeated=True)

  @classmethod
  def _CreateKey(cls, gitiles_host, gitiles_project, gitiles_ref, gitiles_id):
    return ndb.Key(
        cls.__name__, '{}/{}/{}/{}'.format(gitiles_host, gitiles_project,
                                           gitiles_ref, gitiles_id))

  @classmethod
  def Create(cls,
             gitiles_host,
             gitiles_project,
             gitiles_ref,
             gitiles_id,
             commit_position,
             failure_urlsafe_keys=None):
    """Creates an entity for a culprit."""
    return cls(
        key=cls._CreateKey(gitiles_host, gitiles_project, gitiles_ref,
                           gitiles_id),
        gitiles_host=gitiles_host,
        gitiles_project=gitiles_project,
        gitiles_ref=gitiles_ref,
        gitiles_id=gitiles_id,
        commit_position=commit_position,
        failure_urlsafe_keys=failure_urlsafe_keys or [])

  @classmethod
  def Get(cls, gitiles_host, gitiles_project, gitiles_ref, gitiles_id):
    return cls._CreateKey(gitiles_host, gitiles_project, gitiles_ref,
                          gitiles_id).get()

  @classmethod
  def GetOrCreate(cls,
                  gitiles_host,
                  gitiles_project,
                  gitiles_ref,
                  gitiles_id,
                  commit_position=None,
                  failure_urlsafe_keys=None):
    """Gets or Creates a Culprit entity.

    If failure_urlsafe_keys provided, update the culprit as well.
    """
    updated = False
    culprit = cls.Get(gitiles_host, gitiles_project, gitiles_ref, gitiles_id)
    if not culprit:
      culprit = cls.Create(gitiles_host, gitiles_project, gitiles_ref,
                           gitiles_id, commit_position)
      updated = True

    if failure_urlsafe_keys:
      culprit.failure_urlsafe_keys = list(
          set(culprit.failure_urlsafe_keys) | set(failure_urlsafe_keys))
      updated = True

    if updated:
      culprit.put()
    return culprit
