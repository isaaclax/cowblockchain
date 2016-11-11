package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type Owner struct {
	OwnerID string 'json:"ownerID"'
	CowsOwned []Cows 'json:"cows"'
	//TODO(isaac) should never have more policies than cows
}

type Cow struct {
	OwnerID string 'json:"ownerID"'
	SensorID string 'json:"sensorID"'
}

type Policy struct {
	PolicyID string 'json:"policyID"'
	Cow []Cow 'json:"cow"'
	Premium int 'json:"premium"'
	Value int 'json:"value"'
}

type AllCows struct {
	Catalog []Cow 'json:"cows"'
}

type AllOwners struct {
	Catalog []Owner 'json:"owners"'
}

type AllPolicies struct {
	Catalog []Policy 'json:"policies"'
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
		return "no_hash_can_be_generated"
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
	policiesAsBytes, err := json.Marshal(activePolicies)
	if err != nil {
		return nil, err
	}

	// Create and marshal the pending policies
	var activeCows AllCows
	activeCows.Catalog = cowCatalog
	var pendingAsBytes []byte
	cowsAsBytes, err = json.Marshal(activeCows)
	if err != nil {
		return nil, err
	}

	// Create and marshal incomplete policies
	var activeOwners AllOwners
	activeOwners.Catalog = ownerCatalog
	var incompleteAsBytes []byte
	ownersAsBytes, err = json.Marshal(activeOwners)
	if err != nil {
		return nil, err
	}

	err = stub.PutState(activePoliciesString, activeAsBytes)
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
		return nil, err
}

//==============================================================================
//==============================================================================

func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("Method: SimpleChaincode.Invoke; received: " + function)

	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "generatePolicy" {
		return generatePolicy(stub, args)
	}

	fmt.Println("Invoke did not find a function: " + function)
	return nil, errors.New("Received unknown function invocation")
}

//==============================================================================
//==============================================================================

func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("Method: SimpleChaincode.Query; received: " + function)

	if function == "getActivePolicies" {
		return getPolicies(stub, activePoliciesString)
	} else if function == "getActiveCows" {
		return getCows(stub, activeCowsString)
	} else if function == "getActiveOwners" {
		return getOwners(stub, activeOwnersString)
	}

	fmt.Println("Query did not find a function: " + function)
	return nil, errors.New("Received unknown function query")
}

//==============================================================================
//==============================================================================

func getOwners(stub *shim.ChaincodeStub, ownersString string) ([]byte, error) {
	fmt.Println("Function: getOwners (" + ownersString + ")")

	ownersAsBytes, err := stub.GetState(ownersString)
	if err != nil {
		jsonResp := "{\"Error\": \"Failed to get owners.\"}"
		return nil, errors.New(jsonResp)
	}

	return ownersAsBytes, nil
}

//==============================================================================
//==============================================================================

func getCows(stub *shim.ChaincodeStub, cowsString string) ([]byte, error) {
	fmt.Println("Function: getCows (" + cowsString + ")")

	cowsAsBytes, err := stub.GetState(cowsString)
	if err != nil {
		jsonResp := "{\"Error\": \"Failed to get cows.\"}"
		return nil, errors.New(jsonResp)
	}

	return cowsAsBytes, nil
}

//==============================================================================
//==============================================================================

func getPolicies(stub *shim.ChaincodeStub, policiesString string) ([]byte, error) {
	fmt.Println("Function: getPolicies (" + policiesString + ")")

	policiesAsBytes, err := stub.GetState(policiesString)
	if err != nil {
		jsonResp := "{\"Error\": \"Failed to get policies.\"}"
		return nil, errors.New(jsonResp)
	}

	return policiesAsBytes, nil
}

//==============================================================================
//==============================================================================

func (t *SimpleChaincode) registerOwner(stub *shim.ChaincodeStub, args []string) ([]byte, error){
	//TODO(isaac) assign the owner a unique ID #
	//TODO(isaac) make an owner object with this ID #
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2 (First & Last name)"
	}

	var owner Owner
	owner.OwnerID = makeHash(args)

}

//==============================================================================
//==============================================================================

