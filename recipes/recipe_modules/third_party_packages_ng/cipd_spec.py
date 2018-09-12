# Copyright 2018 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

class CIPDSpec(object):
  """CIPDSpec represents a single CIPD package.

  It represents a CIPD (pkg, symver) pair, and allows the following operations:
    * Checking its existance (either locally or on the CIPD server)
    * Fetching to a local cached location
    * Building a package into a local cached location
    * Uploading the locally cached package to the server
    * Deploying the locally cached package to disk
  """

  # (pkg, symver) -> instance_id
  _VERSION_CACHE = {}

  def __init__(self, api, pkg, symver):
    self._api = api
    self._pkg = str(pkg)
    self._symver = str(symver)
    assert self._pkg != 'None'
    assert self._symver, 'what: %r %r' % (pkg, symver)
    assert 'None' not in self._symver
    assert self._symver != 'latest'

  @property
  def tag(self):
    return 'version:'+self._symver

  @property
  def _cache_key(self):
    return (self._pkg, self._symver)

  def check(self):
    """Returns True if the package is available locally or on the server."""
    try:
      self.resolve()
      return True
    except self._api.step.StepFailure:
      self._api.step.active_result.presentation.status = (
        self._api.step.SUCCESS)
      self._api.step.active_result.presentation.step_text = (
        'tag %r not found' % (self.tag,))
      return False

  def resolve(self):
    """Returns the instance_id for this CIPDSpec.

    If this CIPDSpec has been built locally, this will return the instance_id of
    the locally built package; otherwise this will resolve the instance_id from
    the CIPD server.

    Returns str of the instance_id, as reported by CIPD.

    Raises StepFailure if we haven't built this CIPDSpec locally and the CIPD
    server doesn't have this version.
    """
    if self._cache_key not in self._VERSION_CACHE:
      # by default, make this return no tags in test mode.
      desc = self._api.cipd.describe(
        self._pkg, self.tag, test_data_tags=(), test_data_refs=())
      self._api.step.active_result.presentation.step_text = (
        'found %r' % desc.pin.instance_id)
      self._VERSION_CACHE[self._cache_key] = desc.pin.instance_id
    return self._VERSION_CACHE[self._cache_key]

  def deploy(self, root):
    """Deploys the CIPD package to disk (at the given root).

    If the package is already cached locally, this deploys from that. Otherwise
    it will fetch the package from the server.
    """
    self._api.cipd.pkg_deploy(root, self._ensure_fetched())

  @property
  def _resolved_instance_id(self):
    """The cached instance_id for this CIPDSpec.

    If `resolve` hasn't been called and if the package hasn't been built locally
    this returns None."""
    return self._VERSION_CACHE.get(self._cache_key)

  @property
  def _local_pkg_path_dir(self):
    """The local cache directory which would hold the CIPD package files for the
    CIPD package. If there are multiple CIPDSpec's with different versions of
    the same package, they'll put their packages in the same directory.

    The packages are stored with their instance_id as the file name.

    This is `[CACHE]/3pp_cipd/name/of/package`.
    """
    ret = self._api.path['cache'].join('3pp_cipd')
    return ret.join(*self._pkg.split('/'))

  def local_pkg_path(self):
    """The path to the local package if it's available on disk. If it's not on
    disk, this returns None.
    """
    iid = self._resolved_instance_id
    ret = self._local_pkg_path_dir.join(iid)
    if iid and self._api.path.exists(ret):
      return ret

  def _ensure_fetched(self):
    """This ensures that this package is fetched locally into the cache.

    Returns a Path to the locally fetched package.

    Raises StepFailure if the package isn't available on the server and isn't
    already fetched locally.
    """
    ret = self.local_pkg_path()
    if ret:
      return ret

    iid = self._resolved_instance_id
    if iid:
      fetch_path = self._local_pkg_path_dir.join(iid)
      ret = fetch_path
    else:
      fetch_path = self._api.path.mkstemp()

    vers = iid if iid else self._symver
    pin = self._api.cipd.pkg_fetch(fetch_path, self._pkg, vers)
    if not iid:
      local_pkg_dir = self._local_pkg_path_dir
      self._api.file.ensure_directory(
        'ensure cipd package cache exists', local_pkg_dir)

      ret = local_pkg_dir.join(pin.instance_id)
      self._api.file.move(
        'move fetched cipd package to cache', fetch_path, ret)

      self._VERSION_CACHE[self._cache_key] = pin.instance_id
    else:
      self._api.path.mock_add_paths(fetch_path)
    return ret

  def build(self, root, install_mode, version_file):
    """Builds the CIPD package from local files.

    This will populate the local CIPD package cache, and will allow uploading
    this package to the server later.

    Args:
      * root (Path) - A path on the local filesystem to the root of the CIPD
        package which is to be created. All files in this directory will be
        included in the CIPD package.
      * install_mode ('symlink'|'copy') - Passed through to the `cipd` recipe
        module. Describes the method that the CIPD client will prefer to deploy
        this package (note that on Windows it will deploy in copy mode
        regardless of this setting).
      * version_file (str) - If provided, will be a forward-slash delimited
        relative path to the root of where the CIPD client should generate
        a version JSON file when it deploys this package. This version file can
        be used by the application to inspect/display its own CIPD instance ID
        and package name.
    """
    pkg_def = self._api.cipd.PackageDefinition(self._pkg, root, install_mode)
    pkg_def.add_dir(root)
    if version_file:
      pkg_def.add_version_file(str(version_file))

    tmpfile = self._api.path.mkstemp()
    pin = self._api.cipd.build_from_pkg(pkg_def, tmpfile)
    self._VERSION_CACHE[self._cache_key] = pin.instance_id

    local_pkg_dir = self._local_pkg_path_dir
    self._api.file.ensure_directory(
      'ensure cipd pkg cache exists', local_pkg_dir)

    local_path = local_pkg_dir.join(pin.instance_id)
    self._api.file.move(
      'mv built cipd pkg to cache', tmpfile, local_path)
