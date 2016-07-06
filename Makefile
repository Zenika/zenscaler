all: deps install
deps:
	go get -u ./...
install:
	go install .
test:
	go install ./... && gometalinter -j4 --deadline 300s ./... 
