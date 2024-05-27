package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// Query handles chaincode query requests.
func (setup OrgSetup) Query12(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CREATE A MOCK")
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Trunk: %v\n", vars["id"])
	type IDRelease struct {
		ID     string `json:ID`
		Issuer string `json:Issuer`
		Owner  string `json:Owner`
		Tipo   string `json:Tipo`
	}
	var prova IDRelease
	err := json.NewDecoder(r.Body).Decode(&prova)
	fmt.Printf("\n" + string(prova.ID))
	fmt.Printf("\n" + prova.Issuer)
	fmt.Printf("\n" + prova.Owner)
	fmt.Printf("\n" + prova.Tipo)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "errore nel parsing del body \n")
	}
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
	fmt.Fprintf(w, "Body: %v\n", prova)
	queryParams := r.URL.Query()
	chainCodeName := queryParams.Get("chaincodeid")
	channelID := queryParams.Get("channelid")
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	fmt.Fprintf(w, "Body: %v\n", prova)

	evaluateResponse, err := contract.SubmitTransaction("CreateAsset", prova.ID, prova.Issuer, prova.Owner, prova.Tipo)
	if err != nil {
		log.Fatal(w, "Error: %s", err)
		os.Exit(1)
	}
	fmt.Fprintf(w, "Response: %s", evaluateResponse)
}

func (setup OrgSetup) Query13(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET ALL MOCK ASSETS")
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Asset: %v\n", vars["id"])
	queryParams := r.URL.Query()
	chainCodeName := queryParams.Get("chaincodeid")
	channelID := queryParams.Get("channelid")
	args := r.URL.Query()["args"]
	fmt.Println("args   %T\n", args)
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	evaluateResponse, err := contract.EvaluateTransaction("GetAllAssets", args...)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	fmt.Fprintf(w, "Response: %s", evaluateResponse)
}
