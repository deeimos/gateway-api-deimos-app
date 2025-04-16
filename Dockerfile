FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gateway-api ./cmd/gateway-api/main.go

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates netcat-openbsd

COPY --from=builder /app/gateway-api ./gateway-api

VOLUME /app/config

EXPOSE 1000

HEALTHCHECK \
  --interval=10s \
  --timeout=3s \
  --start-period=5s \
  --retries=3 \
  CMD nc -z localhost 1000 || exit 1

CMD ["./gateway-api"]
