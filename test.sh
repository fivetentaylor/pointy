#!/bin/bash

docker compose -f docker-compose.test.yml up -d
go run scripts/loadenv/loadenv.go .env.test.local gotestsum -f dots --jsonfile json.log -- -p 1 ${@:-./...}
docker compose -f docker-compose.test.yml down
