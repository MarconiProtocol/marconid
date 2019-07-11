## Background
Marconi Protocol and the Marconi Blockchain facilitates secure network communication, flexible network infrastructure, and the formation of mesh networks. Please see [White-Paper](https://docsend.com/view/5zragmb) for more details about Marconi.

On a high-level, Marconi consists of 4 major components:

*  **Marconid**: the Marconi daemon that interface with the operating system's networking functionality and handles the networking related interactions between the node and its peers on the Marconi networks.
*  **Go-Marconi**: the Marconi blockchain, responsible for finding peers through the bootnodes and connects and sync'ing blocks from the Global Chain.
*  **Middleware**: acts as a bridge between Go-Marconi and Marconid by exposing JSON RPC APIs. 
*  **CLI**: command-line interface to all of the Marconi components, allowing the user to run commands and interact with the Marconi Network. It's also responsible for ensuring all of the required packages are downloaded/updated and any configurations are executed as needed.
 
For more information about the Marconi Architecture and the major components, please refer to [Architecture Overview](https://github.com/MarconiProtocol/wiki/wiki/Architecture-Overview).

## Description
Marconid (Marconi Daemon) creates and manages mPipe connections between nodes. This includes interfacing with the operating system’s networking functionality, establishing peer connections, and handling packet authentication, encryption and decryption. Marconid also contains a processor that executes programmable packets.

*  **core/blockchain** - interface with the blockchain and handles peers updates via RPC calls
*  **core/config** - static and runtime configurations
*  **core/crypto** - cryptography logic such as the Diffie–Hellman key exchange
*  **core/net/core** - core networking related logic such as bridge, tunnel, tap connections
*  **core/net/dht** - marconid's distributed-hash-table for peer discovery
*  **core/peer** - encapsulate peer management functions such as add/remove/update
*  **core/sys** - system level commands interfacing with the operating system

## Pre-requisite
Golang should be already installed and GOPATH/GOROOT are configured.

## Get dependencies
`go get ./...`

## Build Marconid
`go build -o marconid agent/service/main.go`

## Run Tests
Pre-requisite: Golang should be already installed and GOPATH/GOROOT are configured.

`cd util`

`./test.sh`

## Contribution
Marconi is an open-source project and we welcome the community to contribute.  To get involved, please review our [guide to contributing](CONTRIBUTING.md).

## License
[aGPL](LICENSE)

## Other Resources
[Marconi Website](https://marconi.org)

[Blockchain Explorer](https://explorer.marconi.org)
