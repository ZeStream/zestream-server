#!/bin/bash

cp .env-local .env
go mod download
go run main.go
