FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o pettime-api ./cmd/api

# Production image - use same alpine version as golang:1.23-alpine
FROM golang:1.23-alpine

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/pettime-api .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./pettime-api"]