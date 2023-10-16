---
title: Using Codecs
---

Using Codecs
=

The box below presents the codec interface.
This interface encodes Go values into bytes and decode bytes into structs.

<!---
sed -n '26,41p' ./codec.go
-->
```go
type Codec interface {

	// Get the codec key used to identify which coded the payload uses.
	Key() CodecKey

	// Get the Go type of the codec.
	Type() reflect.Type

	// Try to decode the given payload into Go value.
	// The type of the value should be the same one returned by Type().
	Decode(payload []byte) (any, error)

	// Encode a given Go value to payload.
	// The type of the value should be the same one returned by Type().
	Encode(w io.Writer, value any) error
}
```

# JSON Codecs

EggRoll provides a function to create a JSON codec for a struct type automatically.
This function relies on the JSON package from the Go standard library.

<!---
sed -n '48,49p' ./codec.go
-->
```go
// Create a new JSON codec for the struct type.
func NewJSONCodec[T any]() *JSONCodec {
```

## Example

The example below shows the JSON codec being defined for inputs and output types.
These types are stored in a package that is shared betwen the contract and the client.

<!---
sed -n '1,21p' ./examples/textbox/textbox.go
-->
```go
package textbox

import "github.com/gligneul/eggroll"

type Append struct {
	Value string
}

type Clear struct{}

type TextBox struct {
	Value string
}

func Codecs() []eggroll.Codec {
	return []eggroll.Codec{
		eggroll.NewJSONCodec[Clear](),
		eggroll.NewJSONCodec[Append](),
		eggroll.NewJSONCodec[TextBox](),
	}
}
```

### Contract

The contract should specify which codecs it uses by implementing the `Codecs` method.

<!---
sed -n '10,17p' ./examples/textbox/contract/main.go
-->
```go
type TextBoxContract struct {
	eggroll.DefaultContract
	textbox.TextBox
}

func (c *TextBoxContract) Codecs() []eggroll.Codec {
	return textbox.Codecs()
}
```

The contract can use these Go struct in the advance method.
EggRoll decodes the input in the `DecodeInput` call, and automatically encodes the return value of the advance method.

<!---
sed -n '19,31p' ./examples/textbox/contract/main.go
-->
```go
func (c *TextBoxContract) Advance(env eggroll.Env) (any, error) {
	switch input := env.DecodeInput().(type) {
	case *textbox.Clear:
		env.Log("received input clear")
		c.TextBox.Value = ""
	case *textbox.Append:
		env.Logf("received input append with '%v'\n", input.Value)
		c.TextBox.Value += input.Value
	default:
		return nil, fmt.Errorf("invalid input: %v", input)
	}
	return &c.TextBox, nil
}
```

### Client

The client should specify which codecs it uses in the constructor.

```go
	client, signer, err := eggroll.NewDevClient(ctx, textbox.Codecs())
```

The `SendInput` function encode the inputs automatically using the codecs.

```go
	inputs := []any{
		&textbox.Append{Value: "egg"},
		&textbox.Append{Value: "roll"},
	}
	for _, input := range inputs {
		lastInputIndex, err := client.SendInput(ctx, signer, input)
	}
```

The `DecodeReturn` method decodes the advance return value using the codecs.

```go
	textBox := client.DecodeReturn(result).(*textbox.TextBox)
	fmt.Println(textBox.Value) // -> eggroll
```
