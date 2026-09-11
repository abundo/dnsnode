// Package cmdbase provides the CLI parameters and setup shared by every
// dnsnode command-line tool, built on top of github.com/alecthomas/kong.
package cmdbase

import (
	"fmt"
	"io"
	"os"

	"github.com/abundo/dnsnode"
	"github.com/alecthomas/kong"
	kongyaml "github.com/alecthomas/kong-yaml"
	log "github.com/sirupsen/logrus"
)

// Params are the flags shared by every dnsnode CLI subcommand: which
// config file/API credentials to use, and logging verbosity. Embed this
// anonymously in the top-level CLI struct so every subcommand can access
// it via its Run(*Params) method.
type Params struct {
	ConfigFile kong.ConfigFlag `help:"Path to configuration file" short:"c" default:"/etc/dnsnode.yaml"`
	Token      string          `help:"Token for API access" env:"TOKEN"`
	URL        string          `help:"URL to dnsnode API" env:"URL"`
	Debug      bool            `help:"Enable verbose debug logging" short:"d"`
	Loglevel   string          `help:"Set log level" enum:"error,warning,info,debug" default:"info" env:"LOGLEVEL"`
}

// ConfigLoader is a kong.ConfigurationLoader for the --config-file flag.
//
// It wraps kongyaml.Loader so that an already-set environment variable
// wins over a value from the config file: kong's own resolver mechanism
// only skips a flag that was set on the command line, not one set from
// an envar during Reset(), so plain kongyaml.Loader would let the config
// file silently override an envar. Kong flag > envar > config file >
// default (mirroring the old boa-based CLI's precedence) is preserved by
// having this resolver decline to answer for any flag whose envar is set.
func ConfigLoader(r io.Reader) (kong.Resolver, error) {
	inner, err := kongyaml.Loader(r)
	if err != nil {
		return nil, err
	}
	return kong.ResolverFunc(func(ctx *kong.Context, parent *kong.Path, flag *kong.Flag) (any, error) {
		for _, env := range flag.Envs {
			if _, ok := os.LookupEnv(env); ok {
				return nil, nil
			}
		}
		return inner.Resolve(ctx, parent, flag)
	}), nil
}

// NewClient sets up logging from p and returns a DnsNode API client using
// p's token/URL, which kong has already resolved from CLI flags, env vars
// and/or the config file (in that priority order).
func NewClient(p Params) (*dnsnode.DnsNodeClient, error) {
	level := p.Loglevel
	if p.Debug {
		level = "debug"
	}
	logLevel, ok := dnsnode.Loglevels[level]
	if !ok {
		return nil, fmt.Errorf("unknown loglevel %q", level)
	}
	log.SetLevel(logLevel)

	return dnsnode.New(dnsnode.DnsNodeParam{
		Token: p.Token,
		URL:   p.URL,
		Debug: p.Debug,
	}), nil
}
