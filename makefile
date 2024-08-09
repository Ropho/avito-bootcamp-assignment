export GO111MODULE=on
export GOSUMDB=off
export GOPRIVATE=github.com/Ropho/avito-bootcamp-assignment

##################### PROJECT RELATED VARIABLES #####################
GOPATH?=$(HOME)/go
FIRST_GOPATH:=$(firstword $(subst :, ,$(GOPATH)))
GIT_TAG:=$(shell git describe --exact-match --abbrev=0 --tags 2> /dev/null)
GIT_HASH:=$(shell git log --format="%h" -n 1 2> /dev/null)
GIT_LOG:=$(shell git log --decorate --oneline -n1 2> /dev/null | base64 | tr -d '\n')
GIT_BRANCH:=$(shell git branch 2> /dev/null | grep '*' | cut -f2 -d' ')
GO_VERSION:=$(shell go version)
GO_VERSION_SHORT:=$(shell echo $(GO_VERSION) | sed -E 's/.* go(.*) .*/\1/g')
BUILD_TS:=$(shell date +%FT%T%z)


# App version is sanitized CI branch name, if available.
# Otherwise git branch or commit hash is used.
APP_VERSION:=$(if $(CI_COMMIT_REF_SLUG),$(CI_COMMIT_REF_SLUG),$(if $(GIT_BRANCH),$(GIT_BRANCH),$(GIT_HASH)))

BUILD_ENVPARMS:=CGO_ENABLED=0

LOCAL_BIN:=$(CURDIR)/bin

##################### GOLANG-CI RELATED CHECKS #####################
# Check global GOLANGCI-LINT
GOLANGCI_BIN:=$(LOCAL_BIN)/golangci-lint ## local linter binary path
GOLANGCI_TAG:=1.55.0 ## linter version to use
GOLANGCI_LINTER_IMAGE:="golangci/golangci-lint" ## pipeline linter image to use in ci-lint target

# Check local bin version
ifneq ($(wildcard $(GOLANGCI_BIN)),)
GOLANGCI_BIN_VERSION:=$(shell $(GOLANGCI_BIN) --version)
ifneq ($(GOLANGCI_BIN_VERSION),)
GOLANGCI_BIN_VERSION_SHORT:=$(shell echo "$(GOLANGCI_BIN_VERSION)" | sed -E 's/.* version (.*) built from .* on .*/\1/g')
else
GOLANGCI_BIN_VERSION_SHORT:=0
endif
ifneq "$(GOLANGCI_TAG)" "$(word 1, $(sort $(GOLANGCI_TAG) $(GOLANGCI_BIN_VERSION_SHORT)))"
GOLANGCI_BIN:=
endif
endif

# Check global bin version
ifneq (, $(shell which golangci-lint))
GOLANGCI_VERSION:=$(shell golangci-lint --version 2> /dev/null )
ifneq ($(GOLANGCI_VERSION),)
GOLANGCI_VERSION_SHORT:=$(shell echo "$(GOLANGCI_VERSION)"|sed -E 's/.* version (.*) built from .* on .*/\1/g')
else
GOLANGCI_VERSION_SHORT:=0
endif
ifeq "$(GOLANGCI_TAG)" "$(word 1, $(sort $(GOLANGCI_TAG) $(GOLANGCI_VERSION_SHORT)))"
GOLANGCI_BIN:=$(shell which golangci-lint)
endif
endif

ifeq ($(shell uname -s),Linux)
PROTOC:=$(LOCAL_BIN)/protoc
PROTOC_NAME:=protoc-23.4-linux-x86_64.zip
else
PROTOC:=$(LOCAL_BIN)/protoc
PROTOC_NAME:=protoc-23.4-osx-aarch_64.zip
endif
#####################    GO RELATED CHECKS     #####################
# We always use go 1.16+
ifneq ("1.16","$(shell printf "$(GO_VERSION_SHORT)\n1.16" | sort -V | head -1)")
$(info You could run "scratch install go -v" to help you!)
$(error NEED GO VERSION >= 1.16. Found: $(GO_VERSION_SHORT))
endif
#####################    GO RELATED CHECKS     #####################

.PHONY: all
all: test build ## default scratch target: test and build

.PHONY: help
help: ## show this help
# regex for makefile targets
	@printf "────────────────────────`tput bold``tput setaf 2` Make Targets `tput sgr0`────────────────────────────────\n"
	@sed -ne "/@sed/!s/\(^[^#?=]*:[^=]\).*##\(.*\)/`tput setaf 2``tput bold`\1`tput sgr0`\2/p" $(MAKEFILE_LIST)
# regex for makefile variables
	@printf "────────────────────────`tput bold``tput setaf 4` Make Variables `tput sgr0`───────────────────────────────\n"
	@sed -ne "/@sed/!s/\(.*\)[?:]=\(.*\)##\(.*\)/`tput setaf 4``tput bold`\1: `tput setaf 5`\2`tput sgr0`\3/p" $(MAKEFILE_LIST)

