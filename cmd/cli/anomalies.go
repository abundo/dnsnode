package main

import (
	cmdbase "github.com/abundo/dnsnode/cmd"
)

type AnomaliesCmd struct {
	Name string `help:"Zone name" short:"n"`
}

func (c *AnomaliesCmd) Run(p *cmdbase.Params) error {
	client, err := cmdbase.NewClient(*p)
	if err != nil {
		return err
	}
	data, err := client.Anomalies(c.Name)
	if err != nil {
		return err
	}
	Pprint(data)
	return nil
}
