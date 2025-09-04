APP_NAME := mcp2grule

.PHONY: test test-coverage test-coverage-html lint clean

test:
	go test -v ./...

test-coverage:
	go test -v ./... -coverprofile=coverage.out

test-coverage-html: test-coverage
	go tool cover -html=coverage.out

lint:
	golangci-lint run

clean:
	rm -rf coverage.out mcp2grule

build:
	go build -o $(APP_NAME) main.go