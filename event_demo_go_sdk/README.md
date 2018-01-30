# Get chaincode event using fabric go sdk

In this demo, we customize chaincodes to introduce events in it. After that, we use [fabric go sdk](https://github.com/hyperledger/fabric-sdk-go) to connect to our chaincode and read the events trigged by invokes.

# Installation!

  - Setup an private hyperledger network, in this case we use [first-network](https://github.com/hyperledger/fabric-samples/tree/release/first-network) 
  -- ./byfn.sh -m up # create the network with default chaincode installed
  -- Modify docker-compose-cli.yaml and mount a volumne for chainCodeExample.go in your container.
  -- docker exec -it cli bash # enter into the cli console
  -- run following commands:
```
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
CORE_PEER_ADDRESS=peer0.org1.example.com:7051
CORE_PEER_LOCALMSPID="Org1MSP"
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
peer chaincode install -n mycc1 -v 1.0 -p $PATH_TO_chainCodeExample.go
peer chaincode instantiate -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n mycc1 -v 1.0 -c '{"Args":["init","a", "100", "b","200"]}' -P "OR ('Org1MSP.member','Org2MSP.member')"
```
- Install fabric go sdk following the instructions from their github repository.
- compile setup.go and run it with config.yaml


