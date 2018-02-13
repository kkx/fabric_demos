# Execute Lua in chaincode fabric go sdk

In this demo, we run lua script in chaincodes. we use [fabric go sdk](https://github.com/hyperledger/fabric-sdk-go) to connect to our chaincode and send lua code as payload and write the lua execution output in the blockchain ledger.

# Installation

  - Setup an private hyperledger network, in this case we use [first-network](https://github.com/hyperledger/fabric-samples/tree/release/first-network) 
    - ./byfn.sh -m up # create the network with default chaincode installed
    - Modify docker-compose-cli.yaml and mount a volume for chainCodeExample.go in your container.
    - Run `docker exec -it cli` bash # enter into the cli console
    - Run following commands to install customized chaincode(with events):
```
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
CORE_PEER_ADDRESS=peer0.org1.example.com:7051
CORE_PEER_LOCALMSPID="Org1MSP"
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
peer chaincode install -n mycc1 -v 1.0 -p $PATH_TO_chainCodeExample.go
peer chaincode instantiate -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n mycc1 -v 1.0 -c '{"Args":["init","a", "100", "b","200"]}' -P "OR ('Org1MSP.member','Org2MSP.member')"
```
- Install fabric go sdk following the instructions from their github repository.
- Install [go lua interpreter](https://github.com/yuin/gopher-lua) in your peers where where the chaincode are installed and instantiated
- Compile setup.go and run it with myConfig.yaml
- You should receive from the lua code output in you console
