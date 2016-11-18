package main

import (
	"fmt"
	"errors"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
)

type SimpleChaincode struct {
}

type Owner struct {
	OwnerID string `json:"ownerID"`
	CowsOwned []Cow `json:"cows"`
	//TODO(isaac) should never have more policies than cows
}

type Cow struct {
	OwnerID string `json:"ownerID"`
	SensorID string `json:"sensorID"`
}

type Policy struct {
	SensorID string `json:"sensorID"`
	Premium int `json:"premium"`
	Value int `json:"value"`
}

type AllCows struct {
	Catalog []Cow `json:"cows"`
}

type AllOwners struct {
	Catalog []Owner `json:"owners"`
}

type AllPolicies struct {
	Catalog []Policy `json:"policies"`
}

var activePoliciesString = "_activePolicies"
var activeCowsString = "_activeCows"
var activeOwnersString = "_activeOwners"

//==============================================================================
//==============================================================================

func main() {
	fmt.Println("Function: main")

	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting simple chaincode: %s", err)
	}
}

//==============================================================================
//==============================================================================

func makeHash(args []string) string {
	if len(args) < 0 {
		return "cannot generate hash"
	}

	i := 0
	s := ""
	for i < len(args){
		s = s + args[i]
		i = i + 1
	}
	return s
}

//==============================================================================
//==============================================================================

func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("Method: SimpleChaincode.Init")

	// Initialize the catalogs for both pending and active policies
	ownerCatalog := make([]Owner, 0)
	cowCatalog := make([]Cow, 0)
	policyCatalog := make([]Policy, 0)

	//Create and marshal the active policies
	var activePolicies AllPolicies
	activePolicies.Catalog = policyCatalog
	var policiesAsBytes []byte
	policiesAsBytes, err := json.Marshal(activePolicies)
	if err != nil {
		return nil, err
	}

	// Create and marshal the pending policies
	var activeCows AllCows
	activeCows.Catalog = cowCatalog
	var cowsAsBytes []byte
	cowsAsBytes, err = json.Marshal(activeCows)
	if err != nil {
		return nil, err
	}

	// Create and marshal incomplete policies
	var activeOwners AllOwners
	activeOwners.Catalog = ownerCatalog
	var ownersAsBytes []byte
	ownersAsBytes, err = json.Marshal(activeOwners)
	if err != nil {
		return nil, err
	}

	err = stub.PutState(activePoliciesString, policiesAsBytes)
	if err != nil {
		fmt.Println("Failed to initialize policies")
		return nil, err
	}

	err = stub.PutState(activeCowsString, cowsAsBytes)
	if err != nil {
		fmt.Println("Failed to initialize active cows")
		return nil, err
	}

	err = stub.PutState(activeOwnersString, ownersAsBytes)
	if err != nil {
		fmt.Println("Failed to initialize active owners")
		return nil, err}

	fmt.Println("Initializatin Complete")
	return nil, nil
}

//==============================================================================
//==============================================================================

func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("Method: SimpleChaincode.Invoke; received: " + function)

	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "generatePolicy" {
		return generatePolicy(stub, args)
	} else if function == "sensorTriggered" {
		return nil, sensorTriggered(stub, args)
	}

	fmt.Println("Invoke did not find a function: " + function)
	return nil, errors.New("Received unknown function invocation")
}

//==============================================================================
//==============================================================================

func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("Method: SimpleChaincode.Query; received: " + function)

	if function == "getActivePolicies" {
		return getAll(stub, activePoliciesString)
	} else if function == "getActiveCows" {
		return getAll(stub, activeCowsString)
	} else if function == "getActiveOwners" {
		return getAll(stub, activeOwnersString)
	}
	// else if function == "getPolicyOwner" {
	// 	return getPolicyOwner(stub, args)
	// }

	fmt.Println("Query did not find a function: " + function)
	return nil, errors.New("Received unknown function query")
}

