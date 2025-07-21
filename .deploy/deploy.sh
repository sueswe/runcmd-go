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

stages="stp,testta3 lgkk,testta3 stplgk,test"

for UMG in ${stages}
do
    cd /tmp/ || exit 1
    "$HOME"/bin/vicecersa.sh "${UMG}" runcmd \$HOME/bin/ || {
        echo "Status: $?"
        exit 2
    }
done


echo '
### AIX #########################################################
'

# for UMG in ${stages}
# do
#     cd /tmp/ || exit 1
#     "$HOME"/bin/vicecersa.sh "${UMG}" runcmd.aix \$HOME/bin/ runcmd || {
#         echo "Status: $?"
#         exit 2
#     }
# done
