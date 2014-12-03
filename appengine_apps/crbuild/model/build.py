# Copyright 2014 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

from datetime import datetime, timedelta

from google.appengine.ext import ndb
from google.appengine.ext.ndb import msgprop
from protorpc import messages


MAX_LEASE_BUILDS = 100
MAX_LEASE_SECONDS = 3 * 60 * 60  # 3 hours


class BuildStatus(messages.Enum):
  SCHEDULED = 1
  BUILDING = 2
  SUCCESS = 3
  FAILURE = 4
  EXCEPTION = 5  # Infrastructure failure.


LEASABLE_STATUSES = (BuildStatus.SCHEDULED, BuildStatus.BUILDING)


class BuildProperties(ndb.Expando):
  """Arbitrary build properties, a key-value pair map.

  Attributes:
    builder_name (string): special property, that is used for build grouping
      during rendering.
  """
  builder_name = ndb.StringProperty()


class Build(ndb.Model):
  """Describes a build.

  Build keys are autogenerated integers.
  The only Build entity in its entity group.

  Attributes:
    namespace (string): a generic way to distinguish builds. Different build
      namespaces have different permissions.
    properties (BuildProperties): key-value pair map.
      builder_name: a special property used for grouping builds when rendering.
    status (BuildStatus): status of the build
    url (str): a URL to a build-system-specific build, viewable by a human.
    available_since (datetime): the earliest time the build can be leased.
      The moment the build is leased, |available_since| is set to
      (utcnow + lease_duration). On build creation, is set to utcnow.
  """

  namespace = ndb.StringProperty(required=True)
  properties = ndb.StructuredProperty(BuildProperties)
  status = msgprop.EnumProperty(BuildStatus, default=BuildStatus.SCHEDULED)
  url = ndb.StringProperty(indexed=False)
  available_since = ndb.DateTimeProperty(required=True, auto_now_add=True)

  @property
  def builder_name(self):
    return self.properties.builder_name if self.properties else None

  def set_status(self, value):
    """Changes build status and notifies interested parties."""
    if self.status == value:
      return
    self.status = value
    # TODO(nodir): uncomment when model/log.py is added
    # if value == BuildStatus.BUILDING:
    #   BuildStarted().add_to(self)
    # elif value in (BuildStatus.SUCCESS, BuildStatus.FAILURE):
    #  BuildCompleted().add_to(self)

  def is_available(self):
    return self.available_since <= datetime.utcnow()

  def modify_lease(self, lease_seconds):
    """Changes build's lease, updates |available_since|."""
    self.available_since = datetime.utcnow() + timedelta(seconds=lease_seconds)

  @property
  def key_string(self):
    """Returns an opaque key string."""
    return self.key.urlsafe() if self.key else None

  @classmethod
  def parse_key_string(cls, key_string):
    """Parses an opaque key string."""
    key = ndb.Key(urlsafe=key_string)
    assert key.kind() == cls.__name__
    return key

  @classmethod
  def lease(cls, namespaces, lease_seconds=10, max_builds=10):
    """Leases builds.

    Builds are sorted by available_since attribute, oldest first.

    Args:
      lease_seconds (int): lease duration. After lease expires, the Build can be
        leased again.
      max_builds (int): maximum number of builds to return.
      namespaces (list of string): lease only builds with any of |namespaces|.

    Returns:
      A list of Builds.
    """
    assert isinstance(namespaces, list)
    assert namespaces, 'No namespaces specified'
    assert all(isinstance(n, basestring) for n in namespaces), (
        'namespaces must be strings'
    )
    assert isinstance(lease_seconds, (int, float))
    assert lease_seconds < MAX_LEASE_SECONDS, (
        'lease_seconds must not exceed %d' % MAX_LEASE_BUILDS
    )
    assert isinstance(max_builds, int)
    assert max_builds <= MAX_LEASE_BUILDS, (
        'max_builds must not be greater than %s' % MAX_LEASE_BUILDS
    )

    now = datetime.utcnow()
    q = cls.query(
        cls.status.IN([BuildStatus.SCHEDULED, BuildStatus.BUILDING]),
        cls.namespace.IN(namespaces),
        cls.available_since <= now
    )
    q = q.order(Build.available_since) # oldest first.

    new_available_since = now + timedelta(seconds=lease_seconds)

    @ndb.transactional
    def lease_build(build):
      if build.status not in LEASABLE_STATUSES:  # pragma: no cover
        return False
      if not build.is_available():  # pragma: no cover
        return False
      build.available_since = new_available_since
      build.put()
      return True

    builds = []
    # TODO(nodir): either optimize this query using memcache, or return builds
    # without leasing and then lease builds one by one.
    for b in q.fetch(max_builds):
      if lease_build(b):  # pragma: no branch
        builds.append(b)
    return builds

  @ndb.transactional
  def unlease(self):
    """Removes build lease."""
    self.available_since = datetime.utcnow()
    self.put()
