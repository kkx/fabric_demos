package main

import (
    //"path"
    //"strconv"
	//"testing"
	//"time"
    "fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn"

	"github.com/hyperledger/fabric-sdk-go/pkg/config"
	//"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	//chmgmt "github.com/hyperledger/fabric-sdk-go/api/apitxn/chmgmtclient"
	//fcutil "github.com/hyperledger/fabric-sdk-go/pkg/util"

	//"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/events"
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
    

	// Initialize the configuration
	// This will read the config.yaml, in order to tell to
	// the SDK all options and how contact a peer
	c := config.FromFile("config.yaml")
	//if err != nil {
    //  fmt.Errorf("Initialize the config failed: %v", err)
    //}
	sdk, err := fabsdk.New(c)
	if err != nil {
		fmt.Errorf("Create client failed: %v", err)
	}	


    // Channel client is used to query and execute transactions
    _, err = sdk.NewChannelClient(channelID, "User1")
    if err != nil {
        fmt.Println("hello world")
    }
    fmt.Println("2")
    
	// Channel client is used to query and execute transactions (Org1 is default org)
	chClient, err := sdk.NewClient(fabsdk.WithUser("Admin")).Channel(channelID)
	if err != nil {
		fmt.Errorf("Failed to create new channel client: %s", err)
	}

	// Release all channel client resources
	//defer chClient.Close()
    fmt.Println("3")
    var queryArgs = [][]byte{[]byte("b")}

    transientData := "Some data"
    transientDataMap := make(map[string][]byte)
    transientDataMap["result"] = []byte(transientData)

    value, err := chClient.Query(apitxn.Request{ChaincodeID: ccID, Fcn: "query", Args: queryArgs })
	if err != nil {
		fmt.Errorf("Failed to query funds: %s", err)
	}
	fmt.Printf("%s\n", value)

    eventID := "evtSender"
    notifier := make(chan *apitxn.CCEvent)
    rce := chClient.RegisterChaincodeEvent(notifier, ccID, eventID)

    var invokeArgs = [][]byte{[]byte("a"), []byte("b"), []byte("10")}
	value, _, err = chClient.Execute(apitxn.Request{ChaincodeID: ccID, Fcn: "invoke", Args: invokeArgs })

	
    for ccEvent:= range notifier{
        fmt.Printf("Received CC event: %s\n", ccEvent)
        break 
    }

    /*
	select {
		case ccEvent := <-notifier:
			fmt.Printf("Received CC event: %s\n", ccEvent)
		case <-time.After(time.Second * 20):
			fmt.Printf("Did NOT receive CC event for eventId(%s)\n", eventID)
	}
    */


	
	// Unregister chain code event using registration handle
	err = chClient.UnregisterChaincodeEvent(rce)
	if err != nil {
		fmt.Printf("Unregister cc event failed: %s", err)
	}

    value, err = chClient.Query(apitxn.Request{ChaincodeID: ccID, Fcn: "query", Args: queryArgs })
	if err != nil {
		fmt.Errorf("Failed to query funds: %s", err)
	}
	fmt.Printf("%s\n", value)
}

    
