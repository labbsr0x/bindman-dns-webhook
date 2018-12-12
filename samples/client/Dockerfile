# BUILD
FROM golang:1.11-alpine as builder

RUN apk add --no-cache git mercurial

RUN mkdir -p $GOPATH/src/bindman-dns-webhook/client
WORKDIR $GOPATH/src/bindman-dns-webhook/client

ADD . ./

RUN go get -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /client .

# PACK
FROM scratch

ENV BINDMAN_DNS_MANAGER_ADDRESS ""

EXPOSE 7071

COPY --from=builder /client /
CMD ["./client"]
