#!/bin/bash

set -eux

N=${1:-4}
MPWD=$(pwd)
echo "curdir=$MPWD"

docker network create \
  --driver=bridge \
  --subnet=172.77.0.0/16 \
  --ip-range=172.77.5.0/24 \
  --gateway=172.77.5.254 \
  newdagnet

for i in $(seq 1 $N)
do
    srcdir=$MPWD/server/node$i
    docker create --name=node$i --net=newdagnet --ip=172.77.5.$i newdag/node:latest \
    run --ip="172.77.5.$i" --node_port=1337 --proxy_port=1338 --service_port=80 --client_addr="172.77.5.$(($N+$i)):1339" 
    docker cp $srcdir/priv_key.pem       node$i:/
    docker cp $srcdir/newdag_wallet.dat  node$i:/
    docker start node$i

    docker run -d -p 1388$i:1388  --name=client$i --net=newdagnet --ip=172.77.5.$(($N+$i)) -it newdag/client:latest run\
    --client_port=1339 --proxy_addr="172.77.5.$i:1338" 
done

#docker run -it --rm --name=watcher --net=newdagnet --ip=172.77.5.99 newdag/watcher /watch.sh $N

#   ./server run --node_ip="192.168.137.129" --node_port=1337 --proxy_port=1338  --service_port=80 --client_addr="192.168.137.129:1339" 
#   ./client run --client_port=1339 --proxy_addr="192.168.137.129:1338" 
