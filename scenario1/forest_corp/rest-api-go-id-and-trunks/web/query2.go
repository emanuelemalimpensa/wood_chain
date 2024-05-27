package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

type Trunk struct {
	Nodi       string `json:Nodi`
	Tipo       string `json:Tipo`
	Nazione    string `json:Nazione`
	Azienda    string `json:Azienda`
	SizeH      string `json:"H"`
	SizeL      string `json:"L"`
	Peso       string `json:Peso`
	DataTaglio string `json:DataTaglio`
}

type NewOwner struct {
	NewOwn string
}

// Query handles chaincode query requests.
func (setup OrgSetup) Query2(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET ALL ASSETS")
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

// Query handles chaincode query requests.
func (setup OrgSetup) Query3(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET ALL TRUNKS")
	w.WriteHeader(http.StatusOK)
	queryParams := r.URL.Query()
	chainCodeName := queryParams.Get("chaincodeid")
	channelID := queryParams.Get("channelid")

	args := r.URL.Query()["args"]
	fmt.Println("args   %T\n", args)
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	evaluateResponse, err := contract.EvaluateTransaction("GetAllAssetsTrunk", args...)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	fmt.Fprintf(w, "Response: %s", evaluateResponse)
}

// Query handles chaincode query requests.
func (setup OrgSetup) Query4(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET ALL TABLES")
	w.WriteHeader(http.StatusOK)
	queryParams := r.URL.Query()
	chainCodeName := queryParams.Get("chaincodeid")
	channelID := queryParams.Get("channelid")

	args := r.URL.Query()["args"]
	fmt.Println("args   %T\n", args)
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	evaluateResponse, err := contract.EvaluateTransaction("GetAllAssetsTable", args...)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	fmt.Fprintf(w, "Response: %s", evaluateResponse)
}

// Query handles chaincode query requests.
func (setup OrgSetup) Query5(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET TABLE\n")
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Table : %v\n", vars["id"])
	queryParams := r.URL.Query()
	chainCodeName := queryParams.Get("chaincodeid")
	channelID := queryParams.Get("channelid")

	args := r.URL.Query()["args"]
	args = append(args, vars["id"])
	fmt.Println("args   %T\n", args)
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	evaluateResponse, err := contract.EvaluateTransaction("ReadTable", args...)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	fmt.Fprintf(w, "Response: %s", evaluateResponse)
}

// Query handles chaincode query requests.
func (setup OrgSetup) Query6(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET TRUNK\n")
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Trunk: %v\n", vars["id"])
	queryParams := r.URL.Query()
	chainCodeName := queryParams.Get("chaincodeid")
	channelID := queryParams.Get("channelid")

	fmt.Println("args   %T\n", vars["id"])
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, vars["id"])
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	evaluateResponse, err := contract.EvaluateTransaction("ReadAsset", vars["id"])
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	fmt.Fprintf(w, "Response: %s", evaluateResponse)
}

func (setup OrgSetup) Query7(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CREATE TABLES FROM TRUNK\n")
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Trunk: %v\n", vars["id"])
	var prova TableGen
	err := json.NewDecoder(r.Body).Decode(&prova)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "errore nel parsing del body \n")
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
	fmt.Fprintf(w, "Body: %v\n", prova)
	queryParams := r.URL.Query()
	chainCodeName := queryParams.Get("chaincodeid")
	channelID := queryParams.Get("channelid")

	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, vars["id"])
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	fmt.Fprintf(w, "Body: %v\n", prova)
	evaluateResponse, err := contract.SubmitTransaction("GenerateTablesWithoutIdFromTrunk",
		prova.Id,
		strconv.Itoa(prova.Amount),
		strconv.Itoa(prova.SizeH),
		strconv.Itoa(prova.SizeW),
		strconv.Itoa(prova.SizeL))

	fmt.Fprintf(w, "Response: %s", evaluateResponse)
}

