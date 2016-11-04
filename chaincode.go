package main

import(
  "encoding/json"
	"errors"
	"fmt"
  "strconv"
	"time"
//	"string"
	"github.com/satori/go.uuid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type SimpleChaincode struct {}

type Cow struct {
	ID string
	OwnerID string
	SensorID string
}

type Owner struct {
	CowsOwned []Cow
	Policies []Policy
}

type Sensor struct {
	ID string
	CowID string
}

type Policy struct {
	ID string
	CowID []Cow
	OwnerID []Owner
	Premium int32
	Value int32
}

// =======================================================================================================================
// Make Timestamp - create a timestamp in ms
// =======================================================================================================================

func makeTimestamp() int64 {
    return time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
}

//========================================================================================================================
// Main
//========================================================================================================================

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil{
		fmt.Printf("Error starting simple chaincode: %s", err)
	}
}

//========================================================================================================================
// Initialize the state of the 'Policies' variable
//========================================================================================================================

func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting 1")
    }

    err := stub.PutState("hello_world", []byte(args[0]))
    if err != nil {
        return nil, err
    }

    return nil, nil
}

//========================================================================================================================
// Initialize the state of the 'Policies' variable
//========================================================================================================================

func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
    fmt.Println("invoke is running " + function)

    // Handle different functions

		//initialize the chaincode state, used as reset
    if function == "init" {
        return t.Init(stub, "init", args)
    }
		//deletes an entity from its state
		// else if function == "delete" {
		// 	return t.Delete(stub, args)
		// }
		//writes a value to the chaincode state
		else if function == "write" {
        return t.write(stub, args)
    }
		//create a new cow
		// else if function == "init_cow" {
		// 	return t.init_cow(stub, args)
		// }
		//change owner of a cow
		// else if function == "set_user" {
		// 	return t.set_user(stub, args)
		// }
    fmt.Println("invoke did not find func: " + function)

    return nil, errors.New("Received unknown function invocation")
}

//========================================================================================================================
// Registers a cow to the blockchain
//========================================================================================================================

func (t *SimpleChaincode) registerCow(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
		if len(args) < 2 {
			return nil, errors.New("Incorrect number of arguments. Expecting 2. ID of the owner and the ID of the sensor")
		}

		cowID := uuid.NewV4().String()
		ownerID := args[0]
		sensorID := args[1]

		var newCow Cow
		newCow.ID = cowID
		newCow.OwnerID = ownerID
		newCow.SensorID = sensorID

		if err != nil {
        return nil, err
    }
    return nil, nil
}

//========================================================================================================================
// Registers a policy to the blockchain
//========================================================================================================================

func (t *SimpleChaincode) registerPolicy(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) < 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4. ID of the cow, ID of the owner, the premium, and the value of the policy")
	}

	policyID := uuid.NewV4().String()
	cowID := args[0]
	ownerID := args[1]
	premium := args[2]
	value := args[3]

	var newPolicy Policy
	newPolicy.ID = policyID
	newPolicy.CowID = cowID
	newPolicy.OwnerID = ownerID
	newPolicy.Premium = premium
	newPolicy.Value = value

	if err != nil {
			return nil, err
	}
	return nil, nil

}

//========================================================================================================================
// Check the state of the chaincode
//========================================================================================================================

func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
    fmt.Println("query is running " + function)

    // Handle different functions
    if function == "read" {                            //read a variable
        return t.read(stub, args)
    }
    fmt.Println("query did not find func: " + function)

    return nil, errors.New("Received unknown function query")
}

//========================================================================================================================
// Write a value to a variable
//========================================================================================================================

func (t *SimpleChaincode) write(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    var name, value string
    var err error
    fmt.Println("running write()")

    if len(args) != 2 {
        return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the variable and value to set")
    }

    name = args[0]                            //rename for fun
    value = args[1]
    err = stub.PutState(name, []byte(value))  //write the variable into the chaincode state
    if err != nil {
        return nil, err
    }
    return nil, nil
}

//========================================================================================================================
// Read the state of a variable
//========================================================================================================================

func (t *SimpleChaincode) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    var name, jsonResp string
    var err error

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
    }

    name = args[0]
    valAsbytes, err := stub.GetState(name)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
        return nil, errors.New(jsonResp)
    }

    return valAsbytes, nil
}
