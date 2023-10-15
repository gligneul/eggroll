package honeypot

import (
	"context"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gligneul/eggroll"
	"github.com/gligneul/eggroll/eggtest"
	"github.com/holiman/uint256"
)

const testTimeout = 300 * time.Second

func TestHoneypot(t *testing.T) {
	tester := eggtest.NewIntegrationTester(t)
	defer tester.Close()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	client, signer, err := eggroll.NewDevClient(ctx, Codecs())
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Send inputs
	_, err = client.SendDAppAddress(ctx, signer)
	if err != nil {
		t.Fatalf("failed to send dapp address: %v", err)
	}
	_, err = client.SendEther(ctx, signer, big.NewInt(100), nil)
	if err != nil {
		t.Fatalf("failed to send dapp ether: %v", err)
	}
	input := &Withdraw{
		Value: uint256.NewInt(50),
	}
	index, err := client.SendInput(ctx, signer, input)
	if err != nil {
		t.Fatalf("failed to send withdraw: %v", err)
	}

	// Check returned balance
	result, err := client.WaitFor(ctx, index)
	if err != nil {
		t.Fatalf("failed to wait for result: %v", err)
	}
	honeypot, ok := client.DecodeReturn(result).(*Honeypot)
	if !ok {
		t.Fatalf("failed to convert return: %v", result)
	}
	if !honeypot.Balance.Eq(uint256.NewInt(50)) {
		t.Fatal("wrong honeypot balance")
	}

	// Check voucher
	if len(result.Vouchers) != 1 {
		t.Fatal("missing voucher")
	}
	voucher := result.Vouchers[0]
	if voucher.Destination != client.DAppAddress {
		t.Fatal("wrong voucher destination")
	}
	expected := common.Hex2Bytes("522f6815000000000000000000000000f39fd6e51aad88f6f4ce6ab8827279cfffb922660000000000000000000000000000000000000000000000000000000000000032")
	if !reflect.DeepEqual(voucher.Payload, expected) {
		t.Fatalf("wrong voucher payload: %v", common.Bytes2Hex(voucher.Payload))
	}
}
