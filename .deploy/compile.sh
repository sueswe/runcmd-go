#!/bin/bash

source ~/.profile

cd "$HOME"/compile/runcmd-go || {
    echo "Status: $?"
    exit 4
}

echo "------------------------------------"
env | grep PATH
env | grep LOADED
pwd
echo "------------------------------------"

APPRELEASEVERSION=$(git rev-list -1 HEAD)
export APPRELEASEVERSION
echo "REV: $APPRELEASEVERSION"

echo ""
echo "compiling: go build runcmd.go -ldflags -X main.REV=$APPRELEASEVERSION"
go build -ldflags "-X main.REV=$APPRELEASEVERSION" -v -o /tmp/runcmd || {
    echo "Status: $?"
    exit 4
}

echo ""
echo "compiling: GOOS=aix GOARCH=ppc64 go build runcmd.go -ldflags -X main.REV=$APPRELEASEVERSION"
GOOS=aix GOARCH=ppc64 go build -ldflags "-X main.REV=$APPRELEASEVERSION" -v -o /tmp/runcmd.aix || {
    echo "Status: $?"
    exit 4
}

echo ""
echo "compiling: GOOS=windows GOARCH=amd64 go build runcmd.go -ldflags -X main.REV=$APPRELEASEVERSION"
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.REV=$APPRELEASEVERSION" -v -o /tmp/runcmd.win64 || {
    echo "Status: $?"
    exit 4
}
