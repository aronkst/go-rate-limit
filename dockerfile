FROM golang:1.21.5 AS builder

WORKDIR /app

RUN mkdir -p cmd/ratelimit
RUN mkdir -p internal/config
RUN mkdir -p internal/middleware
RUN mkdir -p internal/ratelimit
RUN mkdir -p internal/storage

COPY go.mod ./
COPY go.sum ./
COPY cmd/ratelimit/main.go ./cmd/ratelimit
COPY internal/config/config.go ./internal/config
COPY internal/config/config_test.go ./internal/config
COPY internal/middleware/rate_limiter_middleware.go ./internal/middleware
COPY internal/middleware/rate_limiter_middleware_test.go ./internal/middleware
COPY internal/ratelimit/rate_limit_interface.go ./internal/ratelimit
COPY internal/ratelimit/rate_limit.go ./internal/ratelimit
COPY internal/ratelimit/rate_limit_test.go ./internal/ratelimit
COPY internal/storage/redis_client.go ./internal/storage
COPY internal/storage/storage_interface.go ./internal/storage

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o go-rate-limit cmd/ratelimit/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/go-rate-limit ./

CMD ["./go-rate-limit"]
