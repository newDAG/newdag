#!/bin/bash

N=${1:-4}
DEST=${2:-"server"}
IPBASE=${3:-172.77.5.}
PORT=${4:-1337}
if  [ $N -lt 2 ]
then
   echo "param 1 should > 1"
   exit
fi

set -e

for i in $(seq 1 $N)
do
    dest=$DEST/node$i
    mkdir -p $dest
    ./server/server keygen | sed -n -e "2 w $dest/pub" -e "4,+4 w $dest/priv_key.pem"
    ./client/client wgen
    mv newdag_wallet.dat $dest/
    echo "$IPBASE$i:$PORT" > $dest/addr
done

PFILE=$DEST/peers.json
echo "[" > $PFILE
for i in $(seq 1 $N)
do
    com=","
    if [[ $i == $N ]]; then
        com=""
    fi

    dest=$DEST/node$i
    printf "\t{\n" >> $PFILE
    printf "\t\t\"NetAddr\":\"$(cat $dest/addr)\",\n" >> $PFILE
    printf "\t\t\"PubKeyHex\":\"$(cat $dest/pub)\"\n" >> $PFILE
    printf "\t}%s\n"  $com >> $PFILE
    rm $dest/addr $dest/pub

done
echo "]" >> $PFILE

echo "generate successed!"
