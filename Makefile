SHELL:=/bin/bash

GO ?= go

default: build

# project details
APPNAME = reader
PACKAGE = linked-art-reader

# build variables
BRANCH_NAME	?= $(shell git rev-parse --abbrev-ref HEAD)
BUILD_DATE	?= $(shell date '+%Y-%m-%d %H:%M:%S')
GIT_COMMIT	?= $(shell git rev-list -1 HEAD)
VERSION		?= 0.0.1
AUTHOR      ?= $(shell git log -1 --pretty=format:'%an')

BUILD_OVERRIDES = \
	-X "$(PACKAGE)/pkg/app.Name=$(APPNAME)" \
	-X "$(PACKAGE)/pkg/app.Product=$(PRODUCT)" \
	-X "$(PACKAGE)/pkg/app.Branch=$(BRANCH_NAME)" \
	-X "$(PACKAGE)/pkg/app.BuildDate=$(BUILD_DATE)" \
	-X "$(PACKAGE)/pkg/app.Commit=$(GIT_COMMIT)" \
	-X "$(PACKAGE)/pkg/app.Version=$(VERSION)" \
	-X "$(PACKAGE)/pkg/app.Author=$(AUTHOR)" \

# misc
CMDPATH=$(shell ls -1 cmd/reader/main.go | head -1)
TMP=._tmp
WORKSPACE=$(shell pwd)
BUILD_FOLDER=$(shell echo `pwd`/build)

.ONESHELL:
info:
	@echo "CMDPATH      " $(CMDPATH)
	@echo "BUILD_FOLDER " $(BUILD_FOLDER)
	@echo "AUTHOR       " $(AUTHOR)
	@echo "BRANCH       " $(BRANCH_NAME)
	@echo "BUILD_DATE   " $(BUILD_DATE)
	@echo "GIT_COMMIT   " $(GIT_COMMIT)
	@echo "VERSION      " $(VERSION)

.ONESHELL:
install:
	go get -u github.com/vektra/mockery/...
	go get -u github.com/jroimartin/gocui
	go get -u github.com/alecthomas/template

	# get linter - make sure this version matches our CI tool
	command -v golangci-lint || go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.20.0

	go mod download

.ONSHELL:
compile: info
	@echo "BUILD_OVERRIDES" $(BUILD_OVERRIDES)
	@echo "Running go build...."

	@echo "Building reader"
	CGO_ENABLED=0  GOARCH=amd64 \
		go build -a \
		-ldflags='-w -s $(BUILD_OVERRIDES)' \
		-o $(BUILD_FOLDER)/reader cmd/reader/main.go

build: compile

clean:
	@echo "Cleaning 'build' and 'vendor' directories"
	rm -rf build vendor

run: build
	HOST=0.0.0.0 \
	PORT=8080 \
	LOG_LEVEL=debug \
	LOG_FORMAT=text \
	BASE_PATH=`pwd` \
	ENVIRONMENT_NAME=local \
		$(BUILD_FOLDER)/main

run-dev: build
	PORT=8080 \
	LOG_LEVEL=debug \
	LOG_FORMAT=text \
	DB_ENDPOINT=http://localhost:8000 \
	ENVIRONMENT_NAME=dev \
	ENVIRONMENT_NUMBER=2 \
	DEBUG=true \
	BASE_PATH=`pwd` \
		go run $(CMDPATH)

build-mocks:
	mockery -name DynamoDBAPI -dir ./vendor/github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface -recursive

.PHONY: test
test:
	BASE_PATH=`pwd` go test -cover ./...
#	BASE_PATH=`pwd` go test -cover -v  ./...
#	go test -cover ./...

test-integration:
	go test -cover -v  ./test/...

# Sonarqube friendly lint
lint:
	golangci-lint run ./internal/... ./cmd/... ./pkg/...

# Sonarqube friendly test
.PHONY: test-report
test-report:
	set -o pipefail
	go test -cover -v -coverprofile=coverage.out -json ./internal/... ./pkg/... | tee report.json
