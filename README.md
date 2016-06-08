zscaler
=======

Zscaler aims to be an environement-agnostic, simple yet intelligent scaler.

Target environements are Kubernetes, Rancher, Mesos and Swarm.

Deploy on EC2
-------------

You'll need:
- `ansible 2.1+`
- `docker  1.10.3 (API 1.22)`

First export some parameters:
```
export AWS_ACCESS_KEY_ID='ACME******'
export AWS_SECRET_ACCESS_KEY='acme*************'
export ANSIBLE_HOST_KEY_CHECKING=False
```

Provision swarm cluster:
```
ansible-playbook aws-provision.yaml
```

Swarm socket is at `<master>:4000`, you can check it with `docker -H <master>:4000 info`.

Disallocate cluster (using dynamic inventory)
```
ansible-playbook -i ec2.py aws-terminate.yaml
```

_Project supported by Maximilien Richer, supervised by Sylvain Revereault (Zenika Rennes)_
