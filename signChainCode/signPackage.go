package main

import (
	"fmt"
	"io/ioutil"
	"os"

	//"github.com/spf13/cobra"
	"github.com/hyperledger/fabric/core/common/ccpackage"
	"github.com/hyperledger/fabric/peer/common"
	"github.com/hyperledger/fabric/protos/utils"
	"github.com/hyperledger/fabric/peer/chaincode"
)

func signpackage(ipackageFile string, opackageFile string) error {
	var err error
	cf, err := chaincode.InitCmdFactory(false, false)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(ipackageFile)
	if err != nil {
		return err
	}

	env := utils.UnmarshalEnvelopeOrPanic(b)

	env, err = ccpackage.SignExistingPackage(env, cf.Signer)
	if err != nil {
		return err
	}

	b = utils.MarshalOrPanic(env)
	err = ioutil.WriteFile(opackageFile, b, 0700)
	if err != nil {
		return err
	}

	fmt.Printf("Wrote signed package to %s successfully\n", opackageFile)

	return nil
}

func main(){
	var mspMgrConfigDir = "/Users/kkx/Desktop/fabric-samples/first-network/crypto-config/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/msp/"
	var mspID = "Org2MSP"
	err := common.InitCrypto(mspMgrConfigDir, mspID)
	if err != nil { // Handle errors reading the config file
		fmt.Printf("Cannot run peer because %s", err)
		os.Exit(1)
	}
	signpackage(os.Args[1], os.Args[2])
}

