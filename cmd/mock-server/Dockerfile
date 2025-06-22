# Multi-stage build for mock server
FROM golang:1.21-alpine AS builder

# Install ca-certificates and git
RUN apk --no-cache add ca-certificates git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the mock server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mock-server ./cmd/mock-server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/mock-server .

# Expose port
EXPOSE 8081

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8081/health || exit 1

# Run the mock server
CMD ["./mock-server"]