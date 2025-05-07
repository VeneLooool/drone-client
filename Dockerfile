FROM golang:1.24 AS builder

WORKDIR /app
COPY . .

# ВАЖНО: явно указываем GOOS и GOARCH
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN go build -o drone-client ./cmd/drone-client

FROM debian:bullseye-slim
WORKDIR /app

COPY --from=builder /app/drone-client .

COPY swagger-ui ./swagger-ui/
COPY internal ./internal/

EXPOSE ${DRONE_CLIENT_GRPC_PORT} ${DRONE_CLIENT_HTTP_PORT}

CMD ["./drone-client"]
