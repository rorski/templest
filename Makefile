VERSION=$(shell git describe --abbrev=0 --tags)
GOOS:=$(shell go env GOOS)
GOARCH:=$(shell go env GOARCH)

clean:
	rm -f ./templizer-*

build:
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o ./templizer-${GOOS}-${GOARCH}-${VERSION} ./main.go ./render.go

test:
	go test

.PHONY: build test
