.PHONY: all

all: deps install
test: unit_test linter

deps:
	go get -u ./...
vendors:	
	go get github.com/Masterminds/glide; \
	glide install

install:
	go install --ldflags "-s -w -X github.com/Zenika/zscaler/core.Version=`git describe --tags`" .
linter:
	go install . && gometalinter -j4 --deadline 300s ./...
unit_test:
	go test -v --cover ./api

build:
	CGO_ENABLED=0 GOGC=off go build --ldflags "-s -w -X github.com/Zenika/zscaler/core.Version=`git describe --tags`" .

docker: docker-build docker-image
docker-build:
	docker build -t zscaler-build -f build.Dockerfile . ; \
	docker run -e "CGO_ENABLED=0" -e "GOGC=off" -v $$PWD/build:/build --rm zscaler-build go build --ldflags "-s -w -X github.com/Zenika/zscaler/core.Version=`git describe --tags`" -o /build/zscaler .

docker-image:
	docker build --rm -t zscaler .
