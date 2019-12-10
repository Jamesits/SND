#!/bin/bash
set +Eeuo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )"

export OUT_FILE=${OUT_FILE:-snd}

export GIT_COMMIT=$(git rev-list -1 HEAD | cut -c -8)
export CURRENT_TIME=$(date -u "+%Y-%m-%d %T UTC")
export COMPILE_HOST=$(hostname --fqdn)
export GIT_STATUS=""
if output=$(git status --porcelain) && [ -z "$output" ]; then
	export GIT_STATUS="clean"
else 
	export GIT_STATUS="dirty"
fi

mkdir -p build
mkdir -p /tmp/go/gopath

# set it to the actual goroot, else you will have strange errors complaining cannot load bufio
export GO111MODULE=on

# go get -d ./...
go mod download
go mod verify
go build -ldflags "-s -w -X \"main.versionGitCommitHash=$GIT_COMMIT\" -X \"main.versionCompileTime=$CURRENT_TIME\" -X \"main.versionCompileHost=$COMPILE_HOST\" -X \"main.versionGitStatus=$GIT_STATUS\"" -o "build/$OUT_FILE"

# upx
if which upx; then
	! upx "build/$OUT_FILE"
else
	echo "UPX not installed, compression skipped"
fi

# root required
! setcap 'cap_net_bind_service=+ep' "build/$OUT_FILE"

ls -lh "build/$OUT_FILE"
./"build/$OUT_FILE" -version
