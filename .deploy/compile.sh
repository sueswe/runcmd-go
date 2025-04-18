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


APPRELEASEVERSION=$(git rev-list -1 HEAD runcmd.go)
export APPRELEASEVERSION
echo "REV: $APPRELEASEVERSION"

echo ""
echo "compiling: go build runcmd.go -ldflags -X main.REV=$APPRELEASEVERSION"
go build -ldflags "-X main.REV=$APPRELEASEVERSION" -v -o /tmp/runcmd || {
    echo "Status: $?"
    exit 4
}

echo ""
echo "compiling: go build runcmd.go -ldflags -X main.REV=$APPRELEASEVERSION"
GOOS=aix GOARCH=ppc64 go build -ldflags "-X main.REV=$APPRELEASEVERSION" -v -o /tmp/runcmd.aix || {
    echo "Status: $?"
    exit 4
}

echo ""
echo "compiling: go build runcmd.go -ldflags -X main.REV=$APPRELEASEVERSION"
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.REV=$APPRELEASEVERSION" -v -o /tmp/runcmd.win64 || {
    echo "Status: $?"
    exit 4
}



# nun entsprechende scp's durchführen:
echo "running viceversa.sh stp,testta3,14T4 runcmd \$HOME/bin"
cd /tmp/ || exit 1
"$HOME"/bin/vicecersa.sh stp,testta3,14T4 runcmd \$HOME/bin || {
    echo "Status: $?"
    exit 4
}
