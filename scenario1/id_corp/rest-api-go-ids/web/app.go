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
	// router.HandleFunc("/assets/{category}/{id:[a-zA-Z0-9]+}", setups.Query2).Methods("GET").Schemes("http").Name("PROVA")
	router.HandleFunc("/mock", setups.Query12).Methods("POST").Schemes("http").Name("mock insert")
	router.HandleFunc("/mock", setups.Query13).Methods("GET").Schemes("http").Name("mock fetch")
	fmt.Println("Listening (http://localhost:3000/)...")
	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:3000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
