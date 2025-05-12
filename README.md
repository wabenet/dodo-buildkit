![release](https://img.shields.io/github/v/release/wabenet/dodo-buildkit?sort=semver)
![build](https://img.shields.io/github/actions/workflow/status/wabenet/dodo-buildkit/ci.yaml?branch=main&logo=github)
![codecov](https://img.shields.io/codecov/c/github/wabenet/dodo-buildkit?logo=codecov)
![codeclimate](https://img.shields.io/codeclimate/maintainability/wabenet/dodo-buildkit?logo=codeclimate)
![license](https://img.shields.io/github/license/wabenet/dodo-buildkit)

# dodo-buildkit

Adds support for buildkit as a dodo builder plugin.

In the current version, it will ony detect the default builder that comes with the
docker daemon, similar to the docker cli.

## installation

This plugin is already included in the dodo default distribution.

If you compile your own dodo distribution, you can add this plugin with the
following generate config:

```yaml
plugins:
  - import: github.com/wabenet/dodo-buildkit
```

Alternatively, you can install it as a standalone plugin by downloading the
correct file for your system from the [releases page](https://github.com/wabenet/dodo-buildkit/releases),
then copy it into the dodo plugin directory (`${HOME}/.dodo/plugins`).

## configuration

The plugin will recognize `DOCKER_HOST` and `DOCKER_CERT_PATH` environment
variables, as well as registry authentication configuration in `~/.docker`,
similar to the normal [docker cli](https://docs.docker.com/engine/reference/commandline/cli/).

## license & authors

```text
Copyright 2025 Ole Claussen

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
