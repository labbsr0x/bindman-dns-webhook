# BUILD
FROM golang:1.11-alpine as builder

RUN apk add --no-cache git mercurial

RUN mkdir -p $GOPATH/src/sandman-dns-webhook/client
WORKDIR $GOPATH/src/sandman-dns-webhook/client

ADD . ./

RUN go get -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /client .

# PACK
FROM scratch

ENV SANDMAN_REVERSE_PROXY_ADDRESS ""
ENV SANDMAN_DNS_MANAGER_ADDRESS ""
ENV SANDMAN_DNS_TAGS ""

EXPOSE 7071

COPY --from=builder /client /
CMD ["./client"]
