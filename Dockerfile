FROM golang:1.23 AS builder

RUN go install github.com/air-verse/air@v1.61.5
RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o main ./cmd/ptm/main.go

# ───────────────────────────────

FROM golang:1.23 AS runtime

RUN go install github.com/air-verse/air@v1.61.5

WORKDIR /app

COPY --from=builder /app/main /app/main

COPY --from=builder /go/bin/air /usr/local/bin/air
COPY --from=builder /go/bin/dlv /usr/local/bin/dlv

COPY .env /app/.env

EXPOSE 8080 40000

COPY configs/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

CMD ["entrypoint.sh"]