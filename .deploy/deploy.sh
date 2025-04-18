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



cd /tmp/ || exit 1
"$HOME"/bin/vicecersa.sh stp,testta3,14T4 runcmd \$HOME/bin || {
    echo "Status: $?"
    exit 2
}

cd /tmp/ || exit 1
"$HOME"/bin/vicecersa.sh lgkk,testta3,19Pt runcmd \$HOME/bin || {
    echo "Status: $?"
    exit 2
}

