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

const DefaultTagCreator = "UATag.system"

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
	if function == "init" {
		//return t.Init(stub, "init", args) // init here resets all data !
		return t.read(stub, []string{"system_created_time"})
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

	err stub.PutState("create_tag", []byte(args[0]))

	var err error
	if len(args) < 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting Atleast 2")
	}

	fmt.Println("- start init tag")

	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}

	_tag_Id := strings.ToUpper(args[0])
	_tag_CreatedAt := strings.ToLower(args[1])
	_tag_ChaincodedAt := time.Now().String()
	_tag_Creator, _tag_IssuedTo, _tag_IssuedAt := "", "", ""

	if args[3] != "" {
		_tag_Creator = args[3]
	} else {
		_tag_Creator = DefaultTagCreator
	}

	if len(args[4]) > 0 && len(args[5]) > 0 {
		_tag_IssuedTo = args[4]
		_tag_IssuedAt = args[5]
	}

	str := `{"Id": "` + _tag_Id + `", "CreatedAt": "` + _tag_CreatedAt + `", "ChaincodedAt": "` + _tag_ChaincodedAt + `", "Creator": "` + _tag_Creator + `", "IssuedTo": "` + _tag_IssuedTo + `", "IssuedAt": "` + _tag_IssuedAt + `"}`
	err = stub.PutState(_tag_Id, []byte(str)) //store marble with id as key
	if err != nil {
		return nil, err
	}

	//get the tag index
	tagsAsBytes, err := stub.GetState(tagIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get tag index")
	}
	var tagIndex []string
	json.Unmarshal(tagsAsBytes, &tagIndexStr) //un stringify it aka JSON.parse()

	//append
	tagIndex = append(tagIndex, _tag_Id) //add marble name to index list
	fmt.Println("! tag index: ", tagIndex)
	jsonAsBytes, _ := json.Marshal(tagIndex)
	err = stub.PutState(tagIndexStr, jsonAsBytes) //store name of marble

	fmt.Println("- end init tag")
	return nil, nil
}
