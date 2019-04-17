#!/bin/bash
set -x

GOOS=darwin GOARCH=amd64 go build -gcflags "all=-trimpath=$GOPATH" -ldflags "-X main.gitHash=$(git describe --always --long --dirty)" -o gfw-"$(git describe --always --long --dirty)" main.go

