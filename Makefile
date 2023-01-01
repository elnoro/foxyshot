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
	go build

## build/docker: build in a docker container (no CGO)
.PHONY: build/docker
build/docker:
	docker build -o bin -f builder.Dockerfile .

## install: build and install the app manually
.PHONY: install
install: confirm
	go build -o foxyshot
	mv foxyshot /usr/local/bin

## release/run: create new release on GitHub
.PHONY: release/run
release/run: confirm
	docker run --rm -e GITHUB_TOKEN -v `pwd`:/app -w /app goreleaser/goreleaser release

## release/test: dry run for goreleaser
.PHONY: release/test
release/test:
	docker run --rm -v `pwd`:/app -w /app goreleaser/goreleaser release --snapshot --rm-dist


## lint/golangci: run golangci (Docker required!)
.PHONY: lint/golangci
lint/golangci:
	docker run --rm -v `pwd`:/app -w /app golangci/golangci-lint:v1.50.1 golangci-lint run -v

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