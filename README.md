zenscaler  [![CircleCI](https://circleci.com/gh/Zenika/zenscaler/tree/master.svg?style=svg&circle-token=78b4c3db440a574eea374cc602addd51a6b5e249)](https://circleci.com/gh/Zenika/zenscaler/tree/master) [![Go Report Card](https://goreportcard.com/badge/github.com/Zenika/zenscaler)](https://goreportcard.com/report/github.com/Zenika/zenscaler) [![GoDoc](https://godoc.org/github.com/Zenika/zenscaler?status.svg)](https://godoc.org/github.com/Zenika/zenscaler) [![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/Zenika/zenscaler/blob/master/LICENSE.md)
=======

zenscaler aims to be an environment-agnostic, simple and flexible scaler. It plugs itself on any existing infrastructure, probe metrics and scale services according to configured rules by issuing orders to orchestration engines.
Currently, the only supported target is the docker engine.

> On the 25-08-2016, **zscaler** became **zenscaler**. Please update you repository path accordingly to avoid broken dependencies.

Requirements
-----

- A running docker engine
- `docker-service` scaler require docker `1.12` (with `--replica`)
- `docker-compose` >1.5 (with `scale`) in your path

Try it out !
------------
Clone the repository in your `$GOPATH` or anywhere else
```BASH
git clone git@github.com:Zenika/zenscaler.git
```
Run an example with træfik on Docker, pulling `zenika/zenscaler` from [docker.io](https://hub.docker.com/r/zenika/zenscaler/)
```BASH
cd ./examples/docker-compose/traefik && docker-compose up
```
Now you can open the [Træfik web UI](http://localhost:8080/) and watch the backend scale up and down in real time!

This demo use a mock probe that report sinus-like values over 1min and causes the backend `whoami` to scale between 1 and 10 containers. The other probe is monitoring the CPU used across all `whoami2` containers. You can stress-test it with the following benchmark:

 ```BASH
ab -c 100 -n 10000000 -H 'Host:whoami2.docker.localhost' http://localhost/
```

Usage
-----

### Configuration file

Please refer to the dedicated [wiki page](https://github.com/Zenika/zenscaler/wiki/Configuration#configuration-file) for details.

Use-case configuration files can be found under the `examples/` folder. Here's a sample:

```YAML
orchestrator:
    engine: "docker"                   # only docker is supported ATM
    endpoint: "tcp://localhost:2376"   # optionnal, default to unix socket
    tls-cacert: "ca.pem"               # if using TLS
    tls-cert: "cert.pem"
    tls-key: "key.pem"
scalers:                               # scaler section
    whoami-compose:                    # custom id
        type: "docker-compose-cmd"     # what do we use to scale the service ?
        project: "traefik"             # parameter for docker-compose
        target: "whoami"               # parameter for docker-compose
        config: "docker-compose.yaml"  # parameter for docker-compose
        upper_count_limit: 0  # 0 mean unlimited, default
        lower_count_limit: 1  # default to ensure availability
    whoami2-compose:
        type: "docker-compose-cmd"
        project: "traefik"
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
$ zenscaler [command]
Available Commands:
  dumpconfig  Dump parsed config file to stdout
  start       Start autoscaler
  version     Display version number
Flags:
  -d, --debug   Activate debug output
```

API
---

A REST API is available at startup, listening on `:3000` (change it with `-l` or `--api-port` flag).
You can find examples on the [wiki page](https://github.com/Zenika/zenscaler/wiki/API).

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

You'll need Go 1.5 or above. Older version _may_ work but are still untested.
- Install Goalang and set you `$GOPATH`
- Run `go get github.com/Zenika/zenscaler` and do
```BASH
make all
```
This will download all dependencies and install the binary in `$GOPATH/bin`.

### Docker build and docker image
If you have a docker engine, you can build zScaler inside a container and run it as a docker image. To do so run:
```BASH
make docker # build in docker and create docker image
```

The golang build image is tagged `zenscaler-build` and the production image `zenscaler`.

Aside : Deploy on EC2
-------------

Some `ansible` scripts where crafted ahead of development to bootstrap cluster deployment when needed. You can find them under the `deploy/` directory.

_Project supported by Maximilien Richer, supervised by Sylvain Revereault (Zenika Rennes)_
