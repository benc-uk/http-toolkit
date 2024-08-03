# 🌐 HTTP Toolkit

This is a simple backend service for testing and debugging HTTP requests. It can be deployed when you need to test some networking or HTTP routing in your enviroment, or if you simply need a "Hello World" type thingamajig to deploy. It's basically like [httpbin](https://github.com/postmanlabs/httpbin) but written in Go and it's mine.

Use cases & key features:

- Echo requests back using `/inspect` endpoint
- Get system info with `/info`

Supporting technologies and libraries:

- Go

![](https://img.shields.io/github/license/benc-uk/http-toolkit)
![](https://img.shields.io/github/last-commit/benc-uk/http-toolkit)
![](https://img.shields.io/github/release/benc-uk/http-toolkit)
![](https://img.shields.io/github/checks-status/benc-uk/http-toolkit/main)
![](https://img.shields.io/github/workflow/status/benc-uk/http-toolkit/CI%20Build?label=ci-build)
![](https://img.shields.io/github/workflow/status/benc-uk/http-toolkit/Release%20Assets?label=release)

## 🧑‍💻 Local Development

Use the Makefile

```text
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
```
