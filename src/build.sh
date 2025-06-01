SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

env GOOS='linux' GOARCH='amd64' go build -a -v -o "$SCRIPT_DIR/../bin/rpm-get"
