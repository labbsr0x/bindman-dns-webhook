# BUILD
FROM golang:1.11-alpine as builder

RUN mkdir -p $GOPATH/src/github.com/labbsr0x/sandman-dns-webhook/samples/client
WORKDIR $GOPATH/src/github.com/labbsr0x/sandman-dns-webhook/samples/client

ADD . ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /client .

# PACK
FROM scratch

ENV REVERSE_PROXY_ADDRESS ""
ENV MANAGER_ADDRESS ""

COPY --from=builder /hook /
CMD ["./hook"]
