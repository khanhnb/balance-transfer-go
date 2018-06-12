## Balance transfer go

A sample go app to demonstrate **__fabric-client__** & **__fabric-ca-client__** go SDK APIs

### Prerequisites and setup:

* [Docker](https://www.docker.com/products/overview) - v1.12 or higher
* [Docker Compose](https://docs.docker.com/compose/overview/) - v1.8 or higher
* [Git client](https://git-scm.com/downloads) - needed for clone commands
* **Golang**
* [Download Docker images](http://hyperledger-fabric.readthedocs.io/en/latest/samples.html#binaries)
```
go get github.com/balance-transfer-go
cd $GOPATH/src/github.com/balance-transfer-go
sudo echo "127.0.0.1 ca.org1.example.com ca.org2.example.com orderer.example.com peer0.org1.example.com peer1.org1.example.com peer0.org2.example.com peer1.org2.example.com" >> /etc/hosts
```


Once you have completed the above setup, you will have provisioned a local network with the following docker container configuration:

* 2 CAs
* A SOLO orderer
* 4 peers (2 peers per Org)

#### Artifacts
* Crypto material has been generated using the **cryptogen** tool from Hyperledger Fabric and mounted to all peers, the orderering node and CA containers. More details regarding the cryptogen tool are available [here](http://hyperledger-fabric.readthedocs.io/en/latest/build_network.html#crypto-generator).
* An Orderer genesis block (genesis.block) and channel configuration transaction (mychannel.tx) has been pre generated using the **configtxgen** tool from Hyperledger Fabric and placed within the artifacts folder. More details regarding the configtxgen tool are available [here](http://hyperledger-fabric.readthedocs.io/en/latest/build_network.html#configuration-transaction-generator).

## Running the sample program

##### Terminal Window 1

* Launch the network using docker-compose

```
docker-compose -f artifacts/docker-compose.yaml up -d
```
##### Terminal Window 2

* Install go REST Service

```
go run main.go
```

##### Terminal Window 3

* Execute the REST APIs from the section [Sample REST APIs Requests](https://github.com/hyperledger/fabric-samples/tree/master/balance-transfer#sample-rest-apis-requests)

## Sample REST APIs Requests

### Login Request

* Register and enroll new users in Organization - **Org1**:

`curl -s -X POST http://localhost:4000/users -H "content-type: application/x-www-form-urlencoded" -d 'username=Jim&orgName=org1'`

**OUTPUT:**

```
{
  "Success": true,
  "Message": "Jim enrolled Successfully",
  "Token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTQ4NjU1OTEsInVzZXJuYW1lIjoiSmltIiwib3JnTmFtZSI6Im9yZzEiLCJpYXQiOjE0OTQ4NjE5OTF9.yWaJhFDuTvMQRaZIqg20Is5t-JJ_1BP58yrNLOKxtNI"
}
```

The response contains the success/failure status, an **enrollment Secret** and a **JSON Web Token (JWT)** that is a required string in the Request Headers for subsequent requests.

### Invoke request

```
curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTQ4NjU1OTEsInVzZXJuYW1lIjoiSmltIiwib3JnTmFtZSI6Im9yZzEiLCJpYXQiOjE0OTQ4NjE5OTF9.yWaJhFDuTvMQRaZIqg20Is5t-JJ_1BP58yrNLOKxtNI" \
  -H "content-type: application/json" \
  -d '{
	"fcn":"invoke",
	"args":["a","b","10"]
}'
```
**NOTE:** Ensure that you save the Transaction ID from the response in order to pass this string in the subsequent query transactions.

### Chaincode Query

```
curl -s -X GET \
  "http://localhost:4000/channels/mychannel/chaincodes/mycc?fcn=query&args=%5B%22a%22%5D" \
  -H "authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTQ4NjU1OTEsInVzZXJuYW1lIjoiSmltIiwib3JnTmFtZSI6Im9yZzEiLCJpYXQiOjE0OTQ4NjE5OTF9.yWaJhFDuTvMQRaZIqg20Is5t-JJ_1BP58yrNLOKxtNI" \
  -H "content-type: application/json"
```

### Query Block by BlockNumber

```
curl -s -X GET \
  "http://localhost:4000/channels/mychannel/blocks/1?peer=peer0.org1.example.com" \
  -H "authorization: <put JSON Web Token here>" \
  -H "content-type: application/json"
```

### Query Transaction by TransactionID (updating)


### Query ChainInfo (updating)


### Query Installed chaincodes (updating)


### Query Instantiated chaincodes (updating)


### Query Channels (updating)

