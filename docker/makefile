default: help

images: bin node-image client-image watcher-image 

bin:
	CGO_ENABLED=0 go build -o ./server/server --ldflags '-extldflags "-static"'  ../cmd
	CGO_ENABLED=0 go build -o ./client/client --ldflags '-extldflags "-static"'  ../client

node-image:
	docker rmi -f newdag/node:latest
	docker build --no-cache=true -t newdag/node:latest  server/

client-image:
	docker rmi -f newdag/client:latest
	docker build --no-cache=true -t newdag/client:latest   client/

watcher-image:
	#docker rmi -f newdag/watcher:latest
	#docker build --no-cache=true -t newdag/watcher:latest watcher/

clean: 
	rm -rf ./server/server ./client/client
	./script/clean.sh

run: 
	./run.sh 

stop:
	./script/stop.sh

demo:
	./script/demo.sh 

conf: bin
	./script/conf.sh
	./server/server init
	mv newdag_wallet.dat  ./client/
	cp newdag_chain.db	  ./server/
	mv newdag_chain.db    ./client/

help:
	
	@echo "make images | node-image | client-image | watcher-image| run | demo | stop | clean | conf | bin"

.PHONY: all clean bin conf run stop demo help
