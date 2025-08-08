#!/bin/bash

source ~/.profile

name='runcmd'

cd "$HOME"/compile/${name}-go || {
    echo "Status: $?"
    exit 4
}

echo "------------------------------------"
env | grep PATH
env | grep LOADED
echo "------------------------------------"

masterrtc=0

echo '### LINUX #########################################################'
stages="stp,testta3 lgkk,testta3 stplgk,test"
for UMG in ${stages}
do
    cd /tmp/ || exit 1
    remotecommander.rb -d "\$HOME/bin/" -s "${name}" -g ${UMG} || {
        echo "Status: $?"
        masterrtc=2
    }
done
rm -v /tmp/${name}


echo '### AIX #########################################################'
stages="stp,testta2,11T4 stp,testta2,12Te stp,testta2,13T1 stp,testta2,14T4 stp,testta2,15T5 stp,testta2,16T1 stp,testta2,17T2 stp,testta2,18T2 stp,testta2,19Pt"
cd /tmp/ || exit 1
cp ${name}.aix ${name}
for UMG in ${stages}
do
    remotecommander.rb -d "\$HOME/bin/" -s "${name}" -g ${UMG}  || {
        echo "Status: $?"
        masterrtc=2
    }
done

rm -v /tmp/${name}*
exit ${masterrtc}
