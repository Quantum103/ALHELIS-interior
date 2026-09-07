FROM golang:1.25-alpine AS builder
WORKDIR /build

COPY proto ./proto
COPY design-service ./design-service

WORKDIR /build/design-service
RUN go mod download
RUN go mod tidy
RUN go build -o app ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /build/design-service/app .
CMD ["./app"]