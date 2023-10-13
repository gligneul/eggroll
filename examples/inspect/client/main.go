package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gligneul/eggroll"
)

func main() {
	input := os.Args[1]
	ctx := context.Background()
	client, _ := eggroll.NewDevClient(nil)
	result, _ := client.Inspect(ctx, []byte(input))
	fmt.Println(string(result.RawReturn))
}
