build:
	go build

buildC:
	docker build -o bin -f builder.Dockerfile .

test:
	go test ./... -timeout=30s -race

install:
	go build -o foxyshot
	mv foxyshot /usr/local/bin

lint:
	docker run --rm -v `pwd`:/app -w /app golangci/golangci-lint:v1.45 golangci-lint run -v

coverage:
	go test ./... -coverprofile coverage.out
	go tool cover -html=coverage.out -o coverage.html