func (t *SimpleChaincode) registerCow(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2 (OwnerID & SensorID)")
	}

	//TODO(isaac) assign the cow a unique ID #
	//sensor ID should be argument, can be unique ID #

	ownerID = args[0]					// owner ID # should be passed as an argument
	cowID = args[1]				// the sensor ID # should be passed as an argument and will be assigned to the cow as it's unique ID #

	//TODO(isaac) check to see if the SensorID has already been registered
	// if it has, reject the attempt

	//TODO(isaac) use the CowID, OwnerID, and SensorID to make a cow object
	// convert the cow object to JSON?
	var cow Cow
	cow.ownerID = ownerID
	cow.CowID = cowID

	//TODO(isaac) pull the entire list of cows from the blockchain
	cowsAsBytes, err := getCows(stub, allCowsString)
	if err != nil {
		return nil, err
	}
	// convert it to list of cows (marshal or unmarshal?)
	//TODO(isaac) add the cow to a list of the cows that are currently covered by policies
	// convert the whole list (marshal)

	//TODO(isaac) add the list of cows to the blockchain

	err = stub.PutState(cowID, []byte(str))					//store cow with id as key
	if err != nil {
		return nil, err
	}
}

//==============================================================================
//==============================================================================

func addOwner(stub *shim.ChaincodeStub, owner Owner) error {
	fmt.Println("Function: addOwner")

	ownersAsBytes, err := getOwners(stub, OwnerString)
	if err != nil {
		return err
	}
	fmt.Println("all owners retrieved")

	var listOfOwners AllOwners
	listOfOwners, err = bytesToAllOwners(listOfOwners)
	if err != nil {
		return err
	}
	fmt.Println("all owners derived from bytes")

	listOfOwners.Catalog = append(listOfOwners.Catalog, owner)
	fmt.Println("owner appended to list of owners")

	ownersAsBytes, err = json.Marshal(listOfOwners)
	err = write(stub, OwnerString, ownersAsBytes)
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

	cowsAsBytes, err := getCows(stub, CowString)
	if err != nil {
		return err
	}
	fmt.Println("all cows retrieved")

	var listOfCows AllCows
	listOfCows, err = bytesToAllCows(listOfCows)
	if err != nil {
		return err
	}
	fmt.Println("all cows derived from bytes")

	listOfCows.Catalog = append(listOfCows.Catalog, cow)
	fmt.Println("cow appended to list of cows")

	cowsAsBytes, err = json.Marshal(listOfCows)
	err = write(stub, CowString, cowsAsBytes)
	if err != nil {
		return err
	}
	fmt.Println("list of cows successfully rewritten with new cow")
	return nil
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

func createPolicyObject(args []string) Policy {
	fmt.Println("Function: createPolicyObject")

	//TODO(isaac) I need to enter the cow object, the premium, and the value to the policy object

	var policy Policy

	policy.PolicyID = makeHash(args)
	policy.Cow = args[0]
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

	//TODO(isaac) once everything has been verified as registered, make the policy object
	newPolicy := createPolicyObject(args)

	//TODO(isaac) add the new policy to the list of current policies
	// add the policy to the owner
	policiesAsBytes, err := getPolicies(stub, incompletePoliciesString)
	if err != nil {
		return nil, err
	}

	var incompletePolicies AllPolicies
	incompletePolicies, err = bytesToAllPolicies(incompleteAsBytes)
	if err != nil {
		return nil, err
	}

	// Add the new policy to the list of pending policies
	incompletePolicies.Catalog = append(incompletePolicies.Catalog, newPolicy)
	fmt.Println("New policy appended to incomplete policies. Incomplete policy count: " + strconv.Itoa(len(incompletePolicies.Catalog)))

	err = writePolicies(stub, incompletePoliciesString, incompletePolicies)
	if err != nil {
		return nil, err
	}
	fmt.Println("incomplete policies successfully rewritten with new policy")
	return nil, nil

}

//==============================================================================
//==============================================================================

// func newOwner(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
// 	//TODO(isaac) should take new onwer's ID # and cow's ID # as arguments
//
// 	//TODO(isaac) should first check to see if owner has registered
// 	//TODO(isaac) then it should check to see if the cow is registered in the system and is attached to a policy
//
// 	//TODO(isaac) replace current ownerID in the policy with the new owners ID #
//
// }

//==============================================================================
//==============================================================================

func payOut(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	//TODO(isaac) for now - should take the sensorID as the argument

	//TODO(isaac) find what cow is associated with the sensorID
	//TODO(isaac) verify that the cowID is part of the list of current cows
	//TODO(isaac) check to see what policy is associated with that cowID
	//TODO(isaac) payout the value of the policy

	//TODO(isaac) delete the policy from the list of current policies

	//
}
