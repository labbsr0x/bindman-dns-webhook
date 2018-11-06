# BUILD
FROM golang:1.11-alpine as builder

RUN apk add --no-cache git mercurial \
    && go get github.com/labbsr0x/sandman-dns-webhook/src/hook

WORKDIR $GOPATH/src/github.com/labbsr0x/sandman-dns-webhook/samples/client

ADD . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /client .

# PACK
FROM scratch

ENV REVERSE_PROXY_ADDRESS ""
ENV MANAGER_ADDRESS ""

COPY --from=builder /client /
CMD ["./client"]
