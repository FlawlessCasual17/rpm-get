#!/usr/bin/env bash

SCRIPT_DIR="$(cd -- "$( dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd )"

env GOOS='linux' GOARCH='amd64' go build -a -v -o "$SCRIPT_DIR/build/rpm-get"
