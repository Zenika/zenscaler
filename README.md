zScaler
=======

zScaler aims to be an environement-agnostic, simple yet intelligent scaler.
Target environements are Kubernetes, Rancher, Mesos and Swarm.
Currently, the only supported target is the docker engine.

Usage
-----

### Configuration file

Please refer to the [wiki page](https://github.com/Zenika/zscaler/wiki/Configuration#configuration-file) for details.

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

You'll need Go (1.5+) and an orchestrator:
* doker (api 1.22+) or docker-swarm
    * docker-compose (1.7.1) if you use it
* kubernetes (_TBD_)
* Mesos (_TBD_)

- Install Goalang and set you `$GOPATH`
- Clone this repo in `$GOPATH/src` and do
```BASH
make all
```
This will download all Go dependencies and install the binary in `$GOPATH/bin`.

Aside : Deploy on EC2
-------------

Some ansible scripts where crafted ahead of development to bootstrap cluster deployment when needed. You can find them under the `deploy/` directory.

_Project supported by Maximilien Richer, supervised by Sylvain Revereault (Zenika Rennes)_
