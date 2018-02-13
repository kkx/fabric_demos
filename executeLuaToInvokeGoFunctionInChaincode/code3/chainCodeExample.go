/*
Copyright IBM Corp. 2016 All Rights Reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
		 http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"fmt"
    "encoding/json"
	"github.com/yuin/gopher-lua"
    "net/http"
	//"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type customEvent struct {
    Type        string `json:"type"`
    Description string `json:"description"`
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Init")
	_, args := stub.GetFunctionAndParameters()
    fmt.Printf("%s", args)
	var err error

	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Expecting 0")
	}

	// Write the state to the ledger
	err = stub.PutState("LuaResult", []byte("hola"))
    var event = customEvent{"putState", "Successfully put state Lua code: empty string"}
    eventBytes, err := json.Marshal(&event)
    err = stub.SetEvent("evtSender", eventBytes)

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "invoke" {
		// execute lua code in chaincode
		return t.invoke(stub, args)
	} else if function == "query" {
		// the old "Query" is now implemtned in invoke
		return t.query(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting at least 1")
	}

    luaFuncCode := args[0]
    L := lua.NewState()
	defer L.Close()
    L.SetGlobal("ServiceCall", L.NewFunction(ServiceCall)) 
	if err := L.DoString(luaFuncCode); err != nil {
	    panic(err)
	}

    if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("execute"), // name of Lua function
		NRet:    1,     // number of returned values 
		Protect: true,                  // return err or panic
    }); err != nil {
		panic(err)
	}

	// Get the returned value from the stack and cast it to a lua.LString
    luaFuncResult, ok := L.Get(-1).(lua.LString);
	if ok {
		//fmt.Println(luaFuncResult)
	}

	// save the result
	// Write the state back to the ledger
	err = stub.PutState("LuaResult", []byte(luaFuncResult))

    /*
    var event = customEvent{"putState", "Successfully put state lua func result: " + string(luaFuncResult)}
    eventBytes, err := json.Marshal(&event)
    err = stub.SetEvent("evtSender", eventBytes)
    */

	if err != nil {
		return shim.Error(err.Error())
	}
    jsonResp := "{" + "LuaResult\":\"" + string(luaFuncResult) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)

	return shim.Success([]byte(jsonResp))
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    // Do nothing
    return shim.Success(nil)
}

// query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	// Get the state from the ledger
	LuaResult, err := stub.GetState("LuaResult")
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + "LuaResult"+ "\"}"
		return shim.Error(jsonResp)
	}

	if LuaResult == nil {
		jsonResp := "{\"Error\":\"Nil amount for LuaResult\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{" + "LuaResult\":\"" + string(LuaResult) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success([]byte(jsonResp))
}

func ServiceCall(L *lua.LState) string {
    url := L.ToString(1)
    method := L.ToString(2)
    if method == 'GET' {
        response, err := http.Get(url)
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        L.push(contents)
    }
    else if method == 'POST' {
        response, err := http.Get(url)
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        L.push(contents)
    }
    return 'ok'
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
