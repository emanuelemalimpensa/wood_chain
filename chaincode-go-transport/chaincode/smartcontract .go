package chaincode

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
// Insert struct field in alphabetic order => to achieve determinism across languages
// golang keeps the order when marshal to json but doesn't order automatically
type TransportAgreement struct {
	ID       string       `json:"ID"`
	Seller   *OrgApproval `json:Seller`
	Buyer    *OrgApproval `json:Buyer`
	Tipo     string       `json:Tipo` //incoterm associato
	Data     string       `json:Data`
	Quantita string       `json:Quantita`
	Status   string       `json:Status`
	Emitter  string       `json:Emitter`
}

type OrgApproval struct {
	Org      string `json:Org`
	Approval string `json:Approval`
}

type HistoryQueryResult struct {
	Record   *TransportAgreement `json:"record"`
	TxId     string              `json:"txId"`
	IsDelete bool                `json:"isDelete"`
}

func (s *SmartContract) PopulateAssetOrg(ctx contractapi.TransactionContextInterface, Org string, ID string, Role string) error {
	exists, err := TAExists(ctx, ID)

	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("failed to modify Transport Agreement %s does not exists. Error: %v", ID, err)

	} else {
		callOrg, _ := ctx.GetClientIdentity().GetMSPID()
		ta, err := s.ReadAsset(ctx, ID)
		if callOrg != ta.Emitter {
			return fmt.Errorf("ORGANIZATION NOT ALLOWED TO POPULATE")
		}
		if err != nil {
			return err
		}
		var org OrgApproval
		org.Org = Org
		org.Approval = "false"

		if Role == "Buyer" || Role == "buyer" || Role == "b" || Role == "B" {
			if ta.Buyer == nil {
				ta.Buyer = &org
			}
		}
		if Role == "Seller" || Role == "seller" || Role == "s" || Role == "S" {
			if ta.Seller == nil {
				ta.Seller = &org
			}
		}

		taJSON, err := json.Marshal(ta)
		if err != nil {
			return err
		}
		return ctx.GetStub().PutState(ID, taJSON)

	}
}

func (s *SmartContract) CreateBuffer(ctx contractapi.TransactionContextInterface, ID string, Tipo string, Data string, Quantita string) error {

	exists, err := TAExists(ctx, ID)

	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("failed to put into world state, Transport Agreement %s already exists. Error: %v", ID, err)

	} else {
		emitter, _ := ctx.GetClientIdentity().GetMSPID()
		ta := TransportAgreement{
			ID:       ID,
			Seller:   nil,
			Buyer:    nil,
			Tipo:     Tipo,
			Data:     Data,
			Quantita: Quantita,
			Status:   "Approving",
			Emitter:  emitter,
		}
		taJSON, err := json.Marshal(ta)
		if err != nil {
			return err
		}
		return ctx.GetStub().PutState(ID, taJSON)
	}
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, ID string) error {
	exists, err := TAExists(ctx, ID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", ID)
	}

	return ctx.GetStub().DelState(ID)
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*TransportAgreement, error) {

	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var ts []*TransportAgreement
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var ta TransportAgreement
		err = json.Unmarshal(queryResponse.Value, &ta)
		if err != nil {
			return nil, err
		}

		ts = append(ts, &ta)
	}

	return ts, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAvaiableTAsWithConstrains(ctx contractapi.TransactionContextInterface, myRole string, otherOrg string) ([]*TransportAgreement, error) {

	myOrg, _ := ctx.GetClientIdentity().GetMSPID()
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var ts []*TransportAgreement
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var ta TransportAgreement
		err = json.Unmarshal(queryResponse.Value, &ta)
		if err != nil {
			return nil, err
		}

		// ts = append(ts, &ta)

		if myRole == "Buyer" || myRole == "buyer" || myRole == "b" || myRole == "B" {
			if ta.Buyer.Org == myOrg {
				if ta.Seller.Org == otherOrg {
					ts = append(ts, &ta)
				}
			}
		}
		if myRole == "Seller" || myRole == "seller" || myRole == "s" || myRole == "S" {
			if ta.Seller.Org == myOrg {
				if ta.Buyer.Org == otherOrg {
					ts = append(ts, &ta)
				}
			}
		}
	}

	return ts, nil
}

func (t *SmartContract) GetAssetHistory(ctx contractapi.TransactionContextInterface, assetID string) ([]HistoryQueryResult, error) {

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
		var asset TransportAgreement
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &asset)
			if err != nil {
				return nil, err
			}
		} else {
			asset = TransportAgreement{
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

// func GetCreator(ctx contractapi.TransactionContextInterface) (string, error) {

// 	ctx.GetStub().GetTxID()
// 	l, _ := ctx.GetStub().GetCreator()
// 	to_ret := string(l)
// 	var i int
// 	var count int
// 	for i = 0; i < len(to_ret); i++ {
// 		if to_ret[i] == '\u0012' {
// 			break
// 		}
// 		count++
// 	}
// 	return to_ret[2:count], nil
// }

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	approvals := []TransportAgreement{
		{ID: "t1", Seller: nil, Buyer: nil, Tipo: "ECO", Data: "2020-01-01", Quantita: "10", Status: "Approving"},
		{ID: "t2", Seller: nil, Buyer: nil, Tipo: "ECO", Data: "2020-01-01", Quantita: "10", Status: "Approving"},
		{ID: "t3", Seller: nil, Buyer: nil, Tipo: "ECO", Data: "2020-01-01", Quantita: "10", Status: "Approving"},
	}

	for _, curr_approval := range approvals {
		approvalJSON, err := json.Marshal(curr_approval)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(curr_approval.ID, approvalJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*TransportAgreement, error) {
	taJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if taJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var ta TransportAgreement
	err = json.Unmarshal(taJSON, &ta)
	if err != nil {
		return nil, err
	}

	return &ta, nil
}

// TrunkExists returns true when asset with given ID exists in world state
func TAExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	taJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return taJSON != nil, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) ApproveTransferForMyOrg(ctx contractapi.TransactionContextInterface, ID string) error {
	exists, err := TAExists(ctx, ID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("TransferAgreement %s does not exist", ID)
	}

	ta, _ := s.ReadAsset(ctx, ID)

	org, _ := ctx.GetClientIdentity().GetMSPID()

	if ta.Buyer.Org == org {
		ta.Buyer.Approval = "true"
	} else if ta.Seller.Org == org {
		if ta.Buyer.Approval == "true" {
			//solo se il compratore ha approvato
			//allora e' possibile al venditore approvare
			ta.Seller.Approval = "true"
		} else {
			return fmt.Errorf("Current Org %s trying to approve but Buyer has not approved yet", org)
		}
	} else {
		return fmt.Errorf("Current Org cannot express objective Org will")
	}

	if ta.Seller.Approval == "true" {
		ta.Status = "Approved"
	}

	taJSON, err := json.Marshal(ta)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(ID, taJSON)
}
