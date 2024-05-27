package chaincode

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
// Insert struct field in alphabetic order => to achieve determinism across languages
// golang keeps the order when marshal to json but doesn't order automatically
type Trunk struct {
	AssetType  string  `json:AssetType`
	ID         string  `json:"ID"`
	Nodi       int     `json:Nodi`
	Tipo       string  `json:Tipo`
	Origine    Origine `json:Origine`
	Owner      string  `json:"Owner"`
	SizeH      int     `json:"H"`
	SizeL      int     `json:"L"`
	Peso       int     `json:Peso`
	DataTaglio string  `json:DataTaglio`
}

type Origine struct {
	Nazione string `json:Nazione`
	Azienda string `json:Azienda`
}

// Asset describes basic details of what makes up a simple asset
// Insert struct field in alphabetic order => to achieve determinism across languages
// golang keeps the order when marshal to json but doesn't order automatically
type Table struct {
	AssetType     string  `json:AssetType`
	ID            string  `json:"ID"`
	Tipo          string  `json:Tipo`
	Origine       Origine `json:Origine`
	OrigineTronco Origine `json:OrigineTronco`
	IDtronco      string  `json:ID`
	Owner         string  `json:"Owner"`
	SizeH         int     `json:"H"`
	SizeW         int     `json:"W"`
	SizeL         int     `json:"L"`
	Peso          int     `json:Peso`
	DataTaglio    string  `json:DataTaglio`
}

type HistoryQueryResult struct {
	Record   *Trunk `json:"record"`
	TxId     string `json:"txId"`
	IsDelete bool   `json:"isDelete"`
}

type HistoryQueryResultTable struct {
	Record   *Table `json:"record"`
	TxId     string `json:"txId"`
	IsDelete bool   `json:"isDelete"`
}

type IDRelease struct {
	ID     string `json:ID`
	Issuer string `json:Issuer`
	Owner  string `json:Owner`
	Tipo   string `json:WoodType`
	Free   string `json:FreeState`
}

type IDret struct {
	ID string `json:IDbuff`
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, ID string, Nodi string, Tipo string, NazioneOrigine string, AziendaOrigine string, Owner string, SizeH string, SizeL string, Peso string, DataTaglio string, CCName string) error {
	//exists, err := s.AssetExists(ctx, id)

	exists, err := TrunkExists(ctx, ID)

	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", ID)
	}

	sz := len(ID)
	params := []string{"GetAllAssets"}
	queryArgs := make([][]byte, len(params))
	for i, arg := range params {
		queryArgs[i] = []byte(arg)
	}

	response := ctx.GetStub().InvokeChaincode(CCName, queryArgs, "mychannel")
	org, _ := ctx.GetClientIdentity().GetMSPID()
	json.Marshal(response)
	var ids []*IDRelease
	err = json.Unmarshal(response.Payload, &ids)
	if err != nil {
		return err
	}
	var i int
	for i = 0; i < len(ids); i++ {
		if ids[i].ID == ID[0:sz-1] {
			if ids[i].Owner != org {
				return fmt.Errorf("ID IS NOT ALLOWED, ALREADY BOOKED")
			}
		}
	}

	PesoInsert, _ := strconv.Atoi(Peso)
	NodiInsert, _ := strconv.Atoi(Nodi)
	SHInsert, _ := strconv.Atoi(SizeH)
	SLInsert, _ := strconv.Atoi(SizeL)

	if Owner == "" {
		Owner, _ = ctx.GetClientIdentity().GetMSPID()
	}

	trunk := Trunk{
		ID:         ID,
		Peso:       PesoInsert,
		DataTaglio: DataTaglio,
		Nodi:       NodiInsert,
		Tipo:       Tipo,
		SizeH:      SHInsert,
		SizeL:      SLInsert,
		Owner:      Owner,
	}
	trunk.Origine.Azienda = AziendaOrigine
	trunk.Origine.Nazione = NazioneOrigine
	trunk.AssetType = "Trunk"
	trunkJSON, err := json.Marshal(trunk)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(ID, trunkJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, ID string) error {
	exists, err := TrunkExists(ctx, ID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", ID)
	}

	t, _ := s.ReadAsset(ctx, ID)
	org, _ := ctx.GetClientIdentity().GetMSPID()
	if t.Owner != org {
		return fmt.Errorf("ORG NOT ALLOWED")
	}

	return ctx.GetStub().DelState(ID)
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Trunk, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var trunks []*Trunk
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var trunk Trunk
		err = json.Unmarshal(queryResponse.Value, &trunk)
		if err != nil {
			return nil, err
		}

		trunks = append(trunks, &trunk)
	}

	return trunks, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssetsTrunk(ctx contractapi.TransactionContextInterface) ([]*Trunk, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var trunks []*Trunk
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var trunk Trunk
		err = json.Unmarshal(queryResponse.Value, &trunk)
		if err != nil {
			return nil, err
		}
		if trunk.AssetType == "Trunk" {
			trunks = append(trunks, &trunk)
		}
	}

	return trunks, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssetsTable(ctx contractapi.TransactionContextInterface) ([]*Trunk, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var trunks []*Trunk
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var trunk Trunk
		err = json.Unmarshal(queryResponse.Value, &trunk)
		if err != nil {
			return nil, err
		}
		if trunk.AssetType == "Table" {
			trunks = append(trunks, &trunk)
		}
	}

	return trunks, nil
}

