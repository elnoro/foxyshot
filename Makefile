## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## build/local: build using local go tools (with CGO)
.PHONY: build/local
build/local:
	go build -o bin/local/ .

## build/docker: build in a docker container (no CGO)
.PHONY: build/docker
build/docker:
	docker build -o bin/docker/ -f builder.Dockerfile .

## install: build and install the app manually (requires sudo)
.PHONY: install
install: confirm build/local
	mv bin/local/foxyshot /usr/local/bin

## release/run: create new release on GitHub
.PHONY: release/run
release/run: confirm
	docker run --rm -e GITHUB_TOKEN -v `pwd`:/app -w /app goreleaser/goreleaser release

## release/test: dry run for goreleaser
.PHONY: release/test
release/test:
	docker run --rm -v `pwd`:/app -w /app goreleaser/goreleaser release --snapshot --clean


## lint/golangci: run golangci (Docker required!)
.PHONY: lint/golangci
lint/golangci:
	docker run --rm -v `pwd`:/app -w /app golangci/golangci-lint:v1.52.2 golangci-lint run -v

## lint/deps: tidy & verify
.PHONY: lint/deps
lint/deps:
	go mod tidy
	go mod verify

## test/cov: run tests with race flag
.PHONY: test/run
test/run:
	go test ./... -timeout=30s -race

## test/cov: run tests with coverage and generate coverage.out and coverage.html
.PHONY: test/cov
test/cov:
	go test ./... -coverprofile coverage.out
	go tool cover -html=coverage.out -o coverage.html

## ci/check: run normal checks
ci/check: lint/deps lint/golangci test/cov

## ci/release: run pre-release checks
ci/release: ci/check build/local build/docker release/test