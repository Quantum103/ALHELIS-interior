FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY user-service/go.mod user-service/go.sum ./user-service/
COPY proto ./proto

WORKDIR /app/user-service

RUN go mod download

COPY user-service/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o user-service ./cmd


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/user-service/user-service .

EXPOSE 8082
EXPOSE 50053

CMD ["./user-service"]

