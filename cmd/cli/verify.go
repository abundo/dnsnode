package main

import (
	cmdbase "github.com/abundo/dnsnode/cmd"
)

type VerifyCmd struct{}

func (c *VerifyCmd) Run(p *cmdbase.Params) error {
	Pprint(p)
	return nil
}
