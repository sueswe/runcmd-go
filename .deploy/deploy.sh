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

echo '### LINUX #########################################################'

stages="stp,testta3 lgkk,testta3 stplgk,test"
for UMG in ${stages}
do
    cd /tmp/ || exit 1
    "$HOME"/bin/vicecersa.sh "${UMG}" runcmd \$HOME/bin/ || {
        echo "Status: $?"
        exit 2
    }
done

echo '### AIX #########################################################'

stages="stp,testta2,11T4 stp,testta2,12Te stp,testta2,13T1 stp,testta2,14T4 stp,testta2,15T5 stp,testta2,16T1 stp,testta2,17T2 stp,testta2,18T2 stp,testta2,19Pt"
for UMG in ${stages}
do
    cd /tmp/ || exit 1
    "$HOME"/bin/vicecersa.sh "${UMG}" runcmd.aix \$HOME/bin/ runcmd || {
        echo "Status: $?"
        exit 2
    }
done
