
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


docker  rm $( docker ps -aq)
docker  images | awk '/vp|none|dev-/ { print $3}' | xargs docker rmi -f
rm -fr /var/hyperledger/production/*
rm -fr /home/blockchain/.fabric-ca-client/msp/

rm /opt/gopath/src/github.com/hyperledger/fabric-sdk-go/test/fixtures/fabricca/tlsOrg1/fabric-ca-server.db