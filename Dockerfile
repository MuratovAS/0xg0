FROM golang:1-alpine AS builder

RUN mkdir /0xg0
WORKDIR /0xg0

COPY go.mod go.mod
COPY go.sum go.sum
COPY 0xg0.go 0xg0.go

RUN go build -o 0xg0 ./0xg0.go

FROM alpine:latest

COPY --from=builder /0xg0/0xg0 /usr/bin/

EXPOSE 80

VOLUME ["/storage"]
ENTRYPOINT  ["/usr/bin/0xg0", "-s=/storage"]
