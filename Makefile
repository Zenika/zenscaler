all:
	go build .
test:
	gometalinter -j4 ./...
