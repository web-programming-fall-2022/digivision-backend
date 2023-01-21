help:
	@echo "Please use \`make <ROOT>' where <ROOT> is one of"
	@echo "  update-dependencies    to update dependencies"
	@echo "  download-dependencies  to download the dependencies"
	@echo "  build                  to build the main binary for current platform"
	@echo "  clean                  to remove generated files"
	@echo "  format                 to format code using goimports"
	@echo "  test                   to run unit tests"
	@echo "  docker-build           to build docker image, you should specify VERSION"
	@echo "  docker-push            to push docker image to registry, you should specify VERSION"
	@echo "  deploy-dev             to deploy dvs along with dev services (rabbitmq, postgresql, ...) to docker compose for dev purposes"

SRCS = $(patsubst ./%,%,$(shell find . -name "*.go"))
GOLANGCI_LINT := $(GOPATH)/bin/golangci-lint

build: $(SRCS)
	$(GO_VARS) $(GO) build -a -installsuffix cgo -o dvs -ldflags="$(LD_FLAGS)" -tags static main.go

format:
	 which goimports || GO111MODULE=off GOPROXY=$(GOPROXY) go get -u golang.org/x/tools/cmd/goimports
	 find . -type f -name "*.go" | xargs -n 1 -I R goimports -w R
	 find . -type f -name "*.go" | xargs -n 1 -I R gofmt -s -w R

clean:
	rm -rf dvs

download-dependencies:
	GOSUMDB=off GOPROXY=$(GOPROXY) GOPRIVATE=$(GOPRIVATE) go mod download

update-dependencies:
	GOSUMDB=off GOPROXY=$(GOPROXY) GOPRIVATE=$(GOPRIVATE) go mod tidy -compat=1.17

install-linter:
	which $(GOLANGCI_LINT) || curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/$(GOLANGCI_LINT_VERSION)/install.sh \
		| sed -e '/install -d/d' \
		| sh -s -- -b $(GOPATH)/bin $(GOLANGCI_LINT_VERSION)

lint: install-linter
	GO111MODULE=on CGO_ENABLED=0 $(GOLANGCI_LINT) run

test:
	$(GO_VARS) $(GO) test -p=1 -coverpkg=./... -coverprofile=coverage.out ./...
	$(GO_VARS) $(GO) tool cover -func=coverage.out

docker-build:
	docker build -t $(DOCKER_IMAGE):$(VERSION) -- .
	docker tag $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):latest

docker-push:
	docker push $(DOCKER_IMAGE):$(VERSION)

deploy-dev: .env
	docker-compose up --build

.env:
	@echo "Please provide .env file first, you may use .env.sample file"
	@exit 1

## Project Vars ##########################################################
ROOT := github.com/arimanius/digivision-backend
DOCKER_IMAGE := ghcr.io/arimanius/digivision-backend
.PHONY: help clean update-dependencies test docker

## Commons Vars ##########################################################
OS = $(shell echo $(shell uname -s) | tr '[:upper:]' '[:lower:]')
GO_VARS = GOARCH=amd64 CGO_ENABLED=0 GOOS=$(OS)
GO_PACKAGES := $(shell go list ./...)
GO ?= go
GIT ?= git
COMMIT := $(shell $(GIT) rev-parse HEAD)
BUILD_TIME := $(shell LANG=en_US date +"%F_%T_%z")
LD_FLAGS := -X $(ROOT).Version=$(VERSION) -X $(ROOT).Commit=$(COMMIT) -X $(ROOT).BuildTime=$(BUILD_TIME) -X $(ROOT).Title=dvs
GOLANGCI_LINT_VERSION ?= v1.45.2
