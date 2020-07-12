.PHONY: lint build run test snapshot release
lint:
	golangci-lint run
build:
	go build -o ./build/mirrorman ./cmd/mirrorman/main.go

ARGS ?= server
run: build
	unset HTTP_PROXY HTTPS_PROXY && ./build/mirrorman --config ./tests/.mirrorman.yaml ${ARGS}

test:
	./tests/test.sh

snapshot:
	goreleaser release --snapshot --rm-dist

release:
	goreleaser release --rm-dist