func (setup OrgSetup) Query8(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Change owner of a trunk\n")
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Trunk: %v\n", vars["id"])
	var prova NewOwner
	err := json.NewDecoder(r.Body).Decode(&prova)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "errore nel parsing del body \n")
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
	fmt.Fprintf(w, "Body: %v\n", prova)
	queryParams := r.URL.Query()
	chainCodeName := queryParams.Get("chaincodeid")
	channelID := queryParams.Get("channelid")

	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, vars["id"])
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	fmt.Fprintf(w, "Body: %v\n", prova)
	evaluateResponse, err := contract.SubmitTransaction("TransferAsset", vars["id"], prova.NewOwn)

	fmt.Fprintf(w, "Response: %s", evaluateResponse)
}

func (setup OrgSetup) Query9(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Change owner of a trunk\n")
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Trunk: %v\n", vars["id"])
	var prova NewOwner
	err := json.NewDecoder(r.Body).Decode(&prova)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "errore nel parsing del body \n")
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
	fmt.Fprintf(w, "Body: %v\n", prova)
	queryParams := r.URL.Query()
	chainCodeName := queryParams.Get("chaincodeid")
	channelID := queryParams.Get("channelid")

	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, vars["id"])
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	fmt.Fprintf(w, "Body: %v\n", prova)
	evaluateResponse, err := contract.SubmitTransaction("TransferTable", vars["id"], prova.NewOwn)

	fmt.Fprintf(w, "Response: %s", evaluateResponse)
}

func (setup OrgSetup) Query10(w http.ResponseWriter, r *http.Request) {
	fmt.Println("histroy of a trunk\n")
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Trunk: %v\n", vars["id"])
	w.WriteHeader(http.StatusOK)
	queryParams := r.URL.Query()
	chainCodeName := queryParams.Get("chaincodeid")
	channelID := queryParams.Get("channelid")

	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, vars["id"])
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	evaluateResponse, _ := contract.EvaluateTransaction("GetAssetHistory", vars["id"])

	fmt.Fprintf(w, "Response: %s", evaluateResponse)
}

func (setup OrgSetup) Query11(w http.ResponseWriter, r *http.Request) {
	fmt.Println("history of a table\n")
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Table: %v\n", vars["id"])
	w.WriteHeader(http.StatusOK)
	queryParams := r.URL.Query()
	chainCodeName := queryParams.Get("chaincodeid")
	channelID := queryParams.Get("channelid")

	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, vars["id"])
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	evaluateResponse, _ := contract.EvaluateTransaction("GetAssetHistoryTable", vars["id"])

	fmt.Fprintf(w, "Response: %s", evaluateResponse)
}

func (setup OrgSetup) Audit(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AUDIT REQUEST SUBMISSION")

	type Quantita struct {
		Quantita   string `json:Quantita`
		ReqCompany string `json:ReqCompany`
		Tipo       string `json:Tipo`
	}
	var prova Quantita
	json.NewDecoder(r.Body).Decode(&prova)
	//GENERAZIONE DEGLI ID RICHIESTI CON LA AUDIT REQUEST
	myurl := "http://localhost:27000/api/product/audit"

	bodyVal := fmt.Sprintf("{\"Quantita\":\"%s\",\"ReqCompany\":\"%s\",\"Tipo\":\"%s\"}", prova.Quantita, prova.ReqCompany, prova.Tipo)
	body := strings.NewReader(bodyVal)
	http.Post(myurl, "application/json", body)

}

func (setup OrgSetup) Query12(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GENERATE APPROVED TRUNK\n")
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	var t Trunk
	err := json.NewDecoder(r.Body).Decode(&t)

	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "errore nel parsing del body \n")
	}
	fmt.Fprintf(w, "Body: %v\n", t)
	queryParams := r.URL.Query()
	chainCodeName := queryParams.Get("chaincodeid")
	channelID := queryParams.Get("channelid")

	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, vars["id"])
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	fmt.Fprintf(w, "Body: %v\n", t)
	evaluateResponse, err := contract.SubmitTransaction("GenerateApprovedTrunk", t.Nodi, t.Tipo, t.Nazione, t.Azienda, t.SizeH, t.SizeL, t.Peso, t.DataTaglio, vars["CCname"])

	fmt.Fprintf(w, "Response: %s", evaluateResponse)
}
