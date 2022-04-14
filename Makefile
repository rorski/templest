VERSION=$(shell git describe --abbrev=0 --tags)
GOOS:=$(shell go env GOOS)
GOARCH:=$(shell go env GOARCH)

clean:
	rm -f ./templest-*

build:
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o ./templest-${GOOS}-${GOARCH}-${VERSION} ./main.go

test:
	go test

.PHONY: build test
