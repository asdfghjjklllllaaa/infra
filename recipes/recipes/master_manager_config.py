# Copyright 2015 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

DEPS = [
  'depot_tools/bot_update',
  'depot_tools/gclient',
  'depot_tools/infra_paths',
  'infra_checkout',
  'recipe_engine/buildbucket',
  'recipe_engine/json',
  'recipe_engine/path',
  'recipe_engine/platform',
  'recipe_engine/properties',
  'recipe_engine/python',
  'recipe_engine/step',
]


def RunSteps(api):
  patch_root = 'infra-data-master-manager'
  co = api.infra_checkout.checkout(
      gclient_config_name='infradata_master_manager',
      patch_root=patch_root,
      internal=True)
  co.gclient_runhooks()
  api.python('master manager configuration test',
             co.path.join('infra', 'run.py'),
             ['infra.services.master_manager_launcher',
              '--verify',
              '--ts-mon-endpoint=none',
              '--json-file',
             co.path.join(
                 'infra-data-master-manager',
                 'desired_master_state.json')])


def GenTests(api):
  yield (
      api.test('master_manager_config') +
      api.buildbucket.ci_build(
          project='infra',
          builder='infradata_config',
          git_repo=(
              'https://chrome-internal.googlesource.com/'
              'infradata/master-manager')))
  yield (
      api.test('master_manager_config_patch') +
      api.buildbucket.try_build(
          project='infra',
          builder='infradata_config',
          git_repo=(
              'https://chrome-internal.googlesource.com/'
              'infradata/master-manager')))
