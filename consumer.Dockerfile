FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

RUN go build -o main

FROM ubuntu:latest AS zestream-consumer

WORKDIR /app

COPY --from=builder /app/main ./

RUN apt update 
RUN apt install -y ffmpeg

CMD ["./main", "--consumer"]
