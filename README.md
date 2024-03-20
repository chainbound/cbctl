# `cbctl`

`cbctl` is a CLI tool for interacting with the Chainbound backend.

## Installation

With Go:

```bash
go install github.com/chainbound/cbctl@latest
```

## Usage

First, initialize `cbctl` with your API key. This will write a config file to `~/.config/cbctl/config.toml`:

```bash
cbctl init --key <your_api_key>
```

### Fiber Commands

The `fiber` subcommand let's you interact with the Fiber backend.

- `cbctl fiber quota`

Gets your quota for the current billing period.

- `cbctl fiber trace tx --hash <tx_hash>`

Gets the trace for a transaction. Run `cbctl fiber trace tx --help` for more info.

- `cbctl fiber trace block --hash <block_hash>` or `cbctl fiber trace block --number <block_number>`

Gets a block trace. Run `cbctl fiber trace block --help` for more info.

- `cbctl fiber trace blob --commitment <blob_kzg_commitment>`

Gets a blob trace. Run `cbctl fiber trace blob --help` for more info.
