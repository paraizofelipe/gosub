.PHONY: lint

clean:
	@find . -iname '*.swp' -delete
	@rm gosub

lint:
	golangci-lint run ./...

test:
	go test ./... -count 1 -v

dep:
	go mod tidy
	go mod download

build:
	@go build
