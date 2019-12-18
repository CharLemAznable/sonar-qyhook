#!/usr/bin/env bash

env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o sonar-qyhook.linux.bin
upx --brute sonar-qyhook.linux.bin
