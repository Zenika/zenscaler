zScaler
=======

zScaler aims to be an environment-agnostic, simple and flexible scaler. It plugs itself on any existing infrastructure, probe metrics and scale services according to configured rules by issuing orders to orchestration engines.
Currently, the only supported target is the docker engine.

Requirements
-----
zScaler require Docker with docker-compose (>1.5) in your path or Docker 1.12 installed.

Usage
-----

### Configuration file

Please refer to the dedicated [wiki page](https://github.com/Zenika/zscaler/wiki/Configuration#configuration-file) for details.

Use-case configuration files can be found under the `examples/` folder. Here's a sample:

```YAML
endpoint: "unix:///var/run/docker.sock"
scalers:                               # scaler section
    whoami-compose:                    # custom id
        type: "docker-compose"         # what do we use to scale the service ?
        target: "whoami"               # parameter for docker-compose
        config: "docker-compose.yaml"  # parameter for docker-compose
    whoami2-compose:
        type: "docker-compose"
        target: "whoami2"
        config: "docker-compose.yaml"
rules:                                 # rule section
    whoami-cpu-scale:                  # custom name of the service
        target: "whoami"               # name of service as tagged in orchestrator
        probe: "swarm.cpu_average"     # probe to use
        up: "> 0.75"                   # up rule
        down: "< 0.25"                 # down rule
        scaler: whoami-compose         # refer to any scaler id defined above
        refresh: 3s                    # scaler refresh rate
    whoami2-cpu-scale:
        target: "whoami2"
        probe: "cmd.execute"           # probe can be any binary
        cmd: "./some_script.sh"        # retrieve and write a float to stdout
        up: "> 200"
        down: "< 1.67"
        scaler: whoami2-compose
        refresh: 10s
```

### Command line interface

```BASH
$ zscaler [command]
Available Commands:
  dumpconfig  Dump parsed config file to stdout
  start       Start autoscaler
  version     Display version number
Flags:
  -d, --debug   Activate debug output
```

API
---

A REST API is available at startup, listening on `:3000`.

URL                | HTTP verb | Description
-------------------|-----------|------
/v1/scalers        | GET       | List scalers
/v1/scalers        | POST      | Create scaler
/v1/scalers/:name  | GET       | Describe scalers
/v1/rules          | GET       | List rules
/v1/rules          | POST      | Create rule
/v1/rules/:name    | GET       | Describe rules


Build it
--------

You'll need Go (1.5+).
- Install Goalang and set you `$GOPATH`
- Clone this repo in `$GOPATH/src/github.com/Zenika` and do
```BASH
make all
```
This will download all Go dependencies and install the binary in `$GOPATH/bin`.
Du to the use of `net`, the resulting binary is not static.

### Docker build and docker image
If you have a docker engine, you can build zScaler inside a container and run it as a docker image. To do so run:
```BASH
make docker
```
The build image is tagged `zscaler-build` and the production image `zscaler`.

Aside : Deploy on EC2
-------------

Some ansible scripts where crafted ahead of development to bootstrap cluster deployment when needed. You can find them under the `deploy/` directory.

_Project supported by Maximilien Richer, supervised by Sylvain Revereault (Zenika Rennes)_
