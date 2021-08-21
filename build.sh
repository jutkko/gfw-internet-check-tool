#!/bin/bash
set -x

GOOS=darwin GOARCH=amd64 go build -gcflags "all=-trimpath=$GOPATH" -ldflags "-X main.gitHash=$(git describe --always --long --dirty)" -o bin/gfw-darwin-"$(git describe --always --long --dirty)" main.go
GOOS=linux GOARCH=amd64 go build -gcflags "all=-trimpath=$GOPATH" -ldflags "-X main.gitHash=$(git describe --always --long --dirty)" -o bin/gfw-linux-"$(git describe --always --long --dirty)" main.go

