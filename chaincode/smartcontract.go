package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

type Artist struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DateOfBirth string `json:"dateOfBirth"`
}

// Asset describes basic details of Artist certificate
type Certificate struct {
	ID               string `json:"ID"`
	PHOTO            string `json:"URI"`
	Title            string `json:"title"`
	Owner            string `json:"owner"`
	YearOfProduction int    `json:"appraisedValue"`
	Artist           Artist `json:"artist"`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Certificate{
		{PHOTO: "https://images.unsplash.com/photo-1571115764595-644a1f56a55c?ixid=MnwxMjA3fDB8MHxzZWFyY2h8NHx8cGFpbnRpbmd8ZW58MHx8MHx8&ixlib=rb-1.2.1&w=1000&q=80",
			ID:    "8989s1gjJJHJKHJSGHJDJSAD871238S",
			Title: "original-splash-abstract",
			Owner: "Tomoko Janra",
			Artist: Artist{
				"id-13894047849",
				"Tomoko Janra",
				"20.11.1980"},
			YearOfProduction: 1962},
		{PHOTO: "https://ii1.pepperfry.com/media/catalog/product/o/r/568x625/original-handmade-couple-together-oil-painting-on-canvas-by-gallery99-original-handmade-couple-toget-g9u0rd.jpg",
			ID:    "HDFSJ52151511113GHJDJSAD8712389",
			Title: "original-handmade-couple-toget",
			Owner: "Eyas Jaber",
			Artist: Artist{
				"id-13894047849",
				"Mayer Labor",
				"17.11.1971"},
			YearOfProduction: 1977},
		{PHOTO: "https://www.terrain.org/articles/27/fullsize/15_fs.jpg",
			ID:    "34890DFASHJK148904KLOJKUKOP41IO4",
			Title: "original-piano-farm.",
			Owner: "Sebastian Michel",
			Artist: Artist{
				"id-13899810498",
				"Jubran Jubran",
				"17.05.1965"},
			YearOfProduction: 1990},
		{PHOTO: "https://qph.fs.quoracdn.net/main-qimg-a7cbe92701920a2155718c76c58c7c52",
			ID:    "JFKLHJKGHFAK87419FADSHJLgdf4hjaj1",
			Title: "british-guards",
			Owner: "Monde Cardino",
			Artist: Artist{
				"id-13899810498",
				"Jubran Jubran",
				"17.07.1960"},
			YearOfProduction: 2017},
	}

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		keyState := asset.ID + asset.Owner
		err = ctx.GetStub().PutState(keyState, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, photo string, title string, artist Artist, owner string, yearOfProduction int) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	// owner, _ := ctx.GetClientIdentity().GetID()
	asset := Certificate{
		ID:               id,
		PHOTO:            photo,
		Title:            title,
		Artist:           artist,
		Owner:            owner,
		YearOfProduction: yearOfProduction,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	// keyState := id + asset.Owner
	fmt.Println(id)
	// CreateCompositeKeyStub
	return ctx.GetStub().PutState(id, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Certificate, error) {
	// owner := "owner"
	// fmt.Println((id))
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset Certificate
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, photo string, title string, artist Artist, owner string, yearOfProduction int) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	// overwriting original asset with new asset
	// owner, _ := ctx.GetClientIdentity().GetID()
	asset := Certificate{
		ID:               id,
		PHOTO:            photo,
		Title:            title,
		Artist:           artist,
		Owner:            owner,
		YearOfProduction: yearOfProduction,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	keyState := id + asset.Owner
	return ctx.GetStub().PutState(keyState, assetJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// TransferAsset updates the owner field of asset with given id in world state.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return err
	}

	newAsset := asset
	newAsset.Owner = newOwner
	fmt.Println(newAsset, newOwner)
	assetJSON, err := json.Marshal("string")
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Certificate, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Certificate
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Certificate
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}
