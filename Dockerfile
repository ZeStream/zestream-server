FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main

FROM alpine:latest AS run-http
WORKDIR /app
COPY --from=builder /app/main ./
CMD ["./main"]

FROM alpine:latest AS run-consumer
WORKDIR /app
COPY --from=builder /app/main ./
CMD ["./main", "--consumer"]
