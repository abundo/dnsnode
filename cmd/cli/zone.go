package main

import (
	"errors"
	"strings"

	"github.com/abundo/dnsnode"
	cmdbase "github.com/abundo/dnsnode/cmd"
)

// packPrimaries converts from CLI <IP address>:<TSIG> to a slice of PrimaryType
func packPrimaries(in []string) ([]dnsnode.PrimaryType, error) {
	var primaries []dnsnode.PrimaryType
	for _, p := range in {
		tmp := strings.Split(p, ":")
		if len(tmp) != 2 {
			return nil, errors.New("incorrect primaries format")
		}
		primaries = append(primaries, dnsnode.PrimaryType{IP: tmp[0], TSIG: tmp[1]})
	}
	return primaries, nil
}

type ZoneCmd struct {
	Name string `help:"Zone name" short:"n"`
}

func (c *ZoneCmd) Run(p *cmdbase.Params) error {
	client, err := cmdbase.NewClient(*p)
	if err != nil {
		return err
	}
	data, err := client.Zone(c.Name)
	if err != nil {
		return err
	}
	Pprint(data)
	return nil
}

type ZoneCreateCmd struct {
	Name        string   `help:"Zone name" short:"n"`
	Primaries   []string `help:"<IP address>:<TSIG>"`
	Product     string   `help:"Product"`
	Endcustomer string   `help:"End customer"`
}

func (c *ZoneCreateCmd) Run(p *cmdbase.Params) error {
	primaries, err := packPrimaries(c.Primaries)
	if err != nil {
		return err
	}
	client, err := cmdbase.NewClient(*p)
	if err != nil {
		return err
	}
	return client.ZoneCreate(c.Name, primaries, c.Product, c.Endcustomer)
}

type ZoneUpdateCmd struct {
	Name        string   `help:"Zone name" short:"n"`
	Primaries   []string `help:"<IP address>:<TSIG>" short:"p"`
	Product     string   `help:"Product"`
	Endcustomer string   `help:"End customer"`
}

func (c *ZoneUpdateCmd) Run(p *cmdbase.Params) error {
	primaries, err := packPrimaries(c.Primaries)
	if err != nil {
		return err
	}
	client, err := cmdbase.NewClient(*p)
	if err != nil {
		return err
	}
	return client.ZoneUpdate(c.Name, primaries, c.Product, c.Endcustomer)
}

type ZoneDeleteCmd struct {
	Name string `help:"Zone name" short:"n"`
}

func (c *ZoneDeleteCmd) Run(p *cmdbase.Params) error {
	client, err := cmdbase.NewClient(*p)
	if err != nil {
		return err
	}
	return client.ZoneDelete(c.Name)
}
