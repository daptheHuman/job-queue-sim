FROM golang:1.22.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/worker
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/worker main.go

FROM alpine:3.12
WORKDIR /app
COPY --from=builder /app/bin/worker .
CMD ["./worker"]
