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

echo ""
echo "compiling: go build runcmd.go"
go build runcmd.go || {
    echo "Status: $?"
    exit 4
}

echo ""
echo "compiling: GOOS=aix GOARCH=ppc64 go build -o runcmd.aix"
GOOS=aix GOARCH=ppc64 go build -o runcmd.aix || {
    echo "Status: $?"
    exit 4
}

echo ""
echo "compiling: GOOS=windows GOARCH=amd64 go build -o runcmd.win64"
GOOS=windows GOARCH=amd64 go build -o runcmd.win64 || {
    echo "Status: $?"
    exit 4
}



# nun entsprechende scp's durchführen:
echo "running viceversa.sh stp,testta3,14T4 runcmd \$HOME/bin"
viceversa.sh stp,testta3,14T4 runcmd \$HOME/bin || {
    echo "Status: $?"
    exit 4
}
