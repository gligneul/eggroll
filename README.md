# EggRoll

A high-level, opinionated, map-reduce-based framework for Cartesi Rollups in Go

## Quick Look

Let's look at a simple example: a DApp that receives and accumulates inputs.

In EggRoll, you should first define the types shared between the front and back end.
These types are the input that advances the rollups and the backend state.

```go
type Input struct {
	Value string
}

type State struct {
	Accumulator string
}
```

Then, you should define the reduce function that will run in the DApp backend.
This function receives the current state and the next input that should be processed.
The `eggroll.Roll` function will process the rollups and call the reduce function for each input.

```go
func reduce(e eggroll.Environment, s *acc.State, i *acc.Input) error {
	s.Accumulator += i.Value
	return nil
}

func main() {
	eggroll.Roll(reduce)
}
```

Finally, you can interact with the DApp from the front end using the `eggroll.Client` interface.
The `c.Send` call sends an input to the DApp backend through the blockchain.
The `c.WaitFor` call waits for the input to be processed by the DApp backend.
The `c.Read` call reads the state from the backend.

```go
c := setupClient()
idx, err := c.Send(&acc.Input{Value: "hello"})
if err != nil {
    log.Panic(err)
}

idx, err = c.Send(&acc.Input{Value: ", world"})
if err != nil {
    log.Panic(err)
}

if err := c.WaitFor(idx); err != nil {
    log.Panic(err)
}

var state acc.State
if err := c.Read(&state); err != nil {
    log.Panic(err)
}

fmt.Println(state.Accumulator) // -> hello, world
```

## Requiriments

To use EggRoll, you need [sunodo](https://github.com/sunodo/sunodo/) version 0.8.
