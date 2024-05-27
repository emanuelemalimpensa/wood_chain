package web

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// OrgSetup contains organization's config to interact with the network.
type OrgSetup struct {
	OrgName      string
	MSPID        string
	CryptoPath   string
	CertPath     string
	KeyPath      string
	TLSCertPath  string
	PeerEndpoint string
	GatewayPeer  string
	Gateway      client.Gateway
}

// Serve starts http web server.
func Serve(setups OrgSetup) {
	router := mux.NewRouter()
	router.HandleFunc("/sell/new", setups.Query1).Methods("POST").Schemes("http").Name("sell approval for current org")
	router.HandleFunc("/buyer/{id}/{org}", setups.Query2).Methods("PUT").Schemes("http").Name("set buyer approval")
	router.HandleFunc("/seller/{id}/{org}", setups.Query3).Methods("PUT").Schemes("http").Name("sell seller approval")

	fmt.Println("Listening (http://localhost:3003/)...")
	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:3003",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
