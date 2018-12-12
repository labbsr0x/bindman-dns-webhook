# BUILD
FROM golang:1.11-alpine as builder

RUN apk add --no-cache git mercurial

WORKDIR $GOPATH/src/bindman-dns-webhook/hook

ADD . ./

RUN go get -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /hook .

# PACK
FROM scratch

ENV BINDMAN_DNS_TTL ""

EXPOSE 7070

COPY --from=builder /hook /
CMD ["./hook"]
