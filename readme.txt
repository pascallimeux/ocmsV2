sudo rm -R /tmp/

cd /opt/gopath/src/github.com/hyperledger/fabric-sdk-go/test/fixtures/
docker-compose -f docker-compose.yaml up --force-recreate -d
docker ps

cd /opt/gopath/src/github.com/pascallimeux/ocmsV2
go build ocms.go
./ocms
cd /opt/gopath/src/github.com/hyperledger/fabric-sdk-go/test/fixtures/
docker-compose -f docker-compose.yaml stop

