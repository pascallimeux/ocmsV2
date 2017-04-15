get consent from github
wget https://raw.githubusercontent.com/pascallimeux/consentSC/master/consentv2/consentv2.go 
wget https://raw.githubusercontent.com/pascallimeux/consentSC/master/consentv2/consentv2_test.go
go test
chmod 755 ../consentv2/
chmod 644 *.* 

