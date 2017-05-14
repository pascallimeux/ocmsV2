
git clone https://gerrit.hyperledger.org/r/fabric-sdk-go
git clone https://github.com/pascallimeux/ocmsV2.git
go get -u github.com/kardianos/govendor
rm -R vendor
govendor init
govendor add +external



sudo rm -R /tmp/

cd /opt/gopath/src/github.com/hyperledger/fabric-sdk-go/test/fixtures/
docker-compose -f docker-compose.yaml up --force-recreate -d
docker ps

cd /opt/gopath/src/github.com/pascallimeux/ocmsV2
go build ocms.go
./ocms
cd /opt/gopath/src/github.com/hyperledger/fabric-sdk-go/test/fixtures/
docker-compose -f docker-compose.yaml stop

go get github.com/gorilla/mux
go get github.com/op/go-logging


how to clean environment
remove containers
docker ps -a
sudo docker  rm $( docker ps -aq)

remove images
docker images
sudo docker  images | awk '/vp|none|dev-/ { print $3}' | xargs docker rmi -f

rm -fr /var/hyperledger/production/*
rm -fr /home/blockchain/.fabric-ca-client/msp/

rm /opt/gopath/src/github.com/hyperledger/fabric-sdk-go/test/fixtures/fabricca/tlsOrg1/fabric-ca-server.db


how to create docker image
go build ocms.go
docker build -t ocmsv2 .

verify
docker images

start
docker run -d -p 8000:8000 --name ocms ocmsv2

se connecter au docker
docker exec -it ocms bash

stopper le container
docker kill ocms

effacer le container
docker rm ocms

effacer l'images
docker rmi ocmsv2