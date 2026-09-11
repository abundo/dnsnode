# dnsnode

A Go library and CLI for the [Netnod Dnsnode](https://www.netnod.se/) v3 API.

The library follows the Dnsnode API's own naming conventions as closely as
possible and uses strict typing for requests/responses. The CLI
(`dnsnode_cli`) is a thin wrapper around the library, useful for scripts and
ad-hoc testing. All CLI output is JSON formatted.

External dependencies:
- `github.com/alecthomas/kong` (CLI flag/config parsing)
- `github.com/alecthomas/kong-yaml` (YAML config file support for kong)
- `github.com/sirupsen/logrus` (logging)

# Installation

Checkout the source code:

    cd /opt
    git clone https://github.com/abundo/dnsnode

If you want to use the CLI, build and install it:

    cd /opt/dnsnode
    make
    make install

This builds `build/dnsnode_cli` and installs it to `/usr/bin`.

# Configuration

The CLI needs a token to authenticate against the Dnsnode API. It can be
supplied, in order of priority: as a CLI flag, then falling back to a value
from the configuration file.

- `--config <path>` / `-c <path>` - path to the configuration file (default `/etc/dnsnode.yaml`)
- `--token <token>` - API token
- `--url <url>` - API URL (optional, defaults to `https://dnsnodeapi.netnod.se/apiv3`)
- `--debug` / `-d` - enable debug logging
- `--loglevel <level>` - one of `error`, `warning`, `info`, `debug` (default `info`)

Example configuration file (see [examples/dnsnode.yaml](examples/dnsnode.yaml)):

    ---
    token: <your API token>
    url: https://dnsnodeapi.netnod.se/apiv3/

`url` is optional; if not specified, the default API URL is used.

# ----------------------------------------------------------------------
#   CLI
# ----------------------------------------------------------------------

Run `dnsnode_cli <command> --help` for the full flag list of any command.

## ----- Anomalies ------

### Read anomalies for a zone

Example:

    $ dnsnode_cli anomalies --name example.com

    {
      "configured_sites": 11,
      "configured_distmasters": 2,
      "current_serial": 2025010714,
      "current_timestamp": 1746140501,
      "distmasters": [],
      "sites": []
    }

Here distmasters and sites are empty, i.e. there are no anomalies.


## ----- Product -----

### Read products

Example:

    $ dnsnode_cli product

    [
      "premium-a",
      "se-standard-anycast-a"
    ]


## ----- Sites -----

Parameters
  None

Example:

    $ dnsnode_cli sites

    [
      {
        "name": "DBI",
        "city": "Dubai",
        "countrycode": "AE",
        "lat": "25.26452",
        "long": "55.31167",
        "enabled": true,
        "maintenance": false
      },
      <snip>
    ]


## ----- Statistics -----

### Read statistics

Parameters
  - `--name <zonename>`

Example:

    $ dnsnode_cli statistics --name example.com

    {
      "timestamps": [
        1746120300,
        1746120600,
        1746141600
      ],
      "values": []
    }


## ----- Status -----

### Read status for all zones, grouped by end customer

    $ dnsnode_cli status

Example:

    $ dnsnode_cli status

    ----------------------------------------------------------------------
    !  Company
    ----------------------------------------------------------------------

    Zone example.com
      Master IP: 192.168.1.1
           TSIG: netnod-example1.


### Read status for one zone

    $ dnsnode_cli status --name <zone-name>

Example:

    $ dnsnode_cli status --name example.com

    {
      "configured_sites": 11,
      "configured_distmasters": 2,
      "current_serial": 2025010715,
      "current_timestamp": 1746143648,
      "distmasters": [
        {
          "ipv4_address": "194.146.105.24",
          "ipv6_address": "2a01:3f0:0:27::24",
          "serial": 2025010715,
          "timestamp": 1746143648
        },
        <snip>
      ],
      "sites": [
        {
          "name": "BNX",
          "serial": 2025010715,
          "timestamp": 1746143669
        },
        <snip>
      ]
    }


## ----- TSIG -----

### Read all TSIG

    $ dnsnode_cli tsig
    [
      {
        "key": "VXdkbmp1elR3MHQxTkd1Ynlmayt3VTh1TVVtajFQTVB4SFNjK1A2ZFJlQT0=",
        "name": "netnod-abundo-dist.",
        "alg": "hmac-sha256"
      },
      <snip>
    ]


### Read one TSIG

    $ dnsnode_cli tsig --name netnod-example.
    [
      {
        "key": "<hidden>",
        "name": "netnod-example.",
        "alg": "hmac-sha256"
      }
    ]


### Create TSIG

TSIG names must start with `netnod-`.

    $ dnsnode_cli tsig-create --name <tsig-name> --alg <algorithm> --key <key>

Example:

    $ dnsnode_cli tsig-create --name netnod-example. --alg hmac-sha256 --key <base64 key>

Supported algorithms: `hmac-md5`, `hmac-sha1`, `hmac-sha256`, `hmac-sha512`.


### Update TSIG

`--alg` and `--key` are optional; an empty value leaves that field unchanged.

    $ dnsnode_cli tsig-update --name <tsig-name> --alg <algorithm> --key <key>

Example:

    $ dnsnode_cli tsig-update --name netnod-example. --alg hmac-sha256 --key <base64 key>


### Delete TSIG

    $ dnsnode_cli tsig-delete --name <tsig-name>

Example:

    $ dnsnode_cli tsig-delete --name netnod-se-tsigtest.


## ----- Zone ------

Parameters (create/update):
  - `--masters <ip address>:<tsig>` - can be repeated for multiple masters. TSIG must already exist.
  - `--product <product name>` - one of the names returned by the `product` command.
  - `--endcustomer <string>` - customer identifier.

### List zones

    $ dnsnode_cli zone

Example: array of zones, see "Read one zone" below for the zone shape.


### Read one zone

    $ dnsnode_cli zone --name <zone>

Example:

    $ dnsnode_cli zone --name example.com

    [
      {
        "name": "example.com",
        "masters": [
          { "ip": "192.168.1.1", "tsig": "netnod-example1." }
        ],
        "product": "se-standard-anycast-a",
        "endcustomer": "Company"
      }
    ]


### Create zone

    $ dnsnode_cli zone-create --name <zone-name> --product <product> --endcustomer <name> --masters <ip>:<tsig>

Example:

    $ dnsnode_cli zone-create --name example.com --product se-standard-anycast-a --endcustomer Company --masters 192.168.1.1:netnod-example1.


### Update zone

`--masters`, `--product` and `--endcustomer` are optional; an empty/unset value leaves that field unchanged.

    $ dnsnode_cli zone-update --name <zone-name> --product <product> --endcustomer <name> --masters <ip>:<tsig>

Example:

    $ dnsnode_cli zone-update --name example.com --product se-standard-anycast-a --endcustomer Company --masters 10.0.0.1:netnod-example1.


### Delete zone

    $ dnsnode_cli zone-delete --name <zone-name>

Example:

    $ dnsnode_cli zone-delete --name example.com


# ----------------------------------------------------------------------
#   Library
# ----------------------------------------------------------------------

`cmd/cli/*.go` is a good, complete example of how to use the library.

## Creating a client

    import (
        dnsnode "github.com/abundo/dnsnode"
    )

    client := dnsnode.New(dnsnode.DnsNodeParam{
        Token: "<token>",
        URL:   "",    // optional, defaults to the production API
        Debug: false,
    })

All calls below return `(result, error)` and take no config/context beyond
the client itself.

## Anomalies

    anomalies, err := client.Anomalies("example.com")
    zones, err := client.AnomaliesList() // zones with a serial mismatch

## Product

    products, err := client.Product()

## Sites

    sites, err := client.Sites()       // sites configured for this customer
    all, err := client.SitesAll()      // all sites Netnod offers

## Statistics

    stats, err := client.Statistics("example.com")

## Status

    status, err := client.Status("example.com")

## TSIG

    tsigs, err := client.Tsig("")                         // all TSIGs
    tsig, err := client.Tsig("netnod-example.")           // one TSIG
    err := client.TsigCreate("netnod-example.", "hmac-sha256", key)
    err := client.TsigUpdate("netnod-example.", "hmac-sha256", key) // alg/key optional
    err := client.TsigDelete("netnod-example.")

## Zone

    zones, err := client.Zone("")                 // all zones
    zone, err := client.Zone("example.com")          // one zone
    err := client.ZoneCreate("example.com", primaries, "se-standard-anycast-a", "Company")
    err := client.ZoneUpdate("example.com", primaries, "se-standard-anycast-a", "Company") // fields optional
    err := client.ZoneDelete("example.com")

Where `primaries` is `[]dnsnode.PrimaryType{{IP: "192.168.1.1", TSIG: "netnod-example1."}}`.

## Logs

    logs, err := client.LogsXfer("example.com") // transfer/notify logs for a zone

See [DEV.md](DEV.md) for build/development notes and known limitations.