.PHONY: install-lint
install-lint: ## install golangci-lint binary
ifeq ($(wildcard $(GOLANGCI_BIN)),)
	$(info Downloading golangci-lint v$(GOLANGCI_TAG))
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v$(GOLANGCI_TAG)
GOLANGCI_BIN:=$(LOCAL_BIN)/golangci-lint
endif

.PHONY: .lint
.lint: install-lint
	$(info Running lint...)
	$(GOLANGCI_BIN) run --new-from-rev=origin/master --config=.golangci.yaml ./...

.PHONY: lint
lint: .lint ## run golangci-lint only for files that differ from master

.PHONY: .lint-full
.lint-full: install-lint
	$(GOLANGCI_BIN) run --config=.golangci.yaml ./...


.PHONY: lint-full
lint-full: .lint-full ## run golangci-lint for the whole project

.PHONY: .bin-deps
.bin-deps:
	mkdir -p bin
	$(info Installing binary dependencies...)
	GOBIN=$(LOCAL_BIN) go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest && \
	GOBIN=$(LOCAL_BIN) go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.1

.PHONY: bin-deps
bin-deps: .bin-deps ## install binary dependencies required for code generation

.PHONY: .deps
.deps:
	$(info Install dependencies...)
	go mod download

.PHONY: deps
deps: .deps ## install project dependencies

.PHONY: .update-deps
.update-deps:
	$(info Updating dependencies...)
	wget -nv  --header="Accept: application/yaml" -O oas/openapi.yaml https://app.swaggerhub.com/apiproxy/registry/IVANKOVALENKOWORK/backend-bootcamp/1.0.0?resolved=true&flatten=true&pretty=true

.PHONY: update-deps
update-deps: .update-deps ## update dependencies

.PHONY: .pgen-init
.pgen-init:

.PHONY: .test
.test:
	$(info Running tests...)
	go test -v -tags=unit ./... -covermode=count -coverprofile=coverage.out
	go tool cover -func=coverage.out -o=coverage.out

.PHONY: test
test: .test ## run unit tests


ifndef BIN_DIR
BIN_DIR=./bin
endif


.PHONY: .build
.build:
# сначала собирается основной сервис, скачиваются нужные пакеты и все кладется в кеш для дальнейшего использования
	$(info Building...)
	$(BUILD_ENVPARMS) $(GOX_BIN) -output="$(BIN_DIR)/{{.Dir}}" -osarch="$(HOSTOS)/$(HOSTARCH)" -ldflags "$(LDFLAGS)" ./cmd/devices-api
	@if [ -n "$(CMD_LIST)" ] && [ "$(DISABLE_CMD_LIST_BUILD)" != 1 ]; then\
		$(BUILD_ENVPARMS) $(GOX_BIN) -output="$(BIN_DIR)/{{.Dir}}" -osarch="$(HOSTOS)/$(HOSTARCH)" -ldflags "$(LDFLAGS)" $(CMD_LIST);\
	fi

.PHONY: build
build: .build ## build project

.PHONY: .run
.run:
	$(info Running...)
	$(BUILD_ENVPARMS) go run -ldflags "$(LDFLAGS)" ./cmd/devices-api config_name=local_config.yaml

.PHONY: run
run: .run ## run app locally in development mode

.PHONY: .generate
.generate:
	mkdir -p ./api
# ./bin/oapi-codegen --config=oas/config.yaml oas/openapi.yaml 
# $(LOCAL_BIN)/oapi-codegen --config=./oas/config.yaml ./oas/openapi.yaml
	go mod tidy

.PHONY: generate
generate: update-deps bin-deps .generate ## generate code from proto

# Прогоняет линтер из пайплайна. При необходимости указания кастомного образа
# в Makefile необходимо переопределить переменную GOLANGCI_LINTER_IMAGE (см выше).
# Предполагается, что за свежестью образов линтера пользователь следит самостоятельно.
.PHONY: ci-lint
ci-lint: ## run linter via pipeline linter image
	docker run \
		--rm  \
		--volume $(CURDIR):/code \
		--workdir /code $(GOLANGCI_LINTER_IMAGE) \
		/etc/golangci/linters/lint.sh || exit $?


.PHONY: generate
generate: bin-deps .generate ## generate code from proto

# Use make generate before linting & testing
.PHONY: test-integration-ci
test-integration-ci:
	sudo docker compose -f ./test_integration/docker-compose.yaml up --force-recreate --detach
	sleep 10
	./bin/migrate  -path ./migrations -database "postgres://postgres:postgres@localhost:5432/test_db?sslmode=disable" up	
	sudo go test -v ./test_integration/... -tags=integration_repo;
	./bin/migrate  -path ./migrations -database "postgres://postgres:postgres@localhost:5432/test_db?sslmode=disable" down
	sudo docker compose -f ./test_integration/docker-compose.yaml down