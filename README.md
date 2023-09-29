# EggRoll 🥢

A high-level, opinionated, lambda-based framework for Cartesi Rollups in Go

## Requirements

EggRoll is built on top of the [Cartesi Rollups](https://docs.cartesi.io/cartesi-rollups/) infrastructure version 1.0.
To use EggRoll, you also need [sunodo](https://github.com/sunodo/sunodo/) version 0.8.

## Quick Look

Let's look at a simple example: a DApp that keeps a text box in the blockchain.

In EggRoll, you should first define the types shared between the front and back end.
These types are the inputs that advance the rollups and the backend state.

```go
// Append a value to the text box
type Append struct {
	Value string
}

// Clear the text box
type Clear struct {
}

// Text box shared state
type State struct {
	TextBox string
}
```

Then, you should use the `eggroll.DApp` interface to build the DApp backend.
`Handle` registers a function for each input struct type.
Each handler receives the the mutable state and the respective input struct.
After registering all functions, call `Roll` to start the DApp backend.

```go
dapp := eggroll.SetupDApp[State]()

eggroll.Register(dapp, func(_ eggroll.Env, state *State, _ *Clear) error {
    state.TextBox = ""
    return nil
})

eggroll.Register(dapp, func(_ eggroll.Env, state *State, input *Append) error {
    state.TextBox += input.Value
    return nil
})

log.Panic(dapp.Roll())
```

Finally, you can interact with the DApp from the front end using the `eggroll.Client` interface.
`Send` sends inputs to the DApp backend through the blockchain in a single transaction.
`WaitFor` waits for the given input to be processed by the DApp backend.
`State` reads the DApp backend state from the Cartesi reader node.

```go
client := eggroll.SetupClient[State]()

indices, err := client.Send(
    &Clear{},
    &Append{Value: "egg"},
    &Append{Value: "roll"},
)
if err != nil {
    log.Panic(err)
}

lastInput := indices[len(indices)-1]
if err := client.WaitFor(lastInput); err != nil {
    log.Panic(err)
}

state := client.State()
fmt.Println(state.TextBox) // -> eggroll
```

To run this example, check the README in `./examples/textbox`.
