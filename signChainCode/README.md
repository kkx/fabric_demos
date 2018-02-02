# Instructions for testing package and sign chaincodes scripts

This demo use the same fabric code as our library. no sdk is used.

For package chaincode with creator's signature: 
- change variables mspMgrConfigDir and mspID to set a signer of your own environment/fabric network
- go build packagePackage.go
- add test code folder chaincode to your gopath
- ./packagePackage chaincode ccpack.out

For sign signed chaincode package with a new owner signature: 
- change variables mspMgrConfigDir and mspID to set a new signer of your own environment/fabric network
- go build signPackage.go
- ./packagePackage ccpack.out ccpackNew.out


You can test ccpackNew.out using an [extractor](https://github.com/jlgarciasan/hyperledger-utils)


