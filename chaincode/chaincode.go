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
	"time"
)

// SimpleChaincode -> simple Chaincode implementation
type SimpleChaincode struct {
}

// declaration of Tag object
type Tag struct {
	Id        string `json:"id"`         //the fieldtags are needed to keep case from bouncing around
	CreatedAt string `json:"created_at"` // creation date of tag -> when it was placed to chaincode
	Creator   string `json:"creator"`    // creator -> who created? Obiously, Uatag
	IssuedTo  string `json:"issued_to"`  // Company name issued to
	IssuedAt  string `json:"issued_at"`  // the date when tag was issued to company
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

// INIT THE WHOLE NEW NETWORK CHAINCODE and test if writeable
func (t *SimpleChaincode) init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	var Aval int
	var err error

	// Initialize the chaincode
	Aval, err = strconv.Atoi("1")
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}

	// Write the state to the ledger
	err = stub.PutState("abc", []byte(strconv.Itoa(Aval))) //making a test var "abc", I find it handy to read/write to it right away to test the network
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
	return nil, nil
}

//// -------------------- OK

// ============================================================================================================================
// Run - Our entry point for Invokcations
// ============================================================================================================================
func (t *SimpleChaincode) Run(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("run is running " + function)

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset
		return t.init(stub, args)
	}
	fmt.Println("run did not find func: " + function) //error

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
		return t.read(stub, "system_created_time")
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
