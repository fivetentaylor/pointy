#!make
include .env
export $(shell sed 's/=.*//' .env)

CGO_ENABLED?=0
GOCMD=CGO_ENABLED=$(CGO_ENABLED) go
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet

MIGRATION_NAME ?=you_need_to_fill_me_in

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: all
all: help

.PHONY: install

## Initial
.PHONY: install
install: ## install dependencies
	brew tap tinygo-org/tools
	brew install golang-migrate protobuf tinygo mkcert nss poppler poppler-qt5 wv unrtf tidy-html5 stripe/stripe-cli/stripe
	go install github.com/air-verse/air@latest
	go install gorm.io/gen/tools/gentool@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/a-h/templ/cmd/templ@latest
	go install gotest.tools/gotestsum@latest
	go get github.com/getsentry/sentry-go@latest
	npm install -g local-ssl-proxy

.PHONY: setup
setup: docker migrate_up setuplocal ## Setup dependencies

.PHONY: install_helpers
install_helpers: ## Install Reviso specific helper
	go install scripts/loadenv/loadenv.go
	go install scripts/assume-role/assume-role.go

.PHONY: install_dev
install_dev: ## installs the dev dependencies (go + npm) 
	go mod download
	(cd pkg/assets/src && npm install)
	(cd pkg/admin/src && npm install)
	(cd frontend/web && npm install)

.PHONY: docker
docker: ## Start docker in the backgroudn
	docker compose up -d

## Dev
.PHONY: dev
dev: docker ## Run dependencies and server (with reflex)
	CGO_ENABLED=$(CGO_ENABLED) air

.PHONY: dev_stripe_listen
dev_stripe_listen: ## Listen for Stripe events
	stripe listen -l --skip-verify --forward-to https://app.reviso.dev:9090/stripe/webhook

.PHONY: build_assets
build_assets: ## Build assets
	go run ./cmd/reviso/main.go --build

.PHONY: build
build: build_assets ## Build
	CGO_ENABLED=0 go build -ldflags '-s -w' -o reviso ./cmd/reviso/main.go

.PHONY: run_standalone
run_standalone: build ## Run server
	./reviso

.PHONY: debug
debug: docker ## Run server in debug
	CGO_ENABLED=$(CGO_ENABLED) LOG_LEVEL=debug air

.PHONY: server
server: ## Run a local verion of the server
	${GOCMD} run cmd/reviso/main.go --server

.PHONY: worker
worker: ## Run a local verion of the worker
	${GOCMD} run cmd/reviso/main.go --worker

.PHONY: certs
certs: ## Install certs for local tls
	sudo ${GOCMD} run scripts/updatehosts/main.go
	mkdir -p dev/certs
	mkcert -install
	cd dev/certs && mkcert *.reviso.dev && mkcert localhost 127.0.0.1 ::1
	cd dev/certs && mv localhost+2.pem public.crt && mv localhost+2-key.pem private.key

.PHONY: web
web: ## Run nextjs dev
	cd frontend/web && export NODE_EXTRA_CA_CERTS="$$(mkcert -CAROOT)/rootCA.pem" && npm run dev

.PHONY: storybook
storybook:
	(cd pkg/assets/src && npm run storybook)

.PHONY: build_web
build_web: ## Run nextjs build
	cd frontend/web && npm run build

.PHONY: lint_web
lint_web: ## Run linting
	(cd frontend/web && npm run lint)
	(cd pkg/assets/src && npm run prettier)

.PHONY: update_wasm_js
update_wasm_js: ## Pull the wasm_exec.js from go
	GOROOT=$$(go env GOROOT); \
	cp "$$GOROOT/misc/wasm/wasm_exec.js" pkg/assets/static/wasm_exec.js;

.PHONY: down
down: ## Cleanup the docker-compose
	docker compose down

## Testing
.PHONY: test
test: dockertestupd ## Run tests
	./test.sh

bootstrap_test: dockertestupd ## Bootstrap the test env
	(cd frontend/web && ./node_modules/.bin/esbuild test-package.ts --bundle --outfile=out.js --platform=browser)
	cp frontend/web/out.js pkg/assets/static/test-package.js
	go run scripts/loadenv/loadenv.go .env.test.local go run scripts/bootstrap/main.go

run_functional_checks: dockertestupd ## Run the functional tests
	loadenv .env.functional gotestsum -f testname ./checks/...

