---
title: Getting Started
---

Getting Started
=

# Dependencies

The first step to using EggRoll is to install the necessary dependencies.

- [Go](https://go.dev/doc/install)
- [Docker](https://www.docker.com/products/docker-desktop/)
- [Sunodo](https://docs.sunodo.io/guide/introduction/installing)

# Template

The [EggRoll template repository](https://github.com/gligneul/eggroll-template/) provides a minimal working example of an EggRoll project.
The box below presents the most important files of the template.

<!---
tree -a ./examples/template
-->
```
./eggroll-template
├── client
│   └── main.go
├── contract
│   └── main.go
├── Dockerfile
├── .dockerignore
├── .gitignore
├── go.mod
├── go.sum
└── integration_test.go
```

The `client` directory contains the code for the DApp client.
Over there, we define the main package of the client, which runs off-chain.
This package uses the `eggroll.Client` struct to interact with the DApp contract.
Go to the [client section](/client) for more details.

The `contract` directory contains the code for the DApp contract.
We define another main package in this directory for the binary that runs on-chain inside the Cartesi Machine.
This package defines the struct that implements the `eggroll.Contract` interface.
Go to the [contract section](/contract) for more details.

The `Dockerfile` uses the sunodo infrastructure to build the DApp contract image. The `.dockerignore` and `.gitignore` files ignore the `.sunodo` directory.

The `go.mod` and `go.sum` files define the template Go modules.
You should rename the module to your project name when cloning this template.

Finally, the `integration_test.go` file contains an integration test for this module.
For more details, go to the [testing section](/testing).

# Running the DApp

To run the DApp, first you need to build it with sunodo.

```sh
$ sunodo build
```

Then, you can use sunodo to run it.

```
$ sunodo run
```

Once you see the message `Press Ctrl+C to stop the node`, you can open another terminal and send an input to the DApp.

```
$ go run ./client hi
hi
```
