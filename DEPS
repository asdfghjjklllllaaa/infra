deps = {
  "build":
    "https://chromium.googlesource.com/chromium/tools/build.git",

  "depot_tools":
    "https://chromium.googlesource.com/chromium/tools/depot_tools.git",

  "infra/appengine/swarming":
    ("https://chromium.googlesource.com/infra/swarming.git"
     "@292e1a46f854be008b1dd46f15e07e023fd1eb8a"),

  # Appengine third_party DEPS
  "infra/appengine/third_party/cloudstorage":
    ("https://chromium.googlesource.com/infra/third_party/cloudstorage.git"
     "@ad74316d12e198e0c7352bd666bbc2ec7938bd65"),

  "infra/appengine/third_party/pipeline":
    ("https://chromium.googlesource.com/infra/third_party/"
     "appengine-pipeline.git"
     "@aebc176319885bb7b6372229999b893947248694"),

  ## external deps
  # v1.11.6
  "infra/bootstrap/virtualenv":
    ("https://github.com/pypa/virtualenv.git"
     "@93cfa83481a1cb934eb14c946d33aef94c21bcb0"),
}

hooks = [
  {
    "pattern": ".",
    "action": [
      "python", "-u", "./infra/bootstrap/bootstrap.py",
      "--deps_file", "infra/bootstrap/deps.pyl", "infra/ENV"
    ],
  },
  {
    "pattern": ".",
    "action": [
      "python", "-u", "./infra/bootstrap/get_appengine.py", "--dest=.",
    ],
  },
]

recursion = 1
