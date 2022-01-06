FROM golang:1.17-alpine as builder


WORKDIR /app

COPY go.mod go.sum ./
RUN  go mod download

COPY main.go main.go

RUN CGO_ENABLED=0 go build -o build/paranotify main.go

FROM alpine:latest

RUN apk add jq curl
COPY --from=builder /app/build/paranotify /usr/local/bin/paranotify

CMD paranotify
