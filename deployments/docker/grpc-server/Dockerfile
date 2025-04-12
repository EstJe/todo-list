# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY ../../../go.mod go.sum ./
RUN go mod download

COPY ../../../app .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

CMD ["./server"]