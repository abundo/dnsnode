package main

import (
	cmdbase "github.com/abundo/dnsnode/cmd"
)

type TsigCmd struct {
	Name string `help:"TSIG name"`
}

func (c *TsigCmd) Run(p *cmdbase.Params) error {
	client, err := cmdbase.NewClient(*p)
	if err != nil {
		return err
	}
	data, err := client.Tsig(c.Name)
	if err != nil {
		return err
	}
	Pprint(data)
	return nil
}

type TsigCreateCmd struct {
	Name string `help:"TSIG name" required:""`
	Alg  string `help:"Algorithm" required:""`
	Key  string `help:"Key" required:""`
}

func (c *TsigCreateCmd) Run(p *cmdbase.Params) error {
	client, err := cmdbase.NewClient(*p)
	if err != nil {
		return err
	}
	return client.TsigCreate(c.Name, c.Alg, c.Key)
}

type TsigUpdateCmd struct {
	Name string `help:"TSIG name" short:"n" required:""`
	Alg  string `help:"TSIG algorithm" required:""`
	Key  string `help:"TSIG key" required:""`
}

func (c *TsigUpdateCmd) Run(p *cmdbase.Params) error {
	client, err := cmdbase.NewClient(*p)
	if err != nil {
		return err
	}
	return client.TsigUpdate(c.Name, c.Alg, c.Key)
}

type TsigDeleteCmd struct {
	Name string `help:"TSIG name" short:"n" required:""`
}

func (c *TsigDeleteCmd) Run(p *cmdbase.Params) error {
	client, err := cmdbase.NewClient(*p)
	if err != nil {
		return err
	}
	return client.TsigDelete(c.Name)
}
