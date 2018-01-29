package main

import (
    "fmt"
    "os"
    "github.com/CognitionFoundry/gohfc" 
	"context"
)


// FabricSetup implementation

const (
	channelID = "mychannel"
	orgName   = "Org1"
	orgAdmin  = "Admin"
	ccID      = "mycc1"
	ChannelConfig = "../channel-artifacts/channel.tx"
)

func main() {
    
    client, err := gohfc.NewFabricClient("./myConfig.yaml")
    if err != nil {
        fmt.Printf("Error loading file: %v", err)
        os.Exit(1)
    }

    pk:="../crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/admincerts/Admin@org1.example.com-cert.pem"
    sk:="../crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/924e28b054f0d0c0cca0efcdded80356e70a3cc888c0ba2576d82483c477e0e0_sk"

    admin, err := gohfc.LoadCertFromFile(pk,sk)
    if err != nil {
        fmt.Println(err)
    }


    channel := &gohfc.Channel{
        MspId:       "Org1MSP",
        ChannelName: "mychannel",
    }

	chaincode := &gohfc.ChainCode{
		Channel: channel,
		Type:    gohfc.ChaincodeSpec_GOLANG,
		Name:    "mycc1",
		Version: "1.0",
		Args:    []string{"query", "b"},
	}

	// identity is identity returned from ca.Enroll or Admin
    result, err := client.Query(admin,  chaincode, []string{"peer0"})

	if err != nil {
		fmt.Print(err)
	}
    fmt.Print(result[0])

    
	chaincodeInvoke := &gohfc.ChainCode{
		Channel: channel,
		Type:    gohfc.ChaincodeSpec_GOLANG,
		Name:    "mycc1",
		Version: "1.0",
		Args:    []string{"invoke", "b", "a", "10"},
	}

	// identity is identity returned from ca.Enroll or Admin
	_, err = client.Invoke(admin, chaincodeInvoke, []string{"peer0"}, "orderer0")
	if err != nil {
		fmt.Print(err)
	}

	// identity is identity returned from ca.Enroll or Admin
    result, err = client.Query(admin,  chaincode, []string{"peer0"})

	if err != nil {
		fmt.Print(err)
	}
    fmt.Print(result[0])


    fmt.Println("from event")
    ch:=make(chan gohfc.BlockEventResponse)
    ctx,_:=context.WithCancel(context.Background())
    err=client.Listen(ctx, admin, "peer0", "Org1MSP", ch)
    for d:= range ch{
        fmt.Println(d.TxID)
        fmt.Println(d.ChainCodeName)
        fmt.Println(d.CCEvents[0].EventName)
        fmt.Printf("%s\n", d.CCEvents[0].EventPayload)
    }
    fmt.Println("event")


}
    
