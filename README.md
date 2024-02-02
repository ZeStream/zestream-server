# ZeStream Server

#### To support development of ZeStream, please [![](https://img.shields.io/static/v1?label=Sponsor&message=%E2%9D%A4&logo=GitHub&color=%23fe8e86)](https://github.com/sponsors/abskrj)


A media streaming service, which can stream on demand video, image (with tansformations) and audio and works as CDN.

## Getting Started

#### Deploying on Docker Container

```bash
# create .env file from .env-template

docker compose up
```

## Dev Setup

1. Install Golang v1.19 from [Go.dev](https://go.dev/doc/install)
2. Fork this repo, and clone the forked repo
3. `cd zestream-server`
4. `go get .`
5. `go run main.go`

## Implementation

1. User calls API to prcoess the video.
2. The API controller queues the event in a message queue (Kafka) and calls the given webhook.
3. A ZeStream worker polls the event from queue which contains file url
4. Worker fetches the file to local disk
5. Another worker starts a FFmpeg thread to process the video
6. Output is stored in local disk
7. After FFmpeg finishes processing, another worker pushes the files to cloud storage
8. CDN is connected to storage for fetching the files.

## How to contribute?

Check out [contribution guidelines](https://github.com/ZeStream/zestream-server/blob/main/CONTRIBUTING.md)
