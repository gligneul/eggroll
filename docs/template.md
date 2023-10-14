---
title: Getting Started
---

Getting Started
===============

## Dependencies

The first step to use EggRoll is to install the necessary dependencies.

- [Go](https://go.dev/doc/install)
- [Docker](https://www.docker.com/products/docker-desktop/)
- [Sunodo](https://docs.sunodo.io/guide/introduction/installing)

## Template

The [EggRoll template repository](https://github.com/gligneul/eggroll-template/) provides a minimal working example of a EggRoll project.
The box below presents the most important files of the template.

<!---
tree ./examples/template
-->
```
./template
├── client
│   └── main.go
├── dapp
│   └── main.go
├── Dockerfile
├── go.mod
├── go.sum
└── integration_test.go
```

- The client directory contains the code for DApp client.
  In there, we define the main package of the client, which runs off-chain.

- The dapp direcoty contains the code for the DApp contract.
  We define another main package in this directory, for the binary that runs on-chain.

- The Dockerfile describes the Docker image used by Sunodo.
  This file

- The go.mod and go.sum files contain
