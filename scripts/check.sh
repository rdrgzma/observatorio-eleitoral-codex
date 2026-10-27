#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
docker compose config --quiet
python3 -m unittest discover -s tests -v
docker run --rm -v "$PWD:/src:ro" ghcr.io/astral-sh/ruff:0.14.2 check /src/services/engine /src/scripts /src/tests
docker run --rm -v "$PWD/services/api:/src" -w /src golang:1.26-bookworm sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test -race ./...'
docker compose --profile tools run --rm engine python -m unittest discover -s tests -v
docker compose run --rm --no-deps painel sh -c 'composer validate --no-check-publish && composer check-platform-reqs && vendor/bin/pint --test && vendor/bin/phpunit --fail-on-warning'
