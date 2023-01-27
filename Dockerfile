FROM golang:1.19

RUN echo "deb http://security.debian.org/debian-security stretch/updates main" >> /etc/apt/sources.list

RUN apt update 
RUN apt install -y ffmpeg

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o main .

WORKDIR /dist

RUN cp /build/main .
RUN cp /build/.env .

EXPOSE 3000

RUN chmod 777 /dist/main
CMD ["/dist/main"]

