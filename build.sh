#!/bin/bash

cd ./internal/app
env GOOS=linux GOARCH=arm go build -o ../../builds/cube2treon_arm32