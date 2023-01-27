#!/bin/bash

cp .env-local .evn
go mod download
go run main.go
