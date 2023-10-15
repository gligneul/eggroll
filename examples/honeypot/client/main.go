package main

import (
	"context"
	"fmt"
	"honeypot"
	"log"
	"math/big"

	"github.com/gligneul/eggroll"
	"github.com/holiman/uint256"
)

func main() {
	ctx := context.Background()
	client, signer, err := eggroll.NewDevClient(ctx, honeypot.Codecs())
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.SendDAppAddress(ctx, signer)
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.SendEther(ctx, signer, big.NewInt(100), nil)
	if err != nil {
		log.Fatal(err)
	}

	input := &honeypot.Withdraw{
		Value: uint256.NewInt(50),
	}
	index, err := client.SendInput(ctx, signer, input)
	if err != nil {
		log.Fatal(err)
	}

	result, err := client.WaitFor(ctx, index)
	if err != nil {
		log.Fatal(err)
	}

	balance := client.DecodeReturn(result)
	fmt.Printf("balance: %v", balance)
}
