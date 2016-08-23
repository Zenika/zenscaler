.PHONY: all

all: deps install
test: unit_test linter

deps:
	go get -u ./...
vendors:
	go get github.com/Masterminds/glide; \
	glide install

install:
	go install --ldflags "-s -w -X github.com/Zenika/zscaler/core.Version=`git describe --tags --dirty`" .
linter:
	go install . && gometalinter -j4 --deadline 300s ./...
unit_test:
	go test -v --cover ./api
	go test -v --cover ./cmd

build:
	CGO_ENABLED=0 GOGC=off go build --ldflags "-s -w -X github.com/Zenika/zscaler/core.Version=`git describe --tags --dirty`" .

docker: docker-build docker-image
docker-build:
	mkdir -p ./build
	docker build -t zscaler-build -f build.Dockerfile . ; \
	docker run -e "CGO_ENABLED=0" -e "GOGC=off" -v $$PWD/build:/build --rm zscaler-build go build --ldflags "-s -w -X github.com/Zenika/zscaler/core.Version=`git describe --tags --dirty`" -o /build/zscaler .

docker-image:
	docker build --rm -t zenika/zscaler:latest .
