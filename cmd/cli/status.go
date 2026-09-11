package main

import (
	"fmt"

	"github.com/abundo/dnsnode"
	cmdbase "github.com/abundo/dnsnode/cmd"
)

type StatusCmd struct {
	Name string `help:"Zone name" short:"n"`
}

func (c *StatusCmd) Run(p *cmdbase.Params) error {
	client, err := cmdbase.NewClient(*p)
	if err != nil {
		return err
	}
	if c.Name != "" {
		data, err := client.Status(c.Name)
		if err != nil {
			return err
		}
		Pprint(data)
		return nil
	}
	return PrintStatus(client, c.Name)
}

// PrintStatus prints zone(s), grouped by end customer
func PrintStatus(client *dnsnode.DnsNodeClient, name string) error {
	zones, err := client.Zone(name)
	if err != nil {
		return err
	}
	endcustomers := make(map[string][]dnsnode.ZoneType)
	for _, zone := range zones {
		endcust := zone.Endcustomer
		_, ok := endcustomers[endcust]
		if !ok {
			endcustomers[endcust] = []dnsnode.ZoneType{}
		}
		endcustomers[endcust] = append(endcustomers[endcust], zone)
	}
	for endcustomer, zones := range endcustomers {
		fmt.Printf("\n")
		fmt.Printf("----------------------------------------------------------------------\n")
		fmt.Printf("!  %s\n", endcustomer)
		fmt.Printf("----------------------------------------------------------------------\n")

		for _, zone := range zones {
			if name != "" && name != zone.Name {
				continue
			}
			fmt.Printf("\nZone %s\n", zone.Name)
			for _, master := range zone.Primaries {
				fmt.Printf("  Master IP: %s\n", master.IP)
				fmt.Printf("       TSIG: %s\n", master.TSIG)
			}
		}
	}
	return nil
}
