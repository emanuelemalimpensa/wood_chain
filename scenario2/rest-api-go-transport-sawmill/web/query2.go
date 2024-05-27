package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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

// Query handles chaincode query requests.
func (setup OrgSetup) Query1(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET MODIFY Transaction agreement")
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Asset: %v as %v\n", vars["id"])
	queryParams := r.URL.Query()
	chainCodeName := queryParams.Get("chaincodeid")
	channelID := queryParams.Get("channelid")

	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	org := "Org2MSP"
	evaluateResponse, err := contract.SubmitTransaction("ApproveTransferForMyOrg", vars["id"], org)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	fmt.Fprintf(w, "Response: %s", evaluateResponse)
}

func (setup OrgSetup) Query2(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SELL REQUEST")
	type Quantita struct {
		Quantita string `json:Quantita`
		Tipo     string `json:Tipo`
	}
	var prova Quantita
	err := json.NewDecoder(r.Body).Decode(&prova)
	if err != nil {
		return
	}

	queryParams := r.URL.Query()
	chainCodeName := queryParams.Get("chaincodeid")
	channelID := queryParams.Get("channelid")

	myurl := "http://localhost:3001/sell?channelid=" + channelID + "&chaincodeid=" + chainCodeName
	bodyVal := fmt.Sprintf("{\"Quantita\":\"%s\",\"Tipo\":\"%s\"}", prova.Quantita, prova.Tipo)
	body := strings.NewReader(bodyVal)
	_, err2 := http.Post(myurl, "application/json", body)
	if err2 != nil {
		log.Fatal(err2)
		os.Exit(1)
	}
}