//==============================================================================
//==============================================================================

func write(stub *shim.ChaincodeStub, name string, value []byte) error {
	fmt.Println("Function: write")

	err := stub.PutState(name, value)
	if err != nil {
		return err
	}
	return nil
}

//==============================================================================
//==============================================================================

// func getOwners(stub *shim.ChaincodeStub, ownersString string) ([]byte, error) {
// 	fmt.Println("Function: getOwners (" + ownersString + ")")
//
// 	ownersAsBytes, err := stub.GetState(ownersString)
// 	if err != nil {
// 		jsonResp := "{\"Error\": \"Failed to get owners.\"}"
// 		return nil, errors.New(jsonResp)
// 	}
//
// 	return ownersAsBytes, nil
// }
//
// //==============================================================================
// //==============================================================================
//
// func getCows(stub *shim.ChaincodeStub) ([]byte, error) {
// 	fmt.Println("Function: getCows (" + activeCowsString + ")")
//
// 	cowsAsBytes, err := stub.GetState(activeCowsString)
// 	if err != nil {
// 		jsonResp := "{\"Error\": \"Failed to get cows.\"}"
// 		return nil, errors.New(jsonResp)
// 	}
//
// 	return cowsAsBytes, nil
// }
//
// //==============================================================================
// //==============================================================================
//
// func getPolicies(stub *shim.ChaincodeStub) ([]byte, error) {
// 	fmt.Println("Function: getPolicies (" + policiesString + ")")
//
// 	policiesAsBytes, err := stub.GetState(policiesString)
// 	if err != nil {
// 		jsonResp := "{\"Error\": \"Failed to get policies.\"}"
// 		return nil, errors.New(jsonResp)
// 	}
//
// 	return policiesAsBytes, nil
// }

//==============================================================================
//==============================================================================

func getAll(stub *shim.ChaincodeStub, objectString string) ([]byte, error) {
	fmt.Println("Function: getAll (" + objectString + ")")

	objectsAsBytes, err := stub.GetState(objectString)
	if err != nil {
		jsonResp := "{\"Error\": \"Failed to get objects.\"}"
		return nil, errors.New(jsonResp)
	}

	return objectsAsBytes, nil
}

//==============================================================================
//==============================================================================

// func getPolicyOwner(stub *shim.ChaincodeStub, args []string) (int, error) {
// 	fmt.Println("Function: getPolicyOwner")
//
// 	policyID := args[0]
//
// 	cows := getCows(stub)
//
// 	var i int
// 	i = 0
// 	for i < len(cows) {
// 		if cows[i].SensorID == policyID {
// 			return cows[i].OwnerID, nil
// 		}
// 		i = i + 1
// 	}
//
// 	return nil, errors.New("No policy found with this policy ID: " + policyID)
// }

//==============================================================================
//==============================================================================

func writePolicies(stub *shim.ChaincodeStub, policies AllPolicies) error {
	fmt.Println("Function: writePolicies")

	policiesAsBytes, err := json.Marshal(policies)
	if err != nil {
		return err
	}
	fmt.Println("policies have been converted to bytes")

	err = write(stub, activePoliciesString, policiesAsBytes)
	if err != nil {
		return err
	}
	fmt.Println("policies written")
	return nil
}

//==============================================================================
//==============================================================================

func writeOwners(stub *shim.ChaincodeStub, owners AllOwners) error {
	fmt.Println("Function: writeOwners")

	ownersAsBytes, err := json.Marshal(owners)
	if err != nil {
		return err
	}
	fmt.Println("owners have been converted to bytes")

	err = write(stub, activeOwnersString, ownersAsBytes)
	if err != nil {
		return err
	}
	fmt.Println("owners written")
	return nil
}

//==============================================================================
//==============================================================================