func (t *SmartContract) GetAssetHistory(ctx contractapi.TransactionContextInterface, assetID string) ([]HistoryQueryResult, error) {

	//https://github.com/hyperledger/fabric-samples/blob/8ca50df4ffec311e59451c2a7ebe210d9e6f0004/asset-transfer-ledger-queries/chaincode-go/asset_transfer_ledger_chaincode.go

	log.Printf("GetAssetHistory: ID %v", assetID)

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(assetID)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Trunk
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &asset)
			if err != nil {
				return nil, err
			}
		} else {
			asset = Trunk{
				ID: assetID,
			}
		}

		record := HistoryQueryResult{
			TxId:     response.TxId,
			Record:   &asset,
			IsDelete: response.IsDelete,
		}
		records = append(records, record)
	}

	return records, nil
}

func (t *SmartContract) GetAssetHistoryTable(ctx contractapi.TransactionContextInterface, assetID string) ([]HistoryQueryResultTable, error) {

	//https://github.com/hyperledger/fabric-samples/blob/8ca50df4ffec311e59451c2a7ebe210d9e6f0004/asset-transfer-ledger-queries/chaincode-go/asset_transfer_ledger_chaincode.go

	log.Printf("GetAssetHistoryTable: ID %v", assetID)

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(assetID)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []HistoryQueryResultTable
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Table
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &asset)
			if err != nil {
				return nil, err
			}
		} else {
			asset = Table{
				ID: assetID,
			}
		}

		record := HistoryQueryResultTable{
			TxId:     response.TxId,
			Record:   &asset,
			IsDelete: response.IsDelete,
		}
		records = append(records, record)
	}

	return records, nil
}

