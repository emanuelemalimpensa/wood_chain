peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
	--tls \
	--cafile \
	"${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
	-C mychannel -n $1 \
	--peerAddresses localhost:7051 \
	--tlsRootCertFiles \
	"${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
	--peerAddresses localhost:9051 \
	--tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
	-c '{"function":"CreateAsset","Args":["asset1","1000","Trunk","Francia","Mild Forest","","100","100","90329013","2012-10-10"]}'
