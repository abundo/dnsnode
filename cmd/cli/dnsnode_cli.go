// Command dnsnode_cli is a CLI for the dnsnode library.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	cmdbase "github.com/abundo/dnsnode/cmd"

	"github.com/alecthomas/kong"
)

// Pprint prints data JSON formatted
func Pprint(data any) {
	s, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(s))
}

// CLI is the top-level dnsnode command tree. cmdbase.Params is embedded so
// its flags (--config-file, --token, ...) are available on every
// invocation, and injected into each subcommand's Run method by kong.
type CLI struct {
	cmdbase.Params

	Anomalies  AnomaliesCmd  `cmd:"" help:"Show anomalies"`
	Product    ProductCmd    `cmd:"" help:"Show products"`
	Sites      SitesCmd      `cmd:"" help:"Show sites"`
	Statistics StatisticsCmd `cmd:"" help:"Show statistics"`
	Status     StatusCmd     `cmd:"" help:"Show zone(s) status"`

	Tsig       TsigCmd       `cmd:"" name:"tsig" help:"Show TSIG(s)"`
	TsigCreate TsigCreateCmd `cmd:"" name:"tsig-create" help:"Create TSIG"`
	TsigUpdate TsigUpdateCmd `cmd:"" name:"tsig-update" help:"Update TSIG"`
	TsigDelete TsigDeleteCmd `cmd:"" name:"tsig-delete" help:"Delete TSIG"`

	Verify VerifyCmd `cmd:"" help:"Verify configuration"`

	Zone       ZoneCmd       `cmd:"" name:"zone" help:"Show zone"`
	ZoneCreate ZoneCreateCmd `cmd:"" name:"zone-create" help:"Create zone"`
	ZoneUpdate ZoneUpdateCmd `cmd:"" name:"zone-update" help:"Update zone"`
	ZoneDelete ZoneDeleteCmd `cmd:"" name:"zone-delete" help:"Delete zone"`
}

func main() {
	var cli CLI
	parser := kong.Must(&cli,
		kong.Name("dnsnode"),
		kong.Description("Manage DNSnode"),
		kong.Configuration(cmdbase.ConfigLoader, "/etc/dnsnode.yaml"),
	)

	kctx, err := parser.Parse(os.Args[1:])
	parser.FatalIfErrorf(err) // genuine CLI misuse: prints usage, exits 1

	if err := kctx.Run(&cli.Params); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
