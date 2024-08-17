# 🌐 HTTP Toolkit

This is a simple backend service for testing and debugging HTTP requests. It can be deployed when you need to test some
networking or HTTP routing in your environment, or if you simply need a "Hello World" type thingamajig to deploy. It's
basically like [httpbin](https://github.com/postmanlabs/httpbin) but written in Go and it's mine.

![](https://img.shields.io/github/license/benc-uk/http-toolkit)
![](https://img.shields.io/github/last-commit/benc-uk/http-toolkit)
![](https://img.shields.io/github/release/benc-uk/http-toolkit)
![](https://img.shields.io/github/actions/workflow/status/benc-uk/http-toolkit/ci-build.yaml?label=ci-build)

## 💾 Install Binaries

- Download released binaries (Linux x86-64) from GitHub:

```bash
curl -sSL https://github.com/benc-uk/http-toolkit/releases/download/v1.0/http-tool -o ./http-tool
chmod +x ./http-tool
```

Alternatively if you have Go installed, the following installs the http-toolkit binary into the current directory:

```bash
GOBIN=$(pwd) go install github.com/benc-uk/http-toolkit/cmd@main && mv ./cmd ./http-toolkit
```

### 📦 Run From Container

Images are published on GitHub [here](https://github.com/benc-uk/http-toolkit/pkgs/container/http-tool)

To run as a container simply run:

```bash
docker run --rm -it -p 8080:8080 ghcr.io/benc-uk/http-tool:v1.0
```

## 🏹 Usage

The routes and features supported are

```text
GET /                - Root URL returns 200/OK
GET /health          - Also returns 200/OK
GET /healthz         - Same

GET /info            - System info as JSON

ANY /inspect         - Returns JSON description of the request
ANY /echo            - Same

ANY /status/{code}   - Return a given status code

GET /word            - Return a random word
GET /word/{count}    - Return several random words
GET /number          - Return a random number between 0-999
GET /number/{max}    - Return a random number between 0 and max
GET /uuid            - Generate a random UUID
GET /uuid/{input}    - Generate a deterministic UUID from input string
GET /uuid/{seconds}  - Delay a response

ANY /auth/basic      - Protected by basic auth, see config for credentials
ANY /auth/jwt        - Protected by JWT (HMAC-SHA256), see config for signing key
```

## 🛠️ Config

Configuration can be done via environmental variables

| Variable            | Description                                                  | Default          |
| ------------------- | ------------------------------------------------------------ | ---------------- |
| PORT                | Port to listen on                                            | "8080"           |
| REQUEST_DEBUG       | Log request details to console                               | true             |
| BODY_DEBUG          | Include body when inspecting requests                        | true             |
| INSPECT_FALLBACK    | Unmatched routes return /inspect rather than 404             | true             |
| ROUTE_PREFIX        | Set prefix before all routes                                 | "/"              |
| BASIC_AUTH_USER     | Username accepted for basic auth                             | "admin"          |
| BASIC_AUTH_PASSWORD | Password for basic auth user                                 | "secret"         |
| JWT_SIGN_KEY        | Signing key used for JWT auth                                | "key_1234567890" |
| CERT_PATH           | Enable TLS, see below                                        | _none_           |
| SPA_PATH            | Enable SPA serving mode, serving the given directory         | _none_           |
| STATIC_PATH         | Enable static file serving mode, serving the given directory | _none_           |

A note on the `INSPECT_FALLBACK` setting, by default this is enabled, this means that going any route not matched by the
app e.g. `/foo/cheese` will result in the same response as going to `/inspect` and that is echoing back details of your
request as JSON. This would include incorrect methods to routes e.g. a POST to `/info`

Any of these settings can also be passed as arguments when starting, run `http-toolkit -help` for details

### Serving static content

The server can act as a simple HTTP file server for SPAs and other static content

- **SPA Mode**  
  Enable with `SPA_PATH` env-var or `-spa-path` argument. The supplied path should point to a SPA containing an index.html and other bundled static files (JS, CSS, images etc). This mode supports client side routing, a common feature of SPAs, which means all requests that aren't for files, will be redirected to index.html
- **Static File Mode**  
  Enable with `STATIC_PATH` env-var or `-static-path` argument. Similar to SPA mode, except requests for missing files or paths will result in a 404, and the contents of directories without an index.html will be listed

If ether of these modes is enabled, all other features are disabled. There is no subpath support so requests must start at the root URL. Sub-paths are supported using route-prefix, but serving SPAs this way is fraught with problems and not recommended.

### Enabling TLS / HTTPS

To enable TLS on the server, set `CERT_PATH` to point to a directory, and this directory should contain both a cert.pem
file and a key.pem file. If found the server starts in TLS mode and will accept HTTPS requests. You can use a self signed
cert of course but you'll get warnings when making requests of course

## 🧑‍💻 Local Development

Use the Makefile, it's super handy and very nice 😎

```
help                 💬 This help message :)
install-tools        🔮 Install dev tools into project .tools directory
lint                 🔍 Lint & format check only, sets exit code on error for CI
lint-fix             📝 Lint & format, attempts to fix errors & modify code
image                📦 Build container image from Dockerfile
push                 📤 Push container image to registry
build                🔨 Run a local build without a container
run                  🏃 Run locally with reload, used for local development
clean                🧹 Clean up, remove dev data and files
release              🚀 Release a new version on GitHub
test                 🧪 Run unit tests
test-report          📜 Run unit tests with report
test-api             🔬 Run integration tests
test-api-report      📜 Run integration tests with report
version              📝 Show current version
```

### 🧪 Testing

This project uses a HTTP file (`api/tests.http`) that can be used a few different ways, you can instal the [VSCode REST CLient](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) or [httpYac](https://marketplace.visualstudio.com/items?itemName=anweber.vscode-httpyac), [httpYac](https://httpyac.github.io/) is preferred as it supports many more features.

You can interactively run & send the requests in the `api/tests.http` file using these extensions, but the main reason to use httpYac, is it has a much richer language & the support of assertions which can turn each of the example requests into integration tests too 👌

For example

```http
GET http://{{ENDPOINT}}/info

?? status == 200
?? body uptime isString
?? body cpuCount isNumber
```

httpYac has a command line tool for running tests and .http files which forms the basis of the `make test-api` and `make test-api-report` makefile targets. It also natively supports .env files, so will load variables from a .env file if one is found.
