#!/bin/bash

set -euo pipefail
IFS=$'\n\t'

export DOCKER_BUILDKIT=1

cleanup() {
    docker-compose -f docker-compose.test.yml down --remove-orphans
}

trap cleanup EXIT

docker-compose -f docker-compose.test.yml build worker
docker-compose -f docker-compose.test.yml up -d worker
docker-compose -f docker-compose.test.yml logs -f worker
