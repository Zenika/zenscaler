all:
	go build .
install:
	go install .
test:
	go install ./... && gometalinter -j4 --deadline 300s ./... 