run_frontend_tests: ## Run the frontend tests
	(cd frontend/web && npm run test); \

run_frontend_tests_with_test_backend: bootstrap_test ## Run the frontend tests with the test backend
	@echo "Killing anything running on port 9191..."	
	lsof -ti:9191 | xargs kill
	@echo "Starting backend test server..."	
	REVISO_MODE=server go run scripts/loadenv/loadenv.go .env.test.local go run cmd/reviso/main.go &
	@backend_pid=$$!; \
	echo "Waiting for backend to become ready..."; \
	until curl --output /dev/null --silent --fail http://localhost:9191; do \
	    printf '.'; \
	    sleep 1; \
	done; \
	echo "Backend is ready."; \
	echo "Running frontend tests..."; \
	(cd frontend/web && npm run test); \
	echo "Tests complete. Shutting down backend test server..."; \
	lsof -ti:9191 | xargs kill

kill_zombie_test_backend: ## Kill the test backend
	lsof -ti:9191 | xargs kill

.PHONY: dockertestup
dockertestup:
	CUR_DIR=$(shell pwd) && \
		go run scripts/loadenv/loadenv.go .env.test.local docker compose -f docker-compose.test.yml up

.PHONY: dockertestupd
dockertestupd:
	CUR_DIR=$(shell pwd) && \
		go run scripts/loadenv/loadenv.go .env.test.local docker compose -f docker-compose.test.yml up -d

.PHONY: dockertestdown
dockertestdown:
	docker compose -f docker-compose.test.yml down

## Demo
demo_temp: ## Run a demo of the ai
	go run demo/ai/main.go


## Gen
gen: gen_db gen_graphql gen_protos ## Generate all the things

gen_graphql:
	${GOCMD} run github.com/99designs/gqlgen generate --config gqlgen.yml
	
gen_db: gen_grom dump_schema ## Generate models & queries + dump schema

gen_grom: ## Generate models & queries
	${GOCMD} run dev/gen_query/main.go

.PHONY: gen_protos
gen_protos:
	@echo $(shell find pkg -type f -name "*.proto" -print | grep -v '/frontend/')
	protoc --go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative \
			$(shell find pkg -type f -name "*.proto" -print | grep -v '/frontend/')

.PHONY: gen_web
gen_web: ## Run graphql codegen
	(cd pkg/assets/src && export NODE_EXTRA_CA_CERTS="$$(mkcert -CAROOT)/rootCA.pem" && npm run gen)
	(cd pkg/assets/src && npm run prettier)

.PHONY: gen_web_debug
gen_web_debug: ## Run graphql codegen (debug to gen_output.log)
	(cd pkg/assets/src && export NODE_EXTRA_CA_CERTS="$$(mkcert -CAROOT)/rootCA.pem" && npm run gen_debug)
	(cd pkg/assets/src && npm run prettier)

gen_templ: ## Run templ codegen
	templ generate -path pkg/

## Database
.PHONY: dump_schema
dump_schema: ## Dump the database schema
	pg_dump --schema-only ${DATABASE_URL} > pkg/db/schema.sql

.PHONY: migrate_up
migrate_up: ## Migrate the database up
	migrate -source file://pkg/db/migrations -database ${DATABASE_URL} up

.PHONY: migrate_down
migrate_down: ## Migrate the database down one step
	migrate -source file://pkg/db/migrations -database ${DATABASE_URL} down 1

.PHONY: migrate_force_version
migrate_force_version: ## Migrate the database up
	migrate -source file://pkg/db/migrations -database ${DATABASE_URL} force ${VERSION}

.PHONY: create_migration
create_migration:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir pkg/db/migrations -seq $$name

.PHONY: stage_migrate_up
stage_migrate_up: ## Migrate the staging database up
	migrate -source file://pkg/db/migrations -database ${STAGE_DATABASE_URL} up

.PHONY: prod_migrate_up
prod_migrate_up: ## Migrate the production database up
	migrate -source file://pkg/db/migrations -database ${PROD_DATABASE_URL} up

.PHONY: setuplocal
setuplocal: ## Create development buckets / tables
	mkdir -p localdata/redis
	go run scripts/setuplocal/main.go

.PHONY: migrate_dynamo
migrate_dynamo: ## Migrate the development table
	 go run scripts/dynamo/main.go
	 go run dev/dynamo/addGSI1/main.go

## Help:
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)
