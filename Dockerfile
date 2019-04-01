# BUILD
FROM golang:1.11-alpine as builder

RUN apk add --no-cache git mercurial

ENV p $GOPATH/src/github.com/labbsr0x/bindman-dns-webhook

ADD ./ ${p}
WORKDIR ${p}
RUN go get -v ./...