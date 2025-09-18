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
stages="stp,entwta3 stp,testta3 lgkk,testta3 stplgk,test"
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
stages="stp,testta2 stp,entwta2"
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
