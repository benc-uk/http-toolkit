# Stage 1: Build the Go binary
FROM golang:1.22-alpine AS builder
ARG VERSION=0.0

# Set the working directory inside the container
WORKDIR /build

# Copy the Go module files and download the dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY ./cmd ./cmd

# Build the Go binary
RUN go build -ldflags "-X main.version=${VERSION}" -o http-toolkit ./cmd

# Stage 2: Ship the Go binary, we don't use scratch for TLS reasons
FROM alpine:latest

# Copy the built binary from the builder stage
COPY --from=builder /build/http-toolkit /http-toolkit

# Set the entry point to run the Go binary
ENTRYPOINT ["/http-toolkit"]