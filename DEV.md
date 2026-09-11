# Developer notes

## Layout

- `dnsnode.go` - the library (package `dnsnode`). One `DnsNodeClient` per
  set of credentials, created with `dnsnode.New(DnsNodeParam)`. Each API
  resource (anomalies, product, sites, statistics, status, tsig, zone, logs)
  gets its own group of methods around the shared `(*DnsNodeClient).call()`
  helper.
- `cmd/cmd_base.go` (package `cmdbase`) - CLI flags shared by every
  subcommand (`--config`, `--token`, `--url`, `--debug`, `--loglevel`) plus
  `NewClient()`, which turns a resolved `cmdbase.Params` into a
  `*dnsnode.DnsNodeClient`. Embed `cmdbase.Params` anonymously in a
  subcommand's own params struct to pick up these flags unprefixed.
- `cmd/cli/*.go` (package `main`) - one file per CLI subcommand
  (`zone.go`, `tsig.go`, `status.go`, ...), each a `kong`-tagged struct with
  a `Run(*cmdbase.Params) error` method. `dnsnode_cli.go` wires them all
  into the root `CLI` struct and is the CLI's `main()`.
- `examples/dnsnode.yaml` - template configuration file.
- `Makefile` - `make` builds `build/dnsnode_cli`; `make install` copies it
  to `/usr/bin`.

## CLI framework: kong

The CLI is built on [`github.com/alecthomas/kong`](https://github.com/alecthomas/kong),
a struct-tag-driven flag parser. Each subcommand is a plain struct (tags:
`help`, `short`, `required`, `default`, `enum`, `env`) with a
`Run(*cmdbase.Params) error` method; `dnsnode_cli.go`'s root `CLI` struct
lists them all with a `cmd:""` tag and embeds `cmdbase.Params` once for the
shared flags.

kong cleanly separates the two failure modes that boa (and cobra generally)
conflate: `parser.Parse()` only fails, and only prints usage, for genuine
CLI misuse (bad/missing flags); a subcommand's `Run()` returning an error is
just a plain Go error with no framework-printed usage attached. `main()`
handles them accordingly:

    kctx, err := parser.Parse(os.Args[1:])
    parser.FatalIfErrorf(err) // parse errors only: prints usage, exits 1
    if err := kctx.Run(&cli.Params); err != nil {
        fmt.Fprintln(os.Stderr, "Error:", err)
        os.Exit(1)
    }

No `SilenceErrors`/`IsUserInputError`-style workaround is needed (boa's
`CmdT.Run()`/`.RunE()` used to dump full usage on *any* `RunFuncE` error,
runtime failures included - that class of bug isn't representable in kong's
API).

`--config-file`'s flag/env/config-file precedence (flags win, then env vars,
then the YAML config file, then the `default:` tag) is implemented with
`kong.ConfigFlag` + `kong.Configuration(cmdbase.ConfigLoader, ...)`.
`cmdbase.ConfigLoader` wraps `kong-yaml`'s loader to skip resolving a flag
whose envar is already set - plain `kongyaml.Loader` would let the config
file silently win over an envar, since kong's resolver step only skips
flags set explicitly on the command line, not ones set from an envar during
`Reset()`.

`abmon` (`~/code/abmon`) has also been migrated from `jessevdk/go-flags` to
kong; see its `internal/check.go` (`CheckOpts`/`NewCheck`) and any of its
`cmd/check_*` binaries for the same pattern applied to Nagios/Icinga-style
one-shot checks instead of API subcommands.

## Building

    make            # build/dnsnode_cli
    make install    # install to /usr/bin
    go build ./...
    go vet ./...

No automated tests exist yet (`go test ./...` finds nothing to run).

## Known gaps / TODO

- `dnsnode_cli verify` (`cmd/cli/verify.go`) only pretty-prints the resolved
  `cmdbase.Params` - it doesn't call the API to check the token/URL are
  actually valid. There's no `Verify()` method on the library client either.
  If real verification is wanted, it likely belongs as a new library method
  (e.g. a cheap authenticated GET) plus wiring it up in `verify.go`.
- `DnsNodeConfig`/`ReadConfigFile()` in `dnsnode.go` are commented out and
  unused - kong now owns config-file loading (via `cmdbase.Params.ConfigFile`
  being a `kong.ConfigFlag`, resolved through `cmdbase.ConfigLoader`), so
  these can likely be deleted rather than resurrected.
- No test coverage.

## Pushing to GitHub

Checked and clean:

- `go build ./...`, `go vet ./...`, and `gofmt -l .` all pass with no output.
- No secrets: `examples/dnsnode.yaml` ships with an empty `token:`.
- `.gitignore` covers `build/`, `pkg/`, and common binary/test-artifact
  extensions; `build/dnsnode_cli` (a committed-looking compiled binary) is
  correctly ignored, not tracked.
- No leftover references to the old hand-rolled cobra `subcmd` package from
  before the boa migration.
- `LICENSE` (AGPL-3.0) is in place at the repo root.
