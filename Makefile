.PHONY: all

all: deps install
test: unit_test linter

deps:
	go get -u ./...
vendors:
	go get github.com/Masterminds/glide; \
	glide install

install:
	go install --ldflags "-s -w -X github.com/Zenika/zenscaler/core.Version=`git describe --tags --dirty`" .
linter:
	go install . && gometalinter -j4 --deadline 300s --vendor	 ./...
unit_test:
	go test -v --cover ./api
	go test -v --cover ./cmd

build:
	CGO_ENABLED=0 GOGC=off go build --ldflags "-s -w -X github.com/Zenika/zenscaler/core.Version=`git describe --tags --dirty`" .

docker: docker-build docker-image
docker-build:
	mkdir -p ./build
	docker build -t zenscaler-build -f build.Dockerfile . ; \
	docker run -e "CGO_ENABLED=0" -e "GOGC=off" -v $$PWD/build:/build --rm zenscaler-build go build --ldflags "-s -w -X github.com/Zenika/zenscaler/core.Version=`git describe --tags --dirty`" -o /build/zenscaler .

docker-image:
	docker build --rm -t zenika/zenscaler:latest .
