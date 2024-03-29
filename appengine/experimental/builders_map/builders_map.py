from google.appengine.ext import db
import copy
import json
import logging
import urllib
import urllib2
import webapp2

import hardcoded

class Map(db.Model):
  content = db.TextProperty()


LUCI_CHROMIUM_TRY_BUILDERS_URL = (
    'https://cr-buildbucket.appspot.com/_ah/api/swarmbucket/v1/'
    'builders?bucket=luci.chromium.try'
)


class MainPage(webapp2.RequestHandler):
  def get(self):
    self.response.headers['Content-Type'] = 'application/json'
    entity = Map.get_by_key_name('builders')
    if entity:
      self.response.write(entity.content)


class BuildersMap(webapp2.RequestHandler):
  def get(self):
    builders = copy.deepcopy(hardcoded.BUILDERS)

    for builder, builder_def in builders.iteritems():
      master = builder_def.pop('master')
      builder_def['bucket'] = 'master.' + master

      # If this builder is ready for LUCI, use the corresponding LUCI bucket
      # instead.
      luci_migration_url = (
          'https://luci-migration.appspot.com'
          '/masters/%s/builders/%s?format=json' % (
              urllib.quote(master), urllib.quote(builder)))
      try:
        logging.debug('GET %s', luci_migration_url)
        builder_info = json.load(urllib2.urlopen(luci_migration_url))
        if builder_info['luci_is_prod']:
          builder_def['bucket'] = builder_info['bucket']
      except urllib2.HTTPError as ex:
        if ex.code == 404:
          # This builder is not being tracked by LUCI
          pass
        else:
          logging.exception(
              'Failed to fetch LUCI migration state of %s/%s',
              master, builder)
          # Retry the entire task.
          raise

    # Add pure luci.chromium.try builders
    chromium_try_builders = json.load(
        urllib2.urlopen(LUCI_CHROMIUM_TRY_BUILDERS_URL))
    assert len(chromium_try_builders['buckets']) == 1
    for b in chromium_try_builders['buckets'][0]['builders']:
      if b['name'] not in builders:
        builders[b['name']] = {
          'bucket': 'luci.chromium.try',
        }

    b_map = Map(content=json.dumps(builders), key_name='builders')
    b_map.put()


app = webapp2.WSGIApplication([
  ('/', MainPage),
  ('/tasks/builders_map', BuildersMap),
], debug=True)

