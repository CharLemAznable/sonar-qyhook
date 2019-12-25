#!/usr/bin/env bash

go build -ldflags="-s -w" -o sonar-qyhook.bin
upx --brute sonar-qyhook.bin
