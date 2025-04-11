#!/bin/bash


source ~/.profile


cd "$HOME"/compile/runcmd-go || {
    echo "Status: $?"
    exit 4
}



echo "------------------------------------"
env | grep PATH
env | grep LOADED
echo "------------------------------------"


export APPRELEASEVERSION=$(git rev-list -1 HEAD)
echo "Release: $APPRELEASEVERSION"

echo ""
echo "compiling: go build runcmd.go"
go build -ldflags "-X main.AppReleaseVersion=$APPRELEASEVERSION" -v -o /tmp/runcmd || {
    echo "Status: $?"
    exit 4
}

echo ""
echo "compiling: GOOS=aix GOARCH=ppc64 go build -o runcmd.aix"
GOOS=aix GOARCH=ppc64 go build -ldflags "-X main.AppReleaseVersion=$APPRELEASEVERSION" -v -o /tmp/runcmd.aix || {
    echo "Status: $?"
    exit 4
}

echo ""
echo "compiling: GOOS=windows GOARCH=amd64 go build -o runcmd.win64"
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.AppReleaseVersion=$APPRELEASEVERSION" -v -o /tmp/runcmd.win64 || {
    echo "Status: $?"
    exit 4
}



# nun entsprechende scp's durchführen:
echo "running viceversa.sh stp,testta3,14T4 runcmd \$HOME/bin"
cd /tmp/ || exit 1
$HOME/bin/vicecersa.sh stp,testta3,14T4 runcmd \$HOME/bin || {
    echo "Status: $?"
    exit 4
}