func (s *SmartContract) GenerateTablesFromTrunk(ctx contractapi.TransactionContextInterface, ID string, amountS string, sizeHS string, sizeWS string, sizeLS string, Nazione string, Azienda string) ([]Table, error) {
	t, _ := s.ReadAsset(ctx, ID)
	org, _ := ctx.GetClientIdentity().GetMSPID()
	if t.Owner != org {
		return nil, fmt.Errorf("ORG NOT ALLOWED")
	}

	log.Printf("Tables generation, default id generation")
	trunk, err := s.ReadAsset(ctx, ID)
	amount, _ := strconv.Atoi(amountS)
	sizeH, _ := strconv.Atoi(sizeHS)
	sizeL, _ := strconv.Atoi(sizeLS)
	sizeW, _ := strconv.Atoi(sizeWS)
	if err != nil {
		return nil, err
	}

	var Origine Origine
	var tables []Table
	Origine.Azienda = Azienda
	Origine.Nazione = Nazione
	for i := 1; i <= amount; i++ {
		// now := time.Now()
		// y, m, d := now.Date()
		t := time.Now()
		year := t.Year()            // type int
		month := t.Month().String() // type time.Month
		day := t.Day()              // type int
		today := strconv.Itoa(year) + "-" + month + "-" + strconv.Itoa(day)

		table := Table{
			ID:         trunk.ID + "_" + strconv.Itoa(i),
			Peso:       trunk.Peso / amount,
			SizeH:      sizeH,
			SizeW:      sizeW,
			SizeL:      sizeL,
			Origine:    Origine,
			IDtronco:   trunk.ID,
			Owner:      Azienda,
			DataTaglio: today,
			Tipo:       trunk.Tipo,
			AssetType:  "Table",
		}

		trunk, _ := ctx.GetStub().GetHistoryForKey(table.IDtronco)
		for trunk.HasNext() {
			curr_trunk, err := trunk.Next()
			if err != nil {
				return nil, err
			}
			if !curr_trunk.IsDelete {
				var trunk Trunk
				//asset.OrigineTronco.Nazione
				json.Unmarshal(curr_trunk.GetValue(), &trunk)
				table.OrigineTronco.Nazione = trunk.Origine.Nazione
				table.OrigineTronco.Azienda = trunk.Origine.Azienda
				// asset.OrigineTronco.Azienda = curr_trunk.record.Origine.Azienda
			}
		}

		tableJSON, err := json.Marshal(table)
		if err != nil {
			return nil, err
		}

		err = ctx.GetStub().PutState(table.ID, tableJSON)
		if err != nil {
			return nil, err
		}

		tables = append(tables, table)
	}

	err = s.DeleteAsset(ctx, ID)

	if err != nil {
		return nil, err
	}

	return tables, nil
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	org, _ := ctx.GetClientIdentity().GetMSPID()
	trunks := []Trunk{
		{ID: "asset1", Nodi: 20, Tipo: "Tek", Origine: Origine{Nazione: "IT", Azienda: "MultiForest Company"}, SizeH: 5, SizeL: 5, Peso: 200, Owner: org, DataTaglio: "2020-01-01", AssetType: "Trunk"},
		{ID: "asset2", Nodi: 20, Tipo: "Tek", Origine: Origine{Nazione: "IT", Azienda: "MultiForest Company"}, SizeH: 5, SizeL: 5, Peso: 200, Owner: org, DataTaglio: "2020-01-01", AssetType: "Trunk"},
		{ID: "asset3", Nodi: 20, Tipo: "Tek", Origine: Origine{Nazione: "IT", Azienda: "MultiForest Company"}, SizeH: 5, SizeL: 5, Peso: 200, Owner: org, DataTaglio: "2020-01-01", AssetType: "Trunk"},
		{ID: "asset4", Nodi: 20, Tipo: "Tek", Origine: Origine{Nazione: "IT", Azienda: "MultiForest Company"}, SizeH: 5, SizeL: 5, Peso: 200, Owner: org, DataTaglio: "2020-01-01", AssetType: "Trunk"},
		{ID: "asset5", Nodi: 20, Tipo: "Tek", Origine: Origine{Nazione: "Italia", Azienda: "MultiForest Company"}, SizeH: 5, SizeL: 5, Peso: 200, Owner: org, DataTaglio: "2020-01-01", AssetType: "Trunk"},
		{ID: "asset6", Nodi: 20, Tipo: "Tek", Origine: Origine{Nazione: "Italia", Azienda: "MultiForest Company"}, SizeH: 5, SizeL: 5, Peso: 200, Owner: org, DataTaglio: "2020-01-01", AssetType: "Trunk"},
		{ID: "asset7", Nodi: 20, Tipo: "Tek", Origine: Origine{Nazione: "Italia", Azienda: "MultiForest Company"}, SizeH: 5, SizeL: 5, Peso: 200, Owner: org, DataTaglio: "2020-01-01", AssetType: "Trunk"},
		{ID: "asset8", Nodi: 20, Tipo: "Tek", Origine: Origine{Nazione: "Italia", Azienda: "MultiForest Company"}, SizeH: 5, SizeL: 5, Peso: 200, Owner: org, DataTaglio: "2020-01-01", AssetType: "Trunk"},
		{ID: "asset9", Nodi: 20, Tipo: "Tek", Origine: Origine{Nazione: "Italia", Azienda: "MultiForest Company"}, SizeH: 5, SizeL: 5, Peso: 200, Owner: org, DataTaglio: "2020-01-01", AssetType: "Trunk"},
		{ID: "asset10", Nodi: 20, Tipo: "Tek", Origine: Origine{Nazione: "Italia", Azienda: "MultiForest Company"}, SizeH: 5, SizeL: 5, Peso: 200, Owner: org, DataTaglio: "2020-01-01", AssetType: "Trunk"},
	}

	for _, trunk := range trunks {
		trunkJSON, err := json.Marshal(trunk)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(trunk.ID, trunkJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Trunk, error) {
	trunkJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if trunkJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var trunk Trunk
	err = json.Unmarshal(trunkJSON, &trunk)
	if err != nil {
		return nil, err
	}

	if trunk.AssetType != "Trunk" {
		return nil, fmt.Errorf("not trunk type")
	}

	return &trunk, nil
}

// ReadTable returns the asset stored in the world state with given id.
func (s *SmartContract) ReadTable(ctx contractapi.TransactionContextInterface, id string) (*Table, error) {
	tableJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if tableJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var table Table
	err = json.Unmarshal(tableJSON, &table)
	if err != nil {
		return nil, err
	}

	if table.AssetType != "Table" {
		return nil, fmt.Errorf("not table type")
	}
	return &table, nil
}

// TransferAsset updates the owner field of asset with given id in world state, and returns the old owner.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) (string, error) {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return "", err
	}

	org, _ := ctx.GetClientIdentity().GetMSPID()
	oldOwner := asset.Owner
	if oldOwner == org {
		asset.Owner = newOwner
	} else {
		return "", fmt.Errorf("ORG NOT ALLOWED")
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(id, assetJSON)
	if err != nil {
		return "", err
	}

	return oldOwner, nil
}

// TransferAsset updates the owner field of asset with given id in world state, and returns the old owner.
func (s *SmartContract) TransferTable(ctx contractapi.TransactionContextInterface, id string, newOwner string) (string, error) {
	asset, err := s.ReadTable(ctx, id)
	org, _ := ctx.GetClientIdentity().GetMSPID()
	oldOwner := asset.Owner
	if oldOwner == org {
		asset.Owner = newOwner
	} else {
		return "", fmt.Errorf("ORG NOT ALLOWED")
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(id, assetJSON)
	if err != nil {
		return "", err
	}

	return oldOwner, nil
}

// TrunkExists returns true when asset with given ID exists in world state
func TrunkExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	trunkJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return trunkJSON != nil, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, ID string, Nodi string, Tipo string, NazioneOrigine string, AziendaOrigine string, Owner string, SizeH string, SizeL string, Peso string, DataTaglio string) error {
	t, _ := s.ReadAsset(ctx, ID)
	org, _ := ctx.GetClientIdentity().GetMSPID()
	if t.Owner != org {
		return fmt.Errorf("ORG NOT ALLOWED")
	}
	exists, err := TrunkExists(ctx, ID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", ID)
	}

	PesoInsert, _ := strconv.Atoi(Peso)
	NodiInsert, _ := strconv.Atoi(Nodi)
	SHInsert, _ := strconv.Atoi(SizeH)
	SLInsert, _ := strconv.Atoi(SizeL)

	trunk := Trunk{
		ID:         ID,
		Peso:       PesoInsert,
		DataTaglio: DataTaglio,
		Nodi:       NodiInsert,
		Tipo:       Tipo,
		SizeH:      SHInsert,
		SizeL:      SLInsert,
		Owner:      Owner,
	}
	trunk.AssetType = "Trunk"
	trunk.Origine.Azienda = AziendaOrigine
	trunk.Origine.Nazione = NazioneOrigine

	trunkJSON, err := json.Marshal(trunk)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(ID, trunkJSON)
}

func (s *SmartContract) GenerateApprovedTrunk(ctx contractapi.TransactionContextInterface, Nodi string, Tipo string, NazioneOrigine string, AziendaOrigine string, SizeH string, SizeL string, Peso string, DataTaglio string, CCName string) (*IDRelease, error) {
	org, _ := ctx.GetClientIdentity().GetMSPID() //Org1MSP
	params := []string{"GetAllFreeIDsForOrg", org}
	queryArgs := make([][]byte, len(params))
	for i, arg := range params {
		queryArgs[i] = []byte(arg)
	}

	response := ctx.GetStub().InvokeChaincode(CCName, queryArgs, "mychannel")

	json.Marshal(response)
	var ids []*IDRelease
	err := json.Unmarshal(response.Payload, &ids)
	if err != nil {
		return nil, err
	}
	if len(ids) > 0 {
		for i, _ := range ids {
			if ids[i].Tipo == Tipo {
				err = s.CreateAsset(ctx, ids[i].ID+"_", Nodi, Tipo, NazioneOrigine, AziendaOrigine, ids[i].Owner, SizeH, SizeL, Peso, DataTaglio, CCName)
				if err != nil {
					return nil, err
				}
				params = []string{"SetAssetAsUsed", ids[i].ID}
				queryArgs = make([][]byte, len(params))
				for i, arg := range params {
					queryArgs[i] = []byte(arg)
				}
				response = ctx.GetStub().InvokeChaincode(CCName, queryArgs, "mychannel")
				return ids[i], nil
			}
		}
		return nil, fmt.Errorf("0 FREE IDs FOR SPECIFIC TRUNK TYPE ")
	} else {
		return nil, fmt.Errorf("0 FREE IDs FOR YOUR ORG [%s]", org)
	}
}

func (s *SmartContract) GetIDOfOwnedTypeTrunks(ctx contractapi.TransactionContextInterface, Tipo string) ([]*IDret, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	Org, _ := ctx.GetClientIdentity().GetMSPID()
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	var trunksID []*IDret
	// org, _ := ctx.GetClientIdentity().GetMSPID()
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var trunk Trunk
		var CurrId IDret
		err = json.Unmarshal(queryResponse.Value, &trunk)
		if err != nil {
			return nil, err
		}
		if trunk.AssetType == "Trunk" {
			if trunk.Tipo == Tipo && trunk.Owner == Org {
				// return &trunk.ID, nil
				CurrId.ID = trunk.ID
				trunksID = append(trunksID, &CurrId)
			}
		}
	}
	if len(trunksID) == 0 {
		return nil, fmt.Errorf("NOT TRUNKS AVAIABLE")
	}
	return trunksID, nil
}

func (s *SmartContract) ChangeOwnerIDs(ctx contractapi.TransactionContextInterface, NewOwner string, Ids []string) ([]string, error) {
	var id int
	org, _ := ctx.GetClientIdentity().GetMSPID()
	for id = 0; id < len(Ids); id++ {
		old, err := s.TransferAsset(ctx, Ids[id], NewOwner)
		if err != nil {
			s.TransferAsset(ctx, Ids[id], org)
			return nil, err
		}
		if old != org {
			s.TransferAsset(ctx, Ids[id], org)
			return nil, fmt.Errorf("Old owner is not the current one, keeping ownership unchanged for ID :: %s", Ids[id])
		}
	}
	return Ids, nil
}
