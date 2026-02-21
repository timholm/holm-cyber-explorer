.PHONY: build build-local push test lint clean

KO_DOCKER_REPO ?= ghcr.io/timholm/ytarchive

build:
	ko build ./cmd/controller ./cmd/worker ./cmd/collector

build-local:
	ko build --local ./cmd/controller ./cmd/worker ./cmd/collector

push:
	KO_DOCKER_REPO=$(KO_DOCKER_REPO) ko build --push ./cmd/controller ./cmd/worker ./cmd/collector

deploy:
	ko apply -f deploy/kubernetes/

test:
	go test -v ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/

dev-ui:
	cd web && npm run dev

build-ui:
	cd web && npm run build
