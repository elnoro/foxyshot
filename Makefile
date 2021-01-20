build:
	go build 

test:
	go test ./... -timeout=30s -race

install:
	go build -o foxyshot
	mv foxyshot /usr/local/bin

coverage:
	go test ./... -coverprofile coverage.out
	go tool cover -html=coverage.out -o coverage.html