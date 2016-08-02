.PHONY: all

all: deps install

deps:
	go get -u ./...
vendors:	
	go get github.com/Masterminds/glide; \
	glide install

install:
	go install .
test:
	go install . && gometalinter -j4 --deadline 300s ./...

build:
	CGO_ENABLED=0 GOGC=off go build --ldflags "-s -w -X github.com/Zenika/zscaler/core.Version=`git describe --tags`" .

docker: docker-build docker-image
docker-build:
	docker build -t zscaler-build -f build.Dockerfile . ; \
	docker run -v $(pwd)/build:/build -it zscaler-build go build -o /build/zscaler .
docker-image:
	docker build -t zscaler .
