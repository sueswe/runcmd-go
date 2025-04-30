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



echo '
### LINUX #########################################################
'

UMG="
stp,testta3,14T4
lgkk,testta3,19Pt
"

cd /tmp/ || exit 1
"$HOME"/bin/vicecersa.sh ${UMG} runcmd \$HOME/bin/ || {
    echo "Status: $?"
    exit 2
}


echo '
### AIX #########################################################
'
