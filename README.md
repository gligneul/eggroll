# EggRoll 🐣🛼

![Build](https://github.com/gligneul/eggroll/actions/workflows/go.yml/badge.svg)
[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/gligneul/eggroll)

A high-level framework for Cartesi Rollups in Go.

EggRoll divides a Cartesi DApp into two parts: contract and client.
The contract runs in the blockchain inside a Cartesi VM and relies on the Cartesi Rollups API.
The client runs off-chain and communicates with the contract using the Cartesi Reader Node APIs and the Ethereum API.
EggRoll provides abstractions for both sides of the DApp.

## Requirements

EggRoll is built on top of the [Cartesi Rollups](https://docs.cartesi.io/cartesi-rollups/) infrastructure version 1.0.
To use EggRoll, you also need [sunodo](https://github.com/sunodo/sunodo/) version 0.9.

## Quick Look

The first step to using EggRoll is defining the contract struct that runs inside the Cartesi VM.
This struct should use the `eggroll.DefaultContract` to implement the optional methods of the `eggroll.Contract` interface.
The only obligatory method is the advance one, which receives the rollup environment and the input.
In the example below, the advance method logs and returns the input.
Then, you should call the `eggroll.Roll` function, passing the contract to start the rollup's main loop.

<!---
cat ./examples/template/dapp/main.go
-->
```
type TemplateContract struct {
	eggroll.DefaultContract
}

func (c *TemplateContract) Advance(env *eggroll.Env, input any) ([]byte, error) {
	inputStr := string(input.([]byte))
	env.Logf("received: %v", inputStr)
	return []byte("echo: " + inputStr), nil
}

func main() {
	eggroll.Roll(&TemplateContract{})
}
```

Off-chain, you can use the `eggroll.DevClient` struct to interact with the contract.
The example below reads the command line's first argument to use as input.
The `client.SendInputBytes` function sends this input to the blockchain.
The `client.WaitFor` reads the rollups node, waiting until it processes the given input.
Once the node processes the input, the code prints the result from the contract.

<!---
cat ./examples/template/client/main.go
-->
```
func main() {
	input := os.Args[1]
	ctx := context.Background()
	client, _ := eggroll.NewDevClient()
	inputIndex, _ := client.SendInputBytes(ctx, []byte(input))
	result, _ := client.WaitFor(ctx, inputIndex)
	fmt.Println(string(result.Result))
}
```

To run this example or start your project, check out the [template repository](https://github.com/gligneul/eggroll-template/).
For more examples, check the examples directory in this repo.
