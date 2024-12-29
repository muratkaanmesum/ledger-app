FROM golang:1.23 AS builder

# Install Air for hot reloading and Delve for debugging
RUN go install github.com/air-verse/air@v1.61.5
RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN GOOS=linux GOARCH=amd64 go build -o main ./cmd/main.go

FROM golang:1.23 AS runtime

# Install Air for hot reloading
RUN go install github.com/air-verse/air@v1.61.5

WORKDIR /app

# Copy built application and other required files
COPY --from=builder /app/main ./main
COPY --from=builder /app /app
COPY --from=builder /go/bin/air /usr/local/bin/air

# Expose ports
EXPOSE 8080 40000

# Use Air for hot reloading
CMD ["air", "-c", ".air.toml"]