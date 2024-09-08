
# Chia-CADT

## FIles

The Chia-CADT folder includes:
- A `Docker` folder containing all Dockerfiles and scripts necessary to set up two instances of CADT running on a Chia node. It also includes the `chia-blockchain` 2.3.1 files required by the Dockerfile.
- `mock_cadt_data.sql` for creating a mock CADT SQL database and `cadtexample.json` with commands to create projects and units on CADT.

## Docker

The Docker setup currently runs two instances of **Chia-CADT** on Chia testnet11, utilizing the latest version from the [CADT develop branch (default branch)](https://github.com/Chia-Network/cadt).

In the `docker-compose.yml` file, there are two configuration options:
1. The first configuration (commented out) deploys only the Chia-CADT instances.
2. The second configuration deploys both **Chia-CADT** instances and dedicated **MySQL databases** for each instance.

In the `modify_config` file, there are two potential CADT configurations:
1. The first configuration (commented out) does not configure CADT to use a database mirror.
2. The second configuration enables CADT to use a **MySQL database mirror**. This configuration can be adjusted to either mirror databases created within the Docker Compose environment via environment variables (same for both MySQL and Chia-CADT) or mirror an external database.

Due to a bug with units data tables CADT Mirror Databses must be configured with `sql_mode="NO_ENGINE_SUBSTITUTION"`.

## CADT API

To properly interact with the **CADT API**, you must have a funded **Chia wallet**:

1. Obtain the wallet’s **Master Public Key** by running the command:
   ```bash
   chia keys show
   ```

2. Fund the wallet using a faucet:
   - Visit the [Testnet 11 Faucet](https://testnet11-faucet.chia.net/) to obtain testnet funds.

### Interacting with the CADT API

Once the wallet is funded, you can interact with the CADT API by following the documentation provided here:
- [CADT RPC API Documentation](https://github.com/Chia-Network/cadt/blob/develop/docs/cadt_rpc_api.md)

When creating an organization in CADT, ensure that **mirrors** are set up for its four respective data stores using the following command:
```bash
chia data add_mirror -i <Store ID> -a 30000 -u [ip:8575] -m .0003
```

An example of project and unit creation is provided in `(cadt_example.json)`

## Integration with Carbon Market Bridge

The **CADT MySQL Database Mirror(s)** can be used as the data source for **Carbon Market Bridge's chaincode** and **application**.

Alternatively, you can use a mock CADT MySQL Mirror Database to experiment with Carbon Market Bridge. SQL commands to create mock data for testing purposes are provided in the file `(mock_cadt_data.sql)`.

# Hyperledger Fabric Setup for Carbon Market Bridge

## Installation

To run the **Carbon Market Bridge** on the Hyperledger Fabric test network, ensure that the following prerequisites are met:

1. **Test Network Prerequisites**: Install the required dependencies as described in the [Hyperledger Fabric Prerequisites](https://hyperledger-fabric.readthedocs.io/en/latest/prereqs.html).
   
2. **Fabric Binaries, Docker Images, and Test Network**: Install the Fabric binaries, Docker images, and the Fabric test network (part of the Fabric Samples) as outlined in the [Hyperledger Fabric Installation Guide](https://hyperledger-fabric.readthedocs.io/en/latest/install.html).

Note that Hyperledger Fabric currently does not support Windows, WSL 2 must be used in this instance.

## Versions

As of this commit, **Carbon Market Bridge** works with the following exact versions of Hyperledger Fabric components for reference. Ensure your setup matches these versions for consistency and compatibility.

### Docker Images:
```bash
$ docker images | grep hyperledger
hyperledger/fabric-peer             2.5.9
hyperledger/fabric-orderer          2.5.9
hyperledger/fabric-ccenv            2.5.9
hyperledger/fabric-baseos           2.5.9
hyperledger/fabric-ca               1.5.12
```

### Binaries:
```bash
$ peer version
peer:
 Version: v2.5.9

$ orderer version
orderer:
 Version: v2.5.9

$ fabric-ca-client version
fabric-ca-client:
 Version: v1.5.12

$ fabric-ca-server version
fabric-ca-server:
 Version: v1.5.12
```

Ensure that your setup matches these versions for consistency in your environment.

# Carbon Market Bridge

## Files

The Carbon Market Bridge folder includes:
- Packaged `unitTokenChaincode` chaincode which features all of the on-chain logic for the CMB platform to be deployed on a Hyperledger Fabric network.
- Two instances of the application for the CMB platform, configured for simultaneous use by Org1 and Org2 clients of the Hyperledger Fabric test-network.

## Carbon Market Bridge Chaincode

### Chaincode MySQL Configuration

For the chaincode to initialize, it requires a connection to the default CADT MySQL Database Mirror. Modify the `unitTokenChaincode/unittoken-chaincode.go` file accordingly:

```go
func initializeDB() error {
    cfg := mysql.Config{
        User:                 "<YOUR_DB_USER>",
        Passwd:               "<YOUR_DB_PASSWORD>",
        Net:                  "tcp",
        Addr:                 "<YOUR_DB_HOST>:<YOUR_DB_PORT>",
        DBName:               "<YOUR_DATABASE_NAME>",
        AllowNativePasswords: true,
    }
}
```

### Chaincode Deployment

The chaincode must be in the same directory as the test-network.

Use the following commands to bring up the test network and deploy the chaincode:

```bash
./network.sh up
./network.sh createChannel
./network.sh deployCC -ccn unittokenchaincode -ccp ../chaincode/unitTokenChaincode -ccl go
```

Use the following command to bring down the test network:

```bash
./network.sh down
```

---

### Payment Token Commands

The CMB application currently does not interface with the payment token functions in the chaincode, commands must therefore be used, examples of the primary commands as org1 of the test-network are given below:

Set client:

```bash
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
```

Create Token:

```bash
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n unittokenchaincode --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c '{"function":"CreateToken","Args":["Token X"]}'
```

Mint Tokens:

```bash
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n unittokenchaincode --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c '{"function":"MintTokens","Args":["Token X", "10000"]}'
```

Transfer Minted Tokens To Payment Token Creator:

```bash
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n unittokenchaincode --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c '{"function":"TransferMintedTokens","Args":["Token X", "10000"]}'
```

Modify the ports to 9050 and mentions of org1 to org2 to interact with the chaincode as org2.

## Carbon Market Bridge Application

### Base Dependencies

#### Install npm:

```bash
sudo apt install npm
```

#### Install Redis:

```bash
sudo apt install redis-server
```

#### Install NVM and Node.js:

```bash
# Install NVM (Node Version Manager)
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.3/install.sh | bash

# Source the NVM script
export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
[ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"

# Install and use Node.js version 20.16.0
nvm install 20.16.0
nvm use 20.16.0
```

#### Install Python and pip:

To install Python and pip, you can use the following commands:

```bash
sudo apt install python3
sudo apt install python3-pip
```

#### Install Python libraries:

Once Python and pip are installed, you can install necessary Python packages:

```bash
pip3 install numpy
pip3 install statsmodels
```

### Configuration

#### Fabric Gateway Configuration:
To configure the application’s Fabric Gateway to connect to a local peer, modify the endpoint, aliases, and certificate directory fields in `settings.json` accordingly. Each organization’s application instance is currently set up to connect to its respective test-network peer, so you only need to replace `<YOUR_HOME_DIRECTORY>` with the correct path.

#### MySQL Database Configuration:
In `dbConfig.ts`, replace the placeholders with your actual MySQL database details:

- `<YOUR_DB_HOST>`: The host address of your MySQL database.
- `<YOUR_DB_USER>`: The username for your MySQL database.
- `<YOUR_DB_PASSWORD>`: The password for your MySQL database.
- `<YOUR_DATABASE_NAME>`: The name of your database (e.g. `CADTDatabase`).

---

### Installation

By default, the application is configured to run on **port 3000** for Org1 and **port 4000** for Org2. You can adjust these settings by modifying the port references in `app.ts` and the frontend files located in the `public` directory.

To install dependencies, compile, and run the application, use the following commands:

```bash
# Install all dependencies (generates node_modules)
npm install

# Compile TypeScript to JavaScript (generates dist folder)
npx tsc

# Start the application
npm start

```

### Commands

If the Redis server cache needs to be reset such as in cases where the Hyperledger Fabric Network is redeployed, the following command can be used.

```bash
redis-cli FLUSHALL
```