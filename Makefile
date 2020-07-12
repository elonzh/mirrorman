.PHONY: build run test snapshot release
build:
	go build -o ./build/mirrorman ./main.go

ARGS ?= server
run: build
	unset HTTP_PROXY HTTPS_PROXY && ./build/mirrorman --config ./tests/.mirrorman.yaml ${ARGS}

test:
	./tests/test.sh

snapshot:
	goreleaser release --snapshot --rm-dist

release:
	goreleaser release --rm-dist
