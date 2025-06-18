# Build stage
FROM golang:1.24-alpine AS builder

# Install build tools
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first (to leverage caching)
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the entire source code
COPY .. .

# Build the Go binary
RUN go build -o job-runner ./cmd/main

# Final image
FROM alpine:latest

# Install CA certificates (for HTTPS, etc.)
RUN apk add --no-cache ca-certificates

# Set working directory
WORKDIR /app

# Copy built binary from builder
COPY --from=builder /app/job-runner .

# Expose metrics port (Prometheus)
EXPOSE 2112

# Run the binary
ENTRYPOINT ["./job-runner"]
