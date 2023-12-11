#!/bin/bash

cd ./internal/app
env GOOS=linux GOARCH=arm go build -o ../../builds/cube2treon_linux_arm32
env GOOS=darwin GOARCH=arm64 go build -o ../../builds/cube2treon_darwin_arm64