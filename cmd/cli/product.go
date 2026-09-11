package main

import (
	cmdbase "github.com/abundo/dnsnode/cmd"
)

type ProductCmd struct{}

func (c *ProductCmd) Run(p *cmdbase.Params) error {
	client, err := cmdbase.NewClient(*p)
	if err != nil {
		return err
	}
	data, err := client.Product()
	if err != nil {
		return err
	}
	Pprint(data)
	return nil
}