func writeCows(stub *shim.ChaincodeStub, cows AllCows) error {
	fmt.Println("Function: writeCows")

	cowsAsBytes, err := json.Marshal(cows)
	if err != nil {
		return err
	}
	fmt.Println("cows have been converted to bytes")

	err = write(stub, activeCowsString, cowsAsBytes)
	if err != nil {
		return err
	}
	fmt.Println("cows written")
	return nil
}

//==============================================================================
//==============================================================================

func (t *SimpleChaincode) registerOwner(stub *shim.ChaincodeStub, args []string) ([]byte, error){
	//TODO(isaac) assign the owner a unique ID #
	//TODO(isaac) make an owner object with this ID #
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2 (First & Last name)")
	}

	var owner Owner
	owner.OwnerID = makeHash(args)
	return nil, nil
}

//==============================================================================
//==============================================================================

func addOwner(stub *shim.ChaincodeStub, owner Owner) error {
	fmt.Println("Function: addOwner")

	ownersAsBytes, err := getAll(stub, activeOwnersString)
	if err != nil {
		return err
	}
	fmt.Println("all owners retrieved")

	var listOfOwners AllOwners
	listOfOwners, err = bytesToAllOwners(ownersAsBytes)
	if err != nil {
		return err
	}
	fmt.Println("all owners derived from bytes")

	listOfOwners.Catalog = append(listOfOwners.Catalog, owner)
	fmt.Println("owner appended to list of owners")

	ownersAsBytes, err = json.Marshal(listOfOwners)
	err = write(stub, activeOwnersString, ownersAsBytes)
	if err != nil {
		return err
	}
	fmt.Println("list of owners successfully rewritten with new owner")
	return nil
}

//==============================================================================
//==============================================================================

func addCow(stub *shim.ChaincodeStub, cow Cow) error {
	fmt.Println("Function: addCow")

	cowsAsBytes, err := getAll(stub, activeCowsString)
	if err != nil {
		return err
	}
	fmt.Println("all cows retrieved")

	var listOfCows AllCows
	listOfCows, err = bytesToAllCows(cowsAsBytes)
	if err != nil {
		return err
	}
	fmt.Println("all cows derived from bytes")

	listOfCows.Catalog = append(listOfCows.Catalog, cow)
	fmt.Println("cow appended to list of cows")

	cowsAsBytes, err = json.Marshal(listOfCows)
	err = write(stub, activeCowsString, cowsAsBytes)
	if err != nil {
		return err
	}
	fmt.Println("list of cows successfully rewritten with new cow")
	return nil
}

//==============================================================================
//==============================================================================

func bytesToAllPolicies(policiesAsBytes []byte) (AllPolicies, error) {
	fmt.Println("Function: bytesToAllPolicies")

	var policies AllPolicies

	err := json.Unmarshal(policiesAsBytes, &policies)
	fmt.Println("json.Unmarshal error:")
	fmt.Println(err)

	return policies, err
}

//==============================================================================
//==============================================================================

func bytesToAllOwners(ownersAsBytes []byte) (AllOwners, error) {
	fmt.Println("Function: bytesToAllOwners")

	var owners AllOwners

	err := json.Unmarshal(ownersAsBytes, &owners)
	fmt.Println("json.Unmarshal error:")
	fmt.Println(err)

	return owners, err
}

//==============================================================================
//==============================================================================

func bytesToAllCows(cowsAsBytes []byte) (AllCows, error) {
	fmt.Println("Function: bytesToAllCows")

	var cows AllCows

	err := json.Unmarshal(cowsAsBytes, &cows)
	fmt.Println("json.Unmarshal error:")
	fmt.Println(err)

	return cows, err
}

//==============================================================================
//==============================================================================

func createPolicyObject(ID string) Policy {
	fmt.Println("Function: createPolicyObject")

	//TODO(isaac) I need to enter the cow object, the premium, and the value to the policy object

	var policy Policy

	policy.SensorID = ID
	policy.Premium = 100
	policy.Value = 5000

	return policy

}

//==============================================================================
//==============================================================================

