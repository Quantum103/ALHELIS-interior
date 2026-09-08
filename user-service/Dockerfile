FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o user-service ./cmd


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/user-service .

EXPOSE 8082
EXPOSE 50053

CMD ["./user-service"]