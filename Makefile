build:
	go build 

test:
	go test ./... -timeout=30s -race

coverage:
	go test ./... -coverprofile coverage.out
	go tool cover -html=coverage.out -o coverage.html