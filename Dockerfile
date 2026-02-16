# Dockerfile (gunakan Go >= versi yang diminta oleh go.mod)
FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git build-base

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/bin/we-saving-api ./cmd/we-saving-api

FROM alpine:3.18
RUN apk add --no-cache ca-certificates netcat-openbsd

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app
COPY --from=builder /app/bin/we-saving-api /app/we-saving-api
COPY docker-entrypoint.sh /app/docker-entrypoint.sh
RUN chmod +x /app/docker-entrypoint.sh

USER appuser
EXPOSE 8080
ENTRYPOINT ["/app/docker-entrypoint.sh"]
CMD ["/app/we-saving-api"]
