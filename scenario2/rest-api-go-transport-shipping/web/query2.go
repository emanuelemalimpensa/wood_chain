package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type TableGen struct {
	Id     string
	Amount int
	SizeH  int
	SizeW  int
	SizeL  int
}

type NewOwner struct {
	NewOwn string
}

var id string = "1000"

// Query handles chaincode query requests.
func (setup OrgSetup) Query1(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST CREATE Transaction agreement")
	w.WriteHeader(http.StatusOK)
	queryParams := r.URL.Query()
	chainCodeName := queryParams.Get("chaincodeid")
	channelID := queryParams.Get("channelid")

	type Quantita struct {
		Quantita string `json:Quantita`
		Tipo     string `json:Tipo`
	}
	var prova Quantita
	err := json.NewDecoder(r.Body).Decode(&prova)

	args := r.URL.Query()["args"]
	fmt.Println("args   %T\n", args)
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	data := "2020-10-09"
	evaluateResponse, err := contract.SubmitTransaction("CreateBuffer", id, prova.Tipo, data, prova.Quantita)
	var id_int int
	id_int, _ = strconv.Atoi(id)
	id_int += 1
	id = strconv.Itoa(id_int)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	fmt.Fprintf(w, "Response: %s", evaluateResponse)
}

// Query handles chaincode query requests.
func (setup OrgSetup) Query2(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PUT MODIFY Transaction agreement BUYER INSERTION")
	vars := mux.Vars(r)

	w.WriteHeader(http.StatusOK)
	queryParams := r.URL.Query()
	chainCodeName := queryParams.Get("chaincodeid")
	channelID := queryParams.Get("channelid")

	args := r.URL.Query()["args"]
	role := "b"
	fmt.Println("args   %T\n", args)
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	evaluateResponse, err := contract.SubmitTransaction("PopulateAssetOrg", vars["org"], vars["id"], role)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	fmt.Fprintf(w, "Response: %s", evaluateResponse)
}

// Query handles chaincode query requests.
func (setup OrgSetup) Query3(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PUT MODIFY Transaction agreement BUYER INSERTION")
	vars := mux.Vars(r)

	w.WriteHeader(http.StatusOK)
	queryParams := r.URL.Query()
	chainCodeName := queryParams.Get("chaincodeid")
	channelID := queryParams.Get("channelid")

	args := r.URL.Query()["args"]
	role := "s"
	fmt.Println("args   %T\n", args)
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	evaluateResponse, err := contract.SubmitTransaction("PopulateAssetOrg", vars["org"], vars["id"], role)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	fmt.Fprintf(w, "Response: %s", evaluateResponse)
}
