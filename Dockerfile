FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /order-api cmd/main.go
FROM alpine:latest

WORKDIR /app

COPY --from=builder /order-api /app/order-api
COPY configs/ /app/configs/
COPY .env .

EXPOSE 8080

ENTRYPOINT ["/app/order-api"]
