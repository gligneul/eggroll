---
title: The Contract Interface
---

The Contract Interface
=

The box below presents the definition of the contract interface.

<!---
sed -n '92,102p' ./contract.go
-->
```go
type Contract interface {

	// Advance the contract state.
	Advance(env Env) (any, error)

	// Inspect the contract state.
	Inspect(env EnvReader) (any, error)

	// Get the codecs required by the contract.
	Codecs() []Codec
}
```

The contract interface specifies the methods required to define a contract.
The `Advance` method receives an input to advance the contract to the next state.
The `Inspect` method returns information about the contract state.
The `Codec` method specifies how to decode inputs and outputs from the contract.

# Advance

The advance method receives a value `Env` which contains the methods to interact with the Cartesi Rollups API.
For instance, the `RawInput` method returns the advance input as bytes;
the `Metadata` method returns the input metadata;
and, the `Sender` method returns the input sender.

After processing the input, the advance method can return a value or an error.
If it returns a value, it will eventually be available to the DApp client off-chain.
If the method returns an error, EggRoll logs the error and reverts the input.

The box below shows a minimal implementation for the advance method.

<!---
sed -n '11,14p' ./examples/inspect/contract/main.go
-->
```go
func (c *TemplateContract) Advance(env eggroll.Env) (any, error) {
	env.Logf("advance: %v", string(env.RawInput()))
	return env.RawInput(), nil
}
```

# Inspect

The inspect method works similarly to the advance method but receives an `EnvReader` value instead.
This interface removes from the env the methods that are not available during an inspect request.
When implementing the inspect method, you should not alter the state of the DApp contract.

The box below shows a minimal implementation for the inspect method.

<!---
sed -n '16,19p' ./examples/inspect/contract/main.go
-->
```go
func (c *TemplateContract) Inspect(env eggroll.EnvReader) (any, error) {
	env.Logf("inspect: %v", string(env.RawInput()))
	return env.RawInput(), nil
}
```

# Codecs

The codecs return the list of the codecs used by the contract.
Check out the [codecs section](/codecs) for more information about codecs.

# Default Contract

EggRoll provides a default contract implementation for the `Inspect` and `Codecs` method.
The implementation of the `Advance` method is obligatory.

<!---
sed -n '104,115p' ./contract.go
-->
```go
// DefaultContract provides a default implementation for optional contract methods.
type DefaultContract struct{}

// Reject inspect request.
func (_ DefaultContract) Inspect(env EnvReader) (any, error) {
	return nil, fmt.Errorf("inspect not supported")
}

// Return empty list of codecs.
func (_ *DefaultContract) Codecs() []Codec {
	return nil
}
```
