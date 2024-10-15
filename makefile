# Internal variables you don't want to change
.EXPORT_ALL_VARIABLES:
.DEFAULT_GOAL := help
SHELL := /bin/bash

REPO_ROOT := $(shell git rev-parse --show-toplevel)
VERSION ?= $(shell git describe --tags --always)
SRC_DIR := ./cmd

GOLINT_PATH ?= $(REPO_ROOT)/.tools/golangci-lint
AIR_PATH ?= $(REPO_ROOT)/.tools/air
HTTPYAC_PATH ?= $(REPO_ROOT)/.tools/node_modules/.bin/httpyac

IMAGE_REG ?= ghcr.io
IMAGE_NAME ?= benc-uk/http-toolkit
IMAGE_TAG ?= $(VERSION)
IMAGE_PREFIX := $(IMAGE_REG)/$(IMAGE_NAME)

help: ## 💬 This help message :)
	@figlet $@ || true
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install-tools: ## 🔮 Install dev tools into project .tools directory
	@figlet $@ || true
	@$(GOLINT_PATH) > /dev/null 2>&1 || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./.tools/
	@$(AIR_PATH) -v > /dev/null 2>&1 || (curl -sSL https://github.com/cosmtrek/air/releases/download/v1.51.0/air_1.51.0_linux_amd64 -o .tools/air && chmod +x .tools/air)
	@$(HTTPYAC_PATH) -v > /dev/null 2>&1 || npm install httpyac@latest --prefix .tools
	
lint: ## 🔍 Lint & format check only, sets exit code on error for CI
	@figlet $@ || true
	$(GOLINT_PATH) run

lint-fix: ## 📝 Lint & format, attempts to fix errors & modify code
	@figlet $@ || true
	$(GOLINT_PATH) run --fix

image: check-vars ## 📦 Build container image
	@figlet $@ || true
	docker build --file ./Dockerfile \
	--build-arg VERSION="$(VERSION)" \
	--tag $(IMAGE_PREFIX):$(IMAGE_TAG) .

push: check-vars ## 📤 Push container image to registry
	@figlet $@ || true
	docker push $(IMAGE_PREFIX):$(IMAGE_TAG)

build: ## 🔨 Build binary
	@figlet $@ || true
	go build -ldflags "-X main.version=$(VERSION)" -o bin/http-toolkit $(SRC_DIR)/...

run: ## 🏃 Run with hot reload, used for local development
	@figlet $@ || true
	$(AIR_PATH) -c .air.toml

clean: ## 🧹 Clean up, remove dev data and files
	@figlet $@ || true
	@rm -rf bin report .tools

release: build ## 🚀 Release a new version on GitHub
	@figlet $@ || true
	git push origin $(VERSION)
	@echo "Releasing version $(VERSION) on GitHub, ctrl+c to cancel"
	@sleep 5
	gh release create "$(VERSION)" --title "$(VERSION)" --latest --notes "Release $(VERSION)"
	gh release upload $(VERSION) bin/http-toolkit

test: ## 🧪 Run unit tests
	@figlet $@ || true
	go test -v ./...

test-report: ## 📜 Run unit tests with report
	@figlet $@ || true
	@go install github.com/jstemmer/go-junit-report/v2@latest
	@mkdir -p report
	go test -v ./... | go-junit-report -set-exit-code -out report/unit-tests.xml

test-api: ## 🔬 Run integration tests 
	@figlet $@ || true
	@fuser -k 8000/tcp || true
	REQUEST_DEBUG=false go run $(SRC_DIR) &
	@sleep 3
	$(HTTPYAC_PATH) api/tests.http --all --output short
	@fuser -k 8000/tcp || true

test-api-report: ## 📜 Run integration tests with report
	@fuser -k 8000/tcp || true
	go run $(SRC_DIR) &
	@sleep 3
	@mkdir -p report
	$(HTTPYAC_PATH) api/tests.http --all --junit > report/api-tests.xml
	@fuser -k 8000/tcp || true

version: ## 🥇 Show current version
	@echo $(VERSION)

generate-spec: ## 🧁 Generate OpenAPI spec
	@figlet $@ || true
	api/specs/generate.sh

check-vars:
	@if [[ -z "${IMAGE_REG}" ]]; then echo "💥 Error! Required variable IMAGE_REG is not set!"; exit 1; fi
	@if [[ -z "${IMAGE_NAME}" ]]; then echo "💥 Error! Required variable IMAGE_NAME is not set!"; exit 1; fi
	@if [[ -z "${IMAGE_TAG}" ]]; then echo "💥 Error! Required variable IMAGE_TAG is not set!"; exit 1; fi
	@if [[ -z "${VERSION}" ]]; then echo "💥 Error! Required variable VERSION is not set!"; exit 1; fi
