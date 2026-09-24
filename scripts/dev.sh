#!/usr/bin/env bash
# Runs a command inside the Go toolchain container (Go is not installed on the host).
set -euo pipefail
cd "$(dirname "$0")/.."
exec docker run --rm -i \
  --user "$(id -u):$(id -g)" \
  -e HOME=/tmp -e GOCACHE=/src/.cache/go-build -e GOMODCACHE=/src/.cache/mod -e GOFLAGS=-buildvcs=false \
  ${SCRAPINGISNOTACRIME_API_KEY:+-e SCRAPINGISNOTACRIME_API_KEY} \
  -v "$PWD":/src -w /src golang:1.23 "$@"
