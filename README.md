zScaler
=======

Zscaler aims to be an environement-agnostic, simple yet intelligent scaler. Target environements are Kubernetes, Rancher, Mesos and Swarm.

Usage
-----

### Configuration file

```YAML
endpoint: "unix:///var/run/docker.sock"
scalers:                               # scaler section
    whoami-compose:                    # custom id
        type: "docker-compose"         # what do we use to scaler ?
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
        probe: "swarm.cpu_average"
        up: "> 2"
        down: "< 1.5"
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

Dependencies
------------

You'll need Go (1.5+) and an orchestrator :
* doker (api 1.22+) or docker-swarm
    * docker-compose (1.7.1) if you use it
* kubernetes (_TBD_)
* Mesos (_TBD_)

Build it
--------
- Install Goalang and set you `$GOPATH`
- Clone this repo in `$GOPATH/src` and do
```BASH
make all
```
This will download go dependency and install the binary in `$GOPATH/bin`.

Aside : Deploy on EC2
-------------

You'll need:
- `ansible 2.1+`
- `docker  1.10.3 (API 1.22)`

First export some parameters:
```BASH
export AWS_ACCESS_KEY_ID='ACME******'
export AWS_SECRET_ACCESS_KEY='acme*************'
export ANSIBLE_HOST_KEY_CHECKING=False
```

Provision swarm cluster:
```BASH
ansible-playbook aws-provision.yaml
```

Swarm socket is at `<master>:4000`, you can check it with `docker -H <master>:4000 info`.

Disallocate cluster (using dynamic inventory)
```BASH
ansible-playbook -i ec2.py aws-terminate.yaml
```

_Project supported by Maximilien Richer, supervised by Sylvain Revereault (Zenika Rennes)_
