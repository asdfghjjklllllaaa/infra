# Copyright 2016 The Chromium Authors. All rights reserved.
# Use of this source code is govered by a BSD-style
# license that can be found in the LICENSE file or at
# https://developers.google.com/open-source/licenses/bsd

runtime: python27
api_version: 1
threadsafe: no

default_expiration: "3600d"

instance_class: F4

define(`_VERSION', `syscmd(`echo $_VERSION')')

ifdef(`PROD', `
automatic_scaling:
  min_idle_instances: 25
  max_pending_latency: 0.2s
')

handlers:
- url: /_ah/api/.*
  script: monorailapp.endpoints

- url: /robots.txt
  static_files: static/robots.txt
  upload: static/robots.txt

- url: /database-maintenance
  static_files: static/database-maintenance.html
  upload: static/database-maintenance.html

- url: /deployed_node_modules
  static_dir: deployed_node_modules
  secure: always
  mime_type: application/javascript

- url: /static/dist
  static_dir: static/dist
  mime_type: application/javascript
  secure: always

- url: /static/js
  static_dir: static/js
  mime_type: application/javascript
  secure: always

- url: /static
  static_dir: static

- url: /_ah/mail/.+
  script: monorailapp.app
  login: admin

- url: /_ah/warmup
  script: monorailapp.app
  login: admin

- url: /.*
  script: monorailapp.app
  secure: always

inbound_services:
- mail
- mail_bounce
ifdef(`PROD', `
- warmup
')

libraries:
- name: endpoints
  version: 1.0
- name: MySQLdb
  version: "latest"
- name: pycrypto
  version: "2.6"
- name: django
  version: 1.9

includes:
- gae_ts_mon

env_variables:
  VERSION_ID: '_VERSION'

skip_files:
- ^(.*/)?#.*#$
- ^(.*/)?.*~$
- ^(.*/)?.*\.py[co]$
- ^(.*/)?.*/RCS/.*$
- ^(.*/)?\..*$
- node_modules/
