#!/bin/bash

docker rmi -f newdag/node:latest
docker rmi -f newdag/client:latest
#docker rmi -f newdag/watcher:latest
#clear <none> images
docker rmi $(docker images | grep "^<none>" | awk "{print $3}")
