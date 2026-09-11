package main

import (
	cmdbase "github.com/abundo/dnsnode/cmd"
)

type StatisticsCmd struct {
	Name string `help:"Zone name" short:"n"`
}

func (c *StatisticsCmd) Run(p *cmdbase.Params) error {
	client, err := cmdbase.NewClient(*p)
	if err != nil {
		return err
	}
	data, err := client.Statistics(c.Name)
	if err != nil {
		return err
	}
	Pprint(data)
	return nil
}
