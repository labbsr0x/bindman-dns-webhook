# BUILD
FROM golang:1.11-alpine as builder

RUN apk add --no-cache git mercurial

WORKDIR $GOPATH/src/sandman-dns-webhook/hook

ADD . ./

RUN go get ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /hook .

# PACK
FROM scratch

ENV TAGS "test"

EXPOSE 7070

COPY --from=builder /hook /
CMD ["./hook"]
