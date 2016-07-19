#!/bin/sh
docker service create --replicas 1 --name helloworld alpine ping docker.com
