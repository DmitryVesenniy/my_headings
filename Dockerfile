# BUILDER
FROM golang:alpine AS builder
LABEL stage=gobuilder
ENV CGO_ENABLED 0
ENV GOOS linux
RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0  

RUN go build -ldflags="-s -w" -o app ./cmd/api

# PUBLICATION
FROM alpine
RUN apk update --no-cache && apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/app /app/app

ENTRYPOINT ["./app/app"]