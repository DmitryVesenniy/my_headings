# BUILDER
FROM golang:alpine AS builder
LABEL stage=gobuilder

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0  

RUN go build -ldflags="-s -w" -o server ./cmd/api

# PUBLICATION
FROM alpine
RUN apk update --no-cache && apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/server /app/server

ENTRYPOINT ["/app/server"]