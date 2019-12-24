#!/usr/bin/env bash

go build -ldflags="-s -w" -o sonar-qyhook.linux.bin
upx --brute sonar-qyhook.linux.bin
