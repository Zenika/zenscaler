all:
	go build .
test:
	go install ./... && gometalinter -j4 --deadline 300s ./... 
