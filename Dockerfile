# Build stage
FROM golang:1.24.1 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o main .

# Final stage (pakai base image lebih ringan)
FROM debian:latest

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
