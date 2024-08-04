# 🌐 HTTP Toolkit

This is a simple backend service for testing and debugging HTTP requests. It can be deployed when you need to test some networking or HTTP routing in your enviroment, or if you simply need a "Hello World" type thingamajig to deploy. It's basically like [httpbin](https://github.com/postmanlabs/httpbin) but written in Go and it's mine.

![](https://img.shields.io/github/license/benc-uk/http-toolkit)
![](https://img.shields.io/github/last-commit/benc-uk/http-toolkit)
![](https://img.shields.io/github/release/benc-uk/http-toolkit)
![](https://img.shields.io/github/actions/workflow/status/benc-uk/http-toolkit/ci-build.yaml?label=ci-build)

## 📦 Quick Deploy

TBA - Deploy from public container

TBA - Deploy as standalone Go binary

## 🏹 Usage

The routes and features supported are

```text
GET /               - Root URL returns 200/OK
GET /health         - Just returns 200/OK
GET /healthz        - Same

GET /info           - System info as JSON

ANY /inspect        - Returns JSON description of the request
ANY /echo           - Same

ANY /status/{code}  - Return a given status code

GET /word           - Return a random word
GET /word/{count}   - Return several random words
GET /number         - Return a random number between 0-999
GET /number/{max}   - Return a random number between 0 and max
GET /uuid           - Generate a random UUID
GET /uuid/{input}   - Generate a deterministic UUID from input string

ANY /auth/basic/    - Protected by basic auth, see config for username & password
```

## 🛠️ Config

All config is done via environmental variables

| Variable            | Description                                              | Default  |
| ------------------- | -------------------------------------------------------- | -------- |
| PORT                | Port to listen on                                        | "8080"   |
| REQUEST_DEBUG       | Log request details to console                           | true     |
| BODY_DEBUG          | Include body when inspecting requests                    | true     |
| INSPECT_FALLBACK    | Unmatched routes return same as /inspect rather than 404 | true     |
| ROUTE_PREFIX        | Set prefix before all routes                             | "/"      |
| BASIC_AUTH_USER     | Username accepted for basic auth                         | "admin"  |
| BASIC_AUTH_PASSWORD | Password for basic auth user                             | "secret" |

## 🧑‍💻 Local Development

Use the Makefile

```
help                 💬 This help message :)
install-tools        🔮 Install dev tools into project .tools directory
lint                 🔍 Lint & format check only, sets exit code on error for CI
lint-fix             📝 Lint & format, attempts to fix errors & modify code
image                📦 Build container image from Dockerfile
push                 📤 Push container image to registry
build                🔨 Run a local build without a container
run                  🏃 Run application, used for local development
clean                🧹 Clean up, remove dev data and files
release              🚀 Release a new version on GitHub
test                 🧪 Run unit tests
test-report          🧪 Run unit tests to JUnit format report/unit-tests.xml
test-api             🔬 Run integration tests
test-api-report      📜 Run integration tests to JUnit format report/api-tests.xml
```
