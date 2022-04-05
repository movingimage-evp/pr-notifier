ARG GO_VERSION=1.18.0

FROM golang:${GO_VERSION}-buster AS build_base

WORKDIR /build
COPY . /build
RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -x -installsuffix cgo -o pr-notifier .

# Build the Go app
RUN go build -o ./out/pr-notifier .

FROM ubuntu:latest

RUN set -x && apt-get update && \
  DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates && \
  rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=build_base /build/pr-notifier .
ENTRYPOINT ["./pr-notifier"]