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
	Id           string `json:"Id"`           //the fieldtags are needed to keep case from bouncing around
	CreatedAt    string `json:"CreatedAt"`    // creation date of tag -> when it was physically created
	ChaincodedAt string `json:"ChaincodedAt"` // creation date of tag -> when it was placed to chaincode
	Creator      string `json:"Creator"`      // creator -> who created? Obiously, Uatag
	IssuedTo     string `json:"IssuedTo"`     // Company name issued to
	IssuedAt     string `json:"IssuedAt"`     // the date when tag was issued to company
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

/// ghcfghgfy

// ============================================================================================================================
// Invoke - ENTRY POINT FOR ALLLLL INVOKATIONS -> THIS IS A ROUTER FOR ALL  INVOKATIONAL FUNCTIONS
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset

		valAsbytes, err := stub.GetState("already_inited")
		if err != nil {
			return t.Init(stub, function, args)
		} else {
			return valAsbytes, nil
		}

	} else if function == "create_tag" {
		return t.create_tag(stub, args)
	} else if function == "assign_to" {
		return t.assign_to(stub, args)
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
//
//
// 			ALL 	TAGS 	RELATES 	FUNCTIONS 	BELOW
//
//
// ============================================================================================================================

// ============================================================================================================================
// 	0. Create Tag
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

	tag_Id := strings.ToUpper(args[0])
	tag_key := "tag_" + tag_Id
	tag_ChaincodedAt := time.Now().String()
	tag_Creator := "UATag.system"

	tag_CreatedAt, tag_IssuedTo, tag_IssuedAt := "", "", ""

	if len(args) == 2 {

		if len(args[1]) > 0 {
			tag_CreatedAt = args[1]
		}

	}

	if len(args) == 3 {

		if len(args[1]) > 0 {
			tag_CreatedAt = args[1]
		}

		if len(args[2]) > 0 {
			tag_Creator = args[2]
		}
	}

	if len(args) == 4 {

		if len(args[1]) > 0 {
			tag_CreatedAt = args[1]
		}

		if len(args[2]) > 0 {
			tag_Creator = args[2]
		}

		if len(args[3]) > 0 {
			tag_IssuedTo = args[3]
		}
	}

	if len(args) == 5 {

		if len(args[1]) > 0 {
			tag_CreatedAt = args[1]
		}

		if len(args[2]) > 0 {
			tag_Creator = args[2]
		}

		if len(args[3]) > 0 {
			tag_IssuedTo = args[3]
		}

		if len(args[4]) > 0 {
			tag_IssuedAt = args[4]
		}
	}

	str := `{"Id": "` + tag_Id + `", "CreatedAt": "` + tag_CreatedAt + `", "ChaincodedAt": "` + tag_ChaincodedAt + `", "Creator": "` + tag_Creator + `", "IssuedTo": "` + tag_IssuedTo + `", "IssuedAt": "` + tag_IssuedAt + `"}`

	err := stub.PutState(tag_key, []byte(str))
	if err != nil {
		return nil, err
	}

	//sore new tag in registry
	tagsAsBytes, errIndex := stub.GetState(tagIndexStr)
	if errIndex != nil {
		return nil, errors.New("Failed to get tags index")
	}
	var tagIndex []string
	json.Unmarshal(tagsAsBytes, &tagIndex)

	//append
	tagIndex = append(tagIndex, tag_key)
	fmt.Println("! tag index: ", tagIndex)
	jsonAsBytes, _ := json.Marshal(tagIndex)
	err = stub.PutState(tagIndexStr, jsonAsBytes)

	fmt.Println("- end init tag")

	return nil, nil
}

func (t *SimpleChaincode) assign_to(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	//   0       1            2
	// "tag_id", "assign_to", "assign_date"
	if len(args) < 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	tag_key := "tag_" + strings.ToUpper(args[0])

	tagAsBytes, err := stub.GetState(tag_key)
	if err != nil {
		return nil, errors.New("Failed to get tag")
	}

	_tag_Obj := Tag{}

	json.Unmarshal(tagAsBytes, &_tag_Obj) //un stringify it aka JSON.parse()

	_tag_Obj.IssuedTo = string(args[1]) //change the assigned to
	_tag_Obj.IssuedAt = string(args[2]) //change the assigned date

	jsonAsBytes, _ := json.Marshal(_tag_Obj)
	err = stub.PutState(tag_key, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
