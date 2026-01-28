.PHONY: build test lint clean docker

build:
	go build -o smtp-gotify ./cmd/smtp-gotify

test:
	go test -v -race -cover ./...

lint:
	golangci-lint run

clean:
	rm -f smtp-gotify

docker:
	docker build -t smtp-gotify .
