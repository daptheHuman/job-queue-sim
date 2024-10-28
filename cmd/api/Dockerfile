FROM golang:1.22.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api main.go

FROM alpine:3.12
WORKDIR /app
COPY --from=builder /app/bin/api .

# migrations folder
COPY --from=builder /app/migrations ./migrations

CMD ["./api"]
