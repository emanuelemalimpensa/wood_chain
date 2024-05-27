package main

import (
	"fmt"
	"rest-api-go/web"
)

func main() {
	//Initialize setup for Org1
	cryptoPath := "/home/manu/Desktop/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com"
	orgConfig := web.OrgSetup{
		OrgName:      "Org2",
		MSPID:        "Org2MSP",
		CertPath:     cryptoPath + "/users/User1@org2.example.com/msp/signcerts/User1@org2.example.com-cert.pem",
		KeyPath:      cryptoPath + "/users/User1@org2.example.com/msp/keystore/",
		TLSCertPath:  cryptoPath + "/peers/peer0.org2.example.com/tls/ca.crt",
		PeerEndpoint: "localhost:9051",
		GatewayPeer:  "peer0.org2.example.com",
	}

	orgSetup, err := web.Initialize(orgConfig)
	if err != nil {
		fmt.Println("Error initializing setup for Org2: ", err)
	}
	web.Serve(web.OrgSetup(*orgSetup))
}
