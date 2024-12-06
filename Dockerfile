# Build stage
FROM golang:1.23 AS builder

# Install Delve debugger
RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy application source code
COPY . .

# Build the application binary with debugging flags
RUN GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o /app/main cmd/main.go

# Runtime stage
FROM debian:bookworm-slim

WORKDIR /app

# Install runtime dependencies
RUN apt-get update && apt-get install -y ca-certificates

# Copy the built binary and Delve debugger
COPY --from=builder /app/main .
COPY --from=builder /go/bin/dlv /usr/local/bin/dlv

# Expose application and debugging ports
EXPOSE 8080 40000

# Default command to run the binary
CMD ["./app/main"]