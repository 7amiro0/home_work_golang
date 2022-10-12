#!/bin/bash

docker-compose down; docker ps -a | awk '{ print $1 }' | xargs docker rm 2> /dev/null

exit 0