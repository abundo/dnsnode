package main

import (
	cmdbase "github.com/abundo/dnsnode/cmd"
)

type SitesCmd struct{}

func (c *SitesCmd) Run(p *cmdbase.Params) error {
	client, err := cmdbase.NewClient(*p)
	if err != nil {
		return err
	}
	data, err := client.Sites()
	if err != nil {
		return err
	}
	Pprint(data)
	return nil
}
