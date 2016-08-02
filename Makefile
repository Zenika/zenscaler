.PHONY: all

all: deps install

deps:
	go get -u ./...
install:
	go install .
test:
	go install . && gometalinter -j4 --deadline 300s ./...

build:
	CGO_ENABLED=0 GOGC=off go build --ldflags "-s -w -X github.com/Zenika/zscaler/core.Version=`git describe --tags`" .
