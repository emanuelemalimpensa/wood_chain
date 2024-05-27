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
	router.HandleFunc("/assets", setups.Query2).Methods("GET").Schemes("http").Name("ALL ASSETS")
	router.HandleFunc("/table", setups.Query4).Methods("GET").Schemes("http").Name("ALL TABLES")
	router.HandleFunc("/trunk", setups.Query3).Methods("GET").Schemes("http").Name("ALL TRUNKS")
	router.HandleFunc("/table/{id}", setups.Query5).Methods("GET").Schemes("http").Name("table id")
	router.HandleFunc("/trunk/{id}", setups.Query6).Methods("GET").Schemes("http").Name("trunk id")
	router.HandleFunc("/trunk/table", setups.Query7).Methods("POST").Schemes("http").Name("trunk id")
	router.HandleFunc("/trunk/{id}/newowner", setups.Query8).Methods("POST").Schemes("http").Name("change owner trunk")
	router.HandleFunc("/table/{id}/newowner", setups.Query9).Methods("POST").Schemes("http").Name("change owner table")
	router.HandleFunc("/history/trunk/{id}", setups.Query10).Methods("GET").Schemes("http").Name("histroy trunk id")
	router.HandleFunc("/history/table/{id}", setups.Query11).Methods("GET").Schemes("http").Name("history table id")
	router.HandleFunc("/audit", setups.Audit).Methods("POST").Schemes("http").Name("Audit Request")
	router.HandleFunc("/trunk/new/{CCname}", setups.Query12).Methods("POST").Schemes("http").Name("approved trunk creation")

	fmt.Println("Listening (http://localhost:3001/)...")
	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:3001",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
