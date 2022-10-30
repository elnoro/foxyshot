build:
	go build

buildC:
	docker build -o bin -f builder.Dockerfile .

releaseC:
	docker run --rm -e GITHUB_TOKEN -v `pwd`:/app -w /app goreleaser/goreleaser release

test-releaseC:
	docker run --rm -v `pwd`:/app -w /app goreleaser/goreleaser release --snapshot --rm-dist

test:
	go test ./... -timeout=30s -race

install:
	go build -o foxyshot
	mv foxyshot /usr/local/bin

lintC:
	docker run --rm -v `pwd`:/app -w /app golangci/golangci-lint:v1.50.1 golangci-lint run -v

coverage:
	go test ./... -coverprofile coverage.out
	go tool cover -html=coverage.out -o coverage.html