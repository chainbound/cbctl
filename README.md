# `cbctl`

`cbctl` is a CLI tool for interacting with the Chainbound backend.

## Installation
With Go:
```bash
go get github.com/chainbound/cbctl
```

## Usage
First, initialize `cbctl` with your API key. This will write a config file to `~/.config/cbctl/config.toml`:
```bash
cbctl init --key <your_api_key>
```

### Fiber Commands
The `fiber` subcommand let's you interact with the Fiber backend.

* `cbctl fiber quota`
Gets your quota for the current billing period.
* `cbctl fiber trace tx --hash <tx_hash>`
Gets the trace for a transaction. Run `cbctl fiber trace tx --help` for more info.
