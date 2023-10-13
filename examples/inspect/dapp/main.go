package main

import (
	"github.com/gligneul/eggroll"
)

type TemplateContract struct {
	eggroll.DefaultContract
}

func (c *TemplateContract) Advance(env *eggroll.Env) (any, error) {
	env.Logf("advance: %v", string(env.RawInput))
	return env.RawInput, nil
}

func (c *TemplateContract) Inspect(env *eggroll.Env) (any, error) {
	env.Logf("inspect: %v", string(env.RawInput))
	return env.RawInput, nil
}

func main() {
	eggroll.Roll(&TemplateContract{})
}
