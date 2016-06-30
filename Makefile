all:
	go build .
test:
	gometalinter -j4 --deadline 300s ./... 
