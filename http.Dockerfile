FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

RUN go build -o main


FROM ubuntu:latest AS zestream-http

WORKDIR /app

COPY --from=builder /app/main ./

CMD ["./main"]
