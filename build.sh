#!/bin/bash
set +Eeuo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )"

export GIT_COMMIT=$(git rev-list -1 HEAD)
export CURRENT_TIME=$(date -u "+%Y-%m-%d %T UTC")

mkdir -p build
mkdir -p /tmp/go
export GOROOT=/tmp/go
export GO111MODULE=on

go get -d ./...
go build -ldflags "-s -w -X main.versionGitCommitHash=$GIT_COMMIT  -X main.versionCompileTime=$CURRENT_TIME" -o build/snd
