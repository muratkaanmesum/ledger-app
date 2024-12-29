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

COPY --from=builder /app/main ./main
COPY --from=builder /app /app
COPY --from=builder /go/bin/air /usr/local/bin/air
COPY --from=builder /go/bin/dlv /usr/local/bin/dlv

EXPOSE 8080 40000

COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

CMD ["entrypoint.sh"]