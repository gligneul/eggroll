package textbox

import (
	"context"
	"testing"
	"time"

	"github.com/gligneul/eggroll"
	"github.com/gligneul/eggroll/eggeth"
	"github.com/gligneul/eggroll/eggtest"
)

const testTimeout = 300 * time.Second

func TestTextBox(t *testing.T) {
	tester := eggtest.NewIntegrationTester(t)
	defer tester.Close()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	client, signer, err := eggroll.NewDevClient(ctx, Codecs())
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	inputs := []any{
		&Append{Value: "egg"},
		&Append{Value: "roll"},
	}
	sendInputsAndVerifyState(t, ctx, client, signer, inputs, "eggroll")

	inputs = []any{
		&Clear{},
		&Append{Value: "hi"},
	}
	sendInputsAndVerifyState(t, ctx, client, signer, inputs, "hi")
}

func sendInputsAndVerifyState(
	t *testing.T, ctx context.Context, client *eggroll.Client,
	signer eggeth.Signer, inputs []any, expectedState string) {

	var lastInputIndex int
	for _, input := range inputs {
		var err error
		lastInputIndex, err = client.SendInput(ctx, signer, input)
		if err != nil {
			t.Fatalf("failed to send input: %v", err)
		}
	}

	result, err := client.WaitFor(ctx, lastInputIndex)
	if err != nil {
		t.Fatalf("failed to wait for input: %v", err)
	}

	textBox, ok := client.DecodeReturn(result).(*TextBox)
	if !ok {
		t.Fatalf("expected TextBox value")
	}
	if textBox.Value != expectedState {
		t.Fatalf("invalid value %v; expected %v", textBox.Value, expectedState)
	}
}
