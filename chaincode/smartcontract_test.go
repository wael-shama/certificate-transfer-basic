package chaincode_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-samples/asset-transfer-basic/chaincode-go/chaincode"
	"github.com/hyperledger/fabric-samples/asset-transfer-basic/chaincode-go/chaincode/mocks"
	"github.com/stretchr/testify/require"
)

//go:generate counterfeiter -o mocks/transaction.go -fake-name TransactionContext . transactionContext
type transactionContext interface {
	contractapi.TransactionContextInterface
}

//go:generate counterfeiter -o mocks/chaincodestub.go -fake-name ChaincodeStub . chaincodeStub
type chaincodeStub interface {
	shim.ChaincodeStubInterface
}

//go:generate counterfeiter -o mocks/statequeryiterator.go -fake-name StateQueryIterator . stateQueryIterator
type stateQueryIterator interface {
	shim.StateQueryIteratorInterface
}

func TestInitLedger(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	assetTransfer := chaincode.SmartContract{}
	err := assetTransfer.InitLedger(transactionContext)
	require.NoError(t, err)

	chaincodeStub.PutStateReturns(fmt.Errorf("failed inserting key"))
	err = assetTransfer.InitLedger(transactionContext)
	require.EqualError(t, err, "failed to put to world state. failed inserting key")
}

func TestCreateAsset(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}

	fmt.Println(transactionContext)
	transactionContext.GetStubReturns(chaincodeStub)

	assetTransfer := chaincode.SmartContract{}

	// Creating new Certificate for photo Passing args :
	// id string, photo string, title string, artist string, yearOfProduction int
	artist := chaincode.Artist{
		"id-13894047849",
		"Tomoko Janra",
		"20.11.1980"}
	err := assetTransfer.CreateAsset(transactionContext, "original-cowBoy", "photo_uri", "title", artist, "owner1", 1962)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns([]byte{}, nil)

	err = assetTransfer.CreateAsset(transactionContext, "original-cowBoy", "photo_uri", "title", artist, "owner1", 1962)
	require.EqualError(t, err, "the asset original-cowBoy already exists")

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	err = assetTransfer.CreateAsset(transactionContext, "id", "photo_uri", "title", artist, "owner1", 1962)
	require.EqualError(t, err, "failed to read from world state: unable to retrieve asset")
}

func TestReadAsset(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}

	transactionContext.GetStubReturns(chaincodeStub)

	artist := chaincode.Artist{
		"id-13894047849",
		"Tomoko Janra",
		"20.11.1980"}

	expectedAsset := &chaincode.Certificate{PHOTO: "https://images.unsplash.com/photo-1571115764595-644a1f56a55c?ixid=MnwxMjA3fDB8MHxzZWFyY2h8NHx8cGFpbnRpbmd8ZW58MHx8MHx8&ixlib=rb-1.2.1&w=1000&q=80",
		ID:               "8989s1gjJJHJKHJSGHJDJSAD871238S",
		Title:            "original-splash-abstract",
		Owner:            "owner",
		Artist:           artist,
		YearOfProduction: 1962}
	bytes, err := json.Marshal(expectedAsset)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)

	assetTransfer := chaincode.SmartContract{}
	_ = assetTransfer.CreateAsset(transactionContext, "8989s1gjJJHJKHJSGHJDJSAD871238S", "https://images.unsplash.com/photo-1571115764595-644a1f56a55c?ixid=MnwxMjA3fDB8MHxzZWFyY2h8NHx8cGFpbnRpbmd8ZW58MHx8MHx8&ixlib=rb-1.2.1&w=1000&q=80", "original-splash-abstract", artist, "Tomoko Janra", 1962)

	asset, err := assetTransfer.ReadAsset(transactionContext, "8989s1gjJJHJKHJSGHJDJSAD871238S")
	require.NoError(t, err)
	require.Equal(t, expectedAsset, asset)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	_, err = assetTransfer.ReadAsset(transactionContext, "")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve asset")

	chaincodeStub.GetStateReturns(nil, nil)
	asset, err = assetTransfer.ReadAsset(transactionContext, "asset1")
	require.EqualError(t, err, "the asset asset1 does not exist")
	require.Nil(t, asset)
}

func TestTransferAsset(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	assetTransfer := chaincode.SmartContract{}
	artist := chaincode.Artist{
		"id-13894047849",
		"Tomoko Janra",
		"20.11.1980"}

	expectedAsset := &chaincode.Certificate{PHOTO: "https://images.unsplash.com/photo-1571115764595-644a1f56a55c?ixid=MnwxMjA3fDB8MHxzZWFyY2h8NHx8cGFpbnRpbmd8ZW58MHx8MHx8&ixlib=rb-1.2.1&w=1000&q=80",
		ID:               "8989s1gjJJHJKHJSGHJDJSAD871238S",
		Title:            "original-splash-abstract",
		Owner:            "owner",
		Artist:           artist,
		YearOfProduction: 1962}
	bytes, err := json.Marshal(expectedAsset)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)

	assetTransfer = chaincode.SmartContract{}
	_ = assetTransfer.CreateAsset(transactionContext, "8989s1gjJJHJKHJSGHJDJSAD871238S", "https://images.unsplash.com/photo-1571115764595-644a1f56a55c?ixid=MnwxMjA3fDB8MHxzZWFyY2h8NHx8cGFpbnRpbmd8ZW58MHx8MHx8&ixlib=rb-1.2.1&w=1000&q=80", "original-splash-abstract", artist, "Tomoko Janra", 1962)

	asset, err := assetTransfer.ReadAsset(transactionContext, "8989s1gjJJHJKHJSGHJDJSAD871238S")
	fmt.Println(asset)
	require.NoError(t, err)
	require.Equal(t, expectedAsset, asset)

	chaincodeStub.GetStateReturns(bytes, nil)
	err = assetTransfer.TransferAsset(transactionContext, "8989s1gjJJHJKHJSGHJDJSAD871238S", "owner2")
	require.NoError(t, err)

	_, err = assetTransfer.ReadAsset(transactionContext, "8989s1gjJJHJKHJSGHJDJSAD871238S")
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	err = assetTransfer.TransferAsset(transactionContext, "", "")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve asset")
}
