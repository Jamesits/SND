#!/bin/bash
set +Eeuo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )"

export GIT_COMMIT=$(git rev-list -1 HEAD | cut -c -8)
export CURRENT_TIME=$(date -u "+%Y-%m-%d %T UTC")

mkdir -p build
mkdir -p /tmp/go/gopath

# set it to the actual goroot, else you will have strange errors complaining cannot load bufio
export GO111MODULE=on

# go get -d ./...
go mod download
go mod verify
go build -ldflags "-s -w -X \"main.versionGitCommitHash=$GIT_COMMIT\" -X \"main.versionCompileTime=$CURRENT_TIME\"" -o build/snd
./build/snd -version
