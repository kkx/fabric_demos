package main

import (
	"fmt"
	"io/ioutil"
	"os"

	//"github.com/spf13/cobra"
	"github.com/golang/protobuf/proto"
    "github.com/hyperledger/fabric/core/chaincode/platforms"
	"github.com/hyperledger/fabric/core/common/ccpackage"
	"github.com/hyperledger/fabric/core/container"
	"github.com/hyperledger/fabric/peer/common"
    //"github.com/hyperledger/fabric/protos/utils"
	pcommon "github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/peer/chaincode"
	"github.com/hyperledger/fabric/common/cauthdsl"
	pb "github.com/hyperledger/fabric/protos/peer"

)

// chaincodePackage creates the chaincode package. On success, the chaincode name
// (hash) is printed to STDOUT for use by subsequent chaincode-related CLI
// commands.
/*func chaincodePackage(cmd *cobra.Command, output_file string, cdsFact ccDepSpecFactory) error {
	if cdsFact == nil {
		return fmt.Errorf("Error chaincode deployment spec factory not specified")
	}

	var err error
	cf, err = chaincode.InitCmdFactory(false, false)
	if err != nil {
		return err
	}
	spec, err := getChaincodeSpec(cmd)
	if err != nil {
		return err
	}

	cds, err := cdsFact(spec)
	if err != nil {
		return fmt.Errorf("Error getting chaincode code %s: %s", chainFuncName, err)
	}

	var bytesToWrite []byte
	if createSignedCCDepSpec {
		bytesToWrite, err = getChaincodeInstallPackage(cds, cf)
		if err != nil {
			return err
		}
	} else {
		bytesToWrite = utils.MarshalOrPanic(cds)
	}

	fmt.Printf("Packaged chaincode into deployment spec of size <%d>, with args = %v", len(bytesToWrite), output_file)
	fileToWrite := output_file 
	err = ioutil.WriteFile(fileToWrite, bytesToWrite, 0700)
	if err != nil {
		fmt.Printf("Failed writing deployment spec to file [%s]: [%s]", fileToWrite, err)
		return err
	}

	return err
}
*/


func getInstantiationPolicy(policy string) (*pcommon.SignaturePolicyEnvelope, error) {
	p, err := cauthdsl.FromString(policy)
	if err != nil {
		return nil, fmt.Errorf("Invalid policy %s, err %s", policy, err)
	}
	return p, nil
}


// checkSpec to see if chaincode resides within current package capture for language.
func checkSpec(spec *pb.ChaincodeSpec) error {
	// Don't allow nil value
	if spec == nil {
		fmt.Printf("Expected chaincode specification, nil received")
	}

	platform, err := platforms.Find(spec.Type)
	if err != nil {
		fmt.Printf("Failed to determine platform type: %s", err)
	}

	return platform.ValidateSpec(spec)
}


// getChaincodeDeploymentSpec get chaincode deployment spec given the chaincode spec
func getChaincodeDeploymentSpec(spec *pb.ChaincodeSpec) (*pb.ChaincodeDeploymentSpec, error) {
	var codePackageBytes []byte
	var err error
	if err = checkSpec(spec); err != nil {
		return nil, err
	}

	codePackageBytes, err = container.GetChaincodePackageBytes(spec)
	if err != nil {
		err = fmt.Errorf("Error getting chaincode package bytes: %s", err)
		return nil, err
	}
    fmt.Printf("%s", codePackageBytes)
	chaincodeDeploymentSpec := &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec, CodePackage: codePackageBytes}
	return chaincodeDeploymentSpec, nil
}

func defaultCDSFactory(spec *pb.ChaincodeSpec) (*pb.ChaincodeDeploymentSpec, error) {
    return getChaincodeDeploymentSpec(spec)
}

func main(){
	var mspMgrConfigDir = "/Users/kkx/Desktop/fabric-samples/first-network/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/msp/"
	var mspID = "Org1MSP"
	err := common.InitCrypto(mspMgrConfigDir, mspID)
	if err != nil { // Handle errors reading the config file
		fmt.Printf("Cannot run peer because %s", err)
		os.Exit(1)
	}

	cf, err := chaincode.InitCmdFactory(false, false)
	if err != nil {
		fmt.Printf("%s", err)
	}

	owner := cf.Signer
	ip := "AND('" + mspID + ".admin')"
	sp, err := getInstantiationPolicy(ip)
	if err != nil {
		fmt.Printf("%s", err)
	}
    chaincodeSpec := &pb.ChaincodeSpec{Type: pb.ChaincodeSpec_GOLANG, ChaincodeId: &pb.ChaincodeID{Path: os.Args[1], Name: "testcc", Version: "0"}, Input: &pb.ChaincodeInput{Args: [][]byte{[]byte("")}}}
    packageBytes, err := ioutil.ReadFile(os.Args[1])
    if err != nil {
		fmt.Printf("%s", err)
    }

    fmt.Printf("%s\n", os.Args[1])
    fmt.Printf("%s\n", packageBytes)

    cds, err := defaultCDSFactory(chaincodeSpec )
    if err != nil {
        fmt.Printf("Error getting chaincode code %s", err)
    }


    var objToWrite proto.Message
	objToWrite, err = ccpackage.OwnerCreateSignedCCDepSpec(cds, sp, owner)
	if err != nil {
		fmt.Errorf(" %s", err)
	}

	bytesToWrite, err := proto.Marshal(objToWrite)
	if err != nil {
		fmt.Errorf("Error marshalling chaincode package : %s", err)
	}

	fmt.Printf("Packaged chaincode into deployment spec of size <%d>, with args = %v", len(bytesToWrite), os.Args[2])
	fileToWrite := os.Args[2] 
	err = ioutil.WriteFile(fileToWrite, bytesToWrite, 0700)
	if err != nil {
		fmt.Printf("Failed writing deployment spec to file [%s]: [%s]", fileToWrite, err)
	}
}

