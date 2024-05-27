package chaincode

import (
	"encoding/json"
	"fmt"

	// "log"
	// "strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

type IDRelease struct {
	ID     string `json:ID`
	Issuer string `json:Issuer`
	Owner  string `json:Owner`
	Tipo   string `json:WoodType`
	Free   string `json:FreeState`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// 	ids := []IDRelease{{IDs: []string{"0001", "0002", "0003", "0004", "0005", "0006", "0007", "0008", "0009", "0010", "0011", "0012", "0013", "0014", "0015", "0016", "0017", "0018", "0019", "0020"}, ID: "asset1", Quantita: "20", Issuer: "EmitterCompany", Owner: "Forest", Tipo: "Tek"},
	// 		{IDs: []string{"1001", "1002", "1003", "1004", "1005", "1006", "1007", "1008", "1009", "1010", "1011", "1012", "1013", "1014", "1015", "1016", "1017", "1018", "1019", "1020"}, ID: "asset2", Quantita: "20", Issuer: "EmitterCompany", Owner: "MultiForest Company2", Tipo: "Tek"}}
	ids := []IDRelease{{ID: "asset1", Issuer: "EmitterCompany", Owner: "Org1MSP", Tipo: "Tek", Free: "true"},
		{ID: "asset3", Issuer: "EmitterCompany", Owner: "Org1MSP", Tipo: "Faggio", Free: "true"},
		{ID: "asset4", Issuer: "EmitterCompany", Owner: "Org1MSP", Tipo: "Faggio", Free: "true"},
		{ID: "asset5", Issuer: "EmitterCompany", Owner: "Org1MSP", Tipo: "Faggio", Free: "true"},
		{ID: "asset2", Issuer: "EmitterCompany", Owner: "Org2MSP", Tipo: "Tek", Free: "true"}}

	for _, id := range ids {
		idJSON, err := json.Marshal(id)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(id.ID, idJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*IDRelease, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var ids []*IDRelease
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var id IDRelease
		err = json.Unmarshal(queryResponse.Value, &id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, &id)

	}
	return ids, nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, ID string, Issuer string, Owner string, Tipo string) error {
	//exists, err := s.AssetExists(ctx, id)

	exists, err := ContractExists(ctx, ID)

	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", ID)
	}

	contratto := IDRelease{
		ID:     ID,
		Tipo:   Tipo,
		Owner:  Owner,
		Issuer: Issuer,
		Free:   "true",
	}
	contractJSON, err := json.Marshal(contratto)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(ID, contractJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*IDRelease, error) {
	idreleaseJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if idreleaseJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var idrelease IDRelease
	err = json.Unmarshal(idreleaseJSON, &idrelease)
	if err != nil {
		return nil, err
	}

	return &idrelease, nil
}

// TrunkExists returns true when asset with given ID exists in world state
func ContractExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	contractJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return contractJSON != nil, nil
}

func (s *SmartContract) GetAllFreeIDsForOrg(ctx contractapi.TransactionContextInterface, org string) ([]*IDRelease, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("NOT WORKING")
	}
	defer resultsIterator.Close()

	var ids []*IDRelease
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var id IDRelease
		err = json.Unmarshal(queryResponse.Value, &id)
		if err != nil {
			return nil, err
		}
		if id.Owner == org && id.Free == "true" {
			ids = append(ids, &id)
		}

	}
	return ids, nil
}

func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, ID string, Issuer string, Owner string, Tipo string, Free string) error {
	exists, err := ContractExists(ctx, ID)

	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("the asset %s does not exist", ID)
	}

	contratto := IDRelease{
		ID:     ID,
		Tipo:   Tipo,
		Owner:  Owner,
		Issuer: Issuer,
		Free:   Free,
	}

	contractJSON, err := json.Marshal(contratto)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(ID, contractJSON)
}

func (s *SmartContract) SetAssetAsUsed(ctx contractapi.TransactionContextInterface, ID string) error {
	idreleaseJSON, err := ctx.GetStub().GetState(ID)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if idreleaseJSON == nil {
		return fmt.Errorf("the asset %s does not exist", ID)
	}

	var idrelease IDRelease
	err = json.Unmarshal(idreleaseJSON, &idrelease)
	if err != nil {
		return err
	}
	idrelease.Free = "false"

	contractJSON, _ := json.Marshal(idrelease)

	return ctx.GetStub().PutState(ID, contractJSON)
}