func generatePolicy(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	//TODO(isaac) exoecting one argument (cowID), might take ownerID
	if len(args) < 1 {
		return nil, errors.New("Expected multiple arguments; arguments received: " +  strconv.Itoa(len(args)))
	}

	// check to make sure that cow doesn't already have a policy associated to it
	//TODO(isaac) check to see if the cow has been registered

	newPolicy := createPolicyObject(args[0])

	//TODO(isaac) add the new policy to the list of current policies
	// add the policy to the owner
	policiesAsBytes, err := getAll(stub, activePoliciesString)
	if err != nil {
		return nil, err
	}

	var policies AllPolicies
	policies, err = bytesToAllPolicies(policiesAsBytes)
	if err != nil {
		return nil, err
	}

	// Add the new policy to the list of pending policies
	policies.Catalog = append(policies.Catalog, newPolicy)
	fmt.Println("New policy appended to policies. Incomplete policy count: " + strconv.Itoa(len(policies.Catalog)))

	err = writePolicies(stub, policies)
	if err != nil {
		return nil, err
	}
	fmt.Println("policies successfully rewritten with new policy")
	return nil, nil

}

//==============================================================================
//==============================================================================

func getCowIndexBySensor(cows []Cow, sensorID string) (int, error) {
	fmt.Println("Function: getCowBySensor")

	var i int
	i = 0
	for i < len(cows) {
		if cows[i].SensorID == sensorID {
			return i, nil
		}
		i = i + 1
	}

	return 0, errors.New("No cow found with this sensor ID: " + sensorID)
}

//==============================================================================
//==============================================================================

func getPolicyIndexByID(policies []Policy, sensorID string) (int, error) {
	fmt.Println("Function: getCowBySensor")

	var i int
	i = 0
	for i < len(policies) {
		if policies[i].SensorID == sensorID {
			return i, nil
		}
		i = i + 1
	}

	return 0, errors.New("No cow found with this sensor ID: " + sensorID)
}

//==============================================================================
//==============================================================================

func sensorTriggered(stub *shim.ChaincodeStub, args []string) error {
	fmt.Println("Function: sensorTriggered")

	sensorID := args[0]

	return cowDeath(stub, sensorID)
}

//==============================================================================
//==============================================================================

func cowDeath(stub *shim.ChaincodeStub, sensorID string) error {
	//TODO needs to call pay out

	var cows AllCows
	cowsAsBytes, err := getAll(stub, activeCowsString)
	cows, err = bytesToAllCows(cowsAsBytes)

	fmt.Println("Function: cowDeath")
	index, err := getCowIndexBySensor(cows.Catalog, sensorID)
	if err != nil {
		return err
	}

	copy(cows.Catalog[:index], cows.Catalog[index + 1:])
	cows.Catalog = cows.Catalog[:len(cows.Catalog) - 1]

	var policies AllPolicies
	policiesAsBytes, err := getAll(stub, activePoliciesString)
	policies, err = bytesToAllPolicies(policiesAsBytes)
	if err != nil {
		return err
	}

	index, err = getPolicyIndexByID(policies.Catalog, sensorID)
	if err != nil {
		return err
	}

	copy(policies.Catalog[:index], policies.Catalog[index + 1:])
	policies.Catalog = policies.Catalog[:len(policies.Catalog) - 1]

	payOut, err := payOut(policies.Catalog, sensorID)
	if err != nil {
		return err
	}

	fmt.Println(payOut)

	writeCows(stub, cows)
	writePolicies(stub, policies)

	return nil
}

//==============================================================================
//==============================================================================

func payOut(policies []Policy, sensorID string) (int, error) {

	var i int
	i = 0
	for i < len(policies) {
		if policies[i].SensorID == sensorID {
			return policies[i].Value, nil
		}
		i = i + 1
	}

	return 0, errors.New("No policy found with this policy ID: " + sensorID)
}
