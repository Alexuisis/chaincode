/*
Copyright Oleksandr Samoilov ossamoylov@gmail.com 2016 All Rights Reserved.
This is proprietary software. You have no rights to use it unless the permission was given.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
	"strings"
	"time"
)

// SimpleChaincode -> simple Chaincode implementation
type SimpleChaincode struct {
}

// declaration of Tag object
type Tag struct {
	Id           string `json:"id"`            //the fieldtags are needed to keep case from bouncing around
	CreatedAt    string `json:"created_at"`    // creation date of tag -> when it was physically created
	ChaincodedAt string `json:"chaincoded_at"` // creation date of tag -> when it was placed to chaincode
	Creator      string `json:"creator"`       // creator -> who created? Obiously, Uatag
	IssuedTo     string `json:"issued_to"`     // Company name issued to
	IssuedAt     string `json:"issued_at"`     // the date when tag was issued to company
}

// ALL TAGS INDEXES
var tagIndexStr = "_tagindex" //name for the key/value that will store a list of all known tags

//MAIN ENTRY POINT
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// ============================================================================================================================
// Init - creates the database and tests it for usage
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	var Aval int
	var err error

	// Initialize the chaincode
	Aval, err = strconv.Atoi("1")
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}

	// Write the state to the ledger
	err = stub.PutState("already_inited", []byte(strconv.Itoa(Aval)))
	if err != nil {
		return nil, err
	}

	// CLEAR AND CREATE NEW TAGS DATABASE
	var empty []string
	jsonAsBytes, _ := json.Marshal(empty) //marshal an emtpy array of strings to clear the index
	err = stub.PutState(tagIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	// WE LOG CREATION DATE & TIME UPON INIT
	_init_time := time.Now().String()
	err = stub.PutState("system_created_time", []byte(_init_time))
	if err != nil {
		return nil, err
	}

	err = stub.PutState("system_author", []byte("Oleksandr Alex Samoilov"))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ============================================================================================================================
// Invoke - ENTRY POINT FOR ALLLLL INVOKATIONS
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset

		valAsbytes, err := stub.GetState("already_inited") //get the var from chaincode state
		if err != nil {
			return valAsbytes, nil
		} else {
			return t.Init(stub, function, args)
		}

	} else if function == "create_tag" {
		return t.create_tag(stub, args)
	}

	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation")
}

// ============================================================================================================================
// Query - Our entry point for Queries
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	} else if function == "system_created" {
		return t.read(stub, []string{"system_created_time"})
	} else if function == "system_author" {
		return t.read(stub, []string{"system_author"})
	}

	fmt.Println("query did not find func: " + function) //error

	return nil, errors.New("Received unknown function query")
}

// ============================================================================================================================
// Read - read a variable from chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetState(name) //get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil //send it onward
}

// ============================================================================================================================
// 	ALL TAGS RELATES FUNCTIONS BELOW
// ============================================================================================================================

/*
	REQ ARGS
	1 -> Tag ID
	2 -> Created At
	NOT REQ ARGS
	3 -> Creator
	4 -> IssuedTo
	5 -> IssuedAt

*/
func (t *SimpleChaincode) create_tag(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	err := stub.PutState("c1", []byte("1"))
	if err != nil {
		return nil, err
	}

	if len(args) < 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting  Atleast 2")
	}

	fmt.Println("- start tag marble")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}

	err = stub.PutState("c2", []byte("2"))
	if err != nil {
		return nil, err
	}

	tag_Id := strings.ToLower(args[0])
	tag_CreatedAt := strings.ToLower(args[1])
	tag_ChaincodedAt := time.Now().String()
	tag_Creator := ""
	tag_IssuedTo := ""
	tag_IssuedAt := ""

	err = stub.PutState("c3", []byte("3"))
	if err != nil {
		return nil, err
	}

	str := `{"Id": "` + tag_Id + `", "CreatedAt": "` + tag_CreatedAt + `", "ChaincodedAt": ` + tag_ChaincodedAt + `, "Creator": "` + tag_Creator + `", "IssuedTo": "` + tag_IssuedTo + `", "IssuedAt": "` + tag_IssuedAt + `"}`

	err = stub.PutState("c4", []byte("4"))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(tag_Id, []byte(str)) //store marble with id as key
	if err != nil {
		return nil, err
	}

	err = stub.PutState("c5", []byte("5"))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
