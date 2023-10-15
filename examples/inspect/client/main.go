package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gligneul/eggroll"
)

func main() {
	input := os.Args[1]
	ctx := context.Background()
	client, _, err := eggroll.NewDevClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	result, err := client.Inspect(ctx, []byte(input))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(result.RawReturn()))
}
