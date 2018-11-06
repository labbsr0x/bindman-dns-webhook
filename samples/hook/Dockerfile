# BUILD
FROM golang:1.11-alpine as builder

RUN apk add --no-cache git mercurial \
    && go get github.com/labbsr0x/sandman-dns-webhook/src/hook

WORKDIR $GOPATH/src/github.com/labbsr0x/sandman-dns-webhook/samples/hook

ADD . ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /hook .

# PACK
FROM scratch

ENV TAGS "test"

COPY --from=builder /hook /
CMD ["./hook"]
