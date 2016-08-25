FROM golang:1.6.3

RUN go get github.com/Masterminds/glide \
    && go get github.com/golang/lint/golint \
    && go get github.com/kisielk/errcheck


WORKDIR /go/src/github.com/Zenika/zenscaler

COPY glide.yaml glide.yaml
COPY glide.lock glide.lock
RUN glide install

COPY . /go/src/github.com/Zenika/zenscaler
