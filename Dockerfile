FROM golang:1.23 AS builder

RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o /app/main cmd/main.go

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates

COPY --from=builder /app/main ./main
COPY --from=builder /go/bin/dlv /usr/local/bin/dlv

EXPOSE 8080 40000

CMD ["./main"]