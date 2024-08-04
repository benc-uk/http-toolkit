# Set ENV to dev, prod, etc. to load .env.$(ENV) file
ENV ?= 
-include .env
export
-include .env.$(ENV)
export

# Internal variables you don't want to change
REPO_ROOT := $(shell git rev-parse --show-toplevel)
SHELL := /bin/bash

SRC_DIR := ./cmd
GOLINT_PATH := $(REPO_ROOT)/.tools/golangci-lint
AIR_PATH := $(REPO_ROOT)/.tools/air
JUNIT_REPORT_PATH := $(REPO_ROOT)/.tools/go-junit-report

.EXPORT_ALL_VARIABLES:
.PHONY: help image push build run lint lint-fix
.DEFAULT_GOAL := help

help: ## 💬 This help message :)
	@figlet $@ || true
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install-tools: ## 🔮 Install dev tools into project .tools directory
	@figlet $@ || true
	@$(GOLINT_PATH) > /dev/null 2>&1 || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./.tools/
	@$(AIR_PATH) -v > /dev/null 2>&1 || ( wget https://github.com/cosmtrek/air/releases/download/v1.51.0/air_1.51.0_linux_amd64 -q -O .tools/air && chmod +x .tools/air )
	@$(JUNIT_REPORT_PATH) -v > /dev/null 2>&1 || GOBIN=$(REPO_ROOT)/.tools go install github.com/jstemmer/go-junit-report/v2@latest
	
	
lint: ## 🔍 Lint & format check only, sets exit code on error for CI
	@figlet $@ || true
	$(GOLINT_PATH) run

lint-fix: ## 📝 Lint & format, attempts to fix errors & modify code
	@figlet $@ || true
	$(GOLINT_PATH) run --fix

image: check-vars ## 📦 Build container image from Dockerfile
	@figlet $@ || true
	docker build --file ./build/Dockerfile \
	--build-arg BUILD_INFO="$(BUILD_INFO)" \
	--build-arg VERSION="$(VERSION)" \
	--tag $(IMAGE_PREFIX):$(IMAGE_TAG) . 

push: check-vars ## 📤 Push container image to registry
	@figlet $@ || true
	docker push $(IMAGE_PREFIX):$(IMAGE_TAG)

build: ## 🔨 Run a local build without a container
	@figlet $@ || true
	go build -o bin/http-tool $(SRC_DIR)/...

run: ## 🏃 Run application, used for local development
	@figlet $@ || true
	$(AIR_PATH) -c .air.toml

clean: ## 🧹 Clean up, remove dev data and files
	@figlet $@ || true
	@rm -rf bin report .tools

release: ## 🚀 Release a new version on GitHub
	@figlet $@ || true
	@echo "Releasing version $(VERSION) on GitHub"
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]
	gh release create "$(VERSION)" --title "v$(VERSION)" \
	--notes-file docs/release-notes.md \
	--latest 

test: ## 🧪 Run unit tests
	@figlet $@ || true
	go test -v ./...

test-report: ## 🧪 Run unit tests to JUnit format report/unit-tests.xml
	@figlet $@ || true
	go install github.com/jstemmer/go-junit-report/v2@latest
	mkdir -p report
	go test -v ./... | $(JUNIT_REPORT_PATH) -set-exit-code -out report/unit-tests.xml

test-api: ## 🔬 Run integration tests 
	@figlet $@ || true
	fuser -k 8080/tcp || true
	REQUEST_DEBUG=false go run $(SRC_DIR) &
	sleep 2
	npx httpyac api/tests.http --all --output short
	fuser -k 8080/tcp || true

test-api-report: ## 📜 Run integration tests to JUnit format report/api-tests.xml
	fuser -k 8080/tcp || true
	go run $(SRC_DIR) &
	sleep 2
	mkdir -p report
	npx httpyac api/tests.http --all --junit > report/api-tests.xml
	fuser -k 8080/tcp || true

check-vars:
	@if [[ -z "${IMAGE_REG}" ]]; then echo "💥 Error! Required variable IMAGE_REG is not set!"; exit 1; fi
	@if [[ -z "${IMAGE_NAME}" ]]; then echo "💥 Error! Required variable IMAGE_NAME is not set!"; exit 1; fi
	@if [[ -z "${IMAGE_TAG}" ]]; then echo "💥 Error! Required variable IMAGE_TAG is not set!"; exit 1; fi
	@if [[ -z "${VERSION}" ]]; then echo "💥 Error! Required variable VERSION is not set!"; exit 1; fi
