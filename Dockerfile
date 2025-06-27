FROM golang:1.23.7 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/bot/bot.go

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates
WORKDIR /app

COPY --from=builder /app/main ./
COPY migrations/* ./migrations/
COPY configs/* ./configs/
COPY .env ./

RUN chmod +x /app/main


RUN ls ./configs
CMD ["./main", "-c", "/app/configs/config.yaml"]