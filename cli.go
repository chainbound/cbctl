package main

import (
	"fmt"
	"os"

	"github.com/chainbound/cbctl/api"

	"github.com/urfave/cli/v2"
)

func NewApp() *cli.App {
	app := &cli.App{
		Name:  "cbctl",
		Usage: "A command line tool for interacting with the Chainbound API",
		Commands: []*cli.Command{
			{
				Name:  "init",
				Usage: "Initialize cbctl",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "url",
						Usage: "The URL of the Chainbound API",
						Value: "http://backend.fiberapi.io:8000",
					},
					&cli.StringFlag{
						Name:     "key",
						Usage:    "Your API key",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					key := c.String("key")
					url := c.String("url")

					path, err := inititalize(url, key)
					if err != nil {
						return err
					}

					fmt.Printf("Initialized cbctl configuration at %s\n", path)

					return nil
				},
			},
			{
				Name:  "fiber",
				Usage: "Interact with the Fiber API",
				Subcommands: []*cli.Command{
					{
						Name:  "quota",
						Usage: "Get your current usage quota",
						Action: func(c *cli.Context) error {
							cfg, err := ReadConfig()
							if err != nil {
								return fmt.Errorf("Error reading config, did you run cbctl init?: %w", err)
							}

							api := api.NewFiberAPI(cfg.Url, cfg.ApiKey)
							quota, err := api.GetQuota()
							if err != nil {
								return fmt.Errorf("Error getting quota: %w", err)
							}

							println("Current billing period's quota:")
							println(fmt.Sprintf("  EgressMB: %d", quota.EgressMB))
							println(fmt.Sprintf("  MaxEgressMB: %d", quota.MaxEgressMB))
							println(fmt.Sprintf("  TransactionCount: %d", quota.TransactionCount))
							println(fmt.Sprintf("  ActiveStreams: %d", quota.ActiveStreams))
							println(fmt.Sprintf("  MaxActiveStreams: %d", quota.MaxActiveStreams))

							return nil
						},
					},
					{
						Name:  "trace",
						Usage: "Trace transactions or blocks",
						Subcommands: []*cli.Command{
							{
								Name:  "tx",
								Usage: "Trace a transaction",
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:     "hash",
										Aliases:  []string{"H"},
										Usage:    "The transaction hash to trace",
										Required: true,
									},
									&cli.StringFlag{
										Name:    "type",
										Aliases: []string{"t"},
										Usage:   "The observation type to trace (p2p | fiber | all)",
										Value:   "all",
									},
									&cli.BoolFlag{
										Name:    "private",
										Aliases: []string{"p"},
										Usage:   "Whether or not the transaction is private (sent with your API key)",
										Value:   false,
									},
									&cli.BoolFlag{
										Name:    "show-source",
										Aliases: []string{"s"},
										Usage:   "Whether or not to show the source of the transaction",
										Value:   false,
									},
								},
								Action: func(c *cli.Context) error {
									cfg, err := ReadConfig()
									if err != nil {
										return fmt.Errorf("Error reading config, did you run cbctl init?: %w", err)
									}

									hash := c.String("hash")
									observationType := c.String("type")
									private := c.Bool("private")
									showSource := c.Bool("show-source")

									api := api.NewFiberAPI(cfg.Url, cfg.ApiKey)
									traces, err := api.TraceTransaction(hash, observationType, private)
									if err != nil {
										return fmt.Errorf("Error getting quota: %w", err)
									}

									printMessageTrace(traces, showSource)

									return nil
								},
							},
							{
								Name:  "block",
								Usage: "Trace a block",
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:    "hash",
										Aliases: []string{"H"},
										Usage:   "The block hash to trace",
									},
									&cli.StringFlag{
										Name:    "number",
										Aliases: []string{"n"},
										Usage:   "The block number to trace",
									},
									&cli.StringFlag{
										Name:    "type",
										Aliases: []string{"t"},
										Usage:   "The observation type to trace (p2p | fiber | all)",
										Value:   "all",
									},
									&cli.BoolFlag{
										Name:    "show-source",
										Aliases: []string{"s"},
										Usage:   "Whether or not to show the source of the transaction",
										Value:   false,
									},
								},
								Action: func(c *cli.Context) error {
									cfg, err := ReadConfig()
									if err != nil {
										return fmt.Errorf("Error reading config, did you run cbctl init?: %w", err)
									}

									hash := c.String("hash")
									number := c.String("number")
									if hash == "" && number == "" {
										return fmt.Errorf("Must specify either a block hash or block number")
									}

									hashOrNumber := number
									if hashOrNumber == "" {
										hashOrNumber = hash
									}

									observationType := c.String("type")
									showSource := c.Bool("show-source")

									api := api.NewFiberAPI(cfg.Url, cfg.ApiKey)
									traces, err := api.TraceBlock(hashOrNumber, observationType)
									if err != nil {
										return fmt.Errorf("Error getting quota: %w", err)
									}

									printMessageTrace(traces, showSource)

									return nil
								},
							},
							{
								Name:  "blob",
								Usage: "Trace a blob sidecar (consensus layer)",
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:     "commitment",
										Aliases:  []string{"C"},
										Usage:    "The blob KZG commitment to trace",
										Required: true,
									},
									&cli.StringFlag{
										Name:    "type",
										Aliases: []string{"t"},
										Usage:   "The observation type to trace (p2p | fiber | all)",
										Value:   "all",
									},
									&cli.BoolFlag{
										Name:    "show-source",
										Aliases: []string{"s"},
										Usage:   "Whether or not to show the source of the transaction",
										Value:   false,
									},
								},
								Action: func(c *cli.Context) error {
									cfg, err := ReadConfig()
									if err != nil {
										return fmt.Errorf("Error reading config, did you run cbctl init?: %w", err)
									}

									commitment := c.String("commitment")
									if commitment == "" {
										return fmt.Errorf("Must specify a blob commitment")
									}

									observationType := c.String("type")
									showSource := c.Bool("show-source")

									api := api.NewFiberAPI(cfg.Url, cfg.ApiKey)
									traces, err := api.TraceBlob(commitment, observationType)
									if err != nil {
										return fmt.Errorf("Error getting quota: %w", err)
									}

									printMessageTrace(traces, showSource)

									return nil
								},
							},
						},
					},
				},
			},
		},
	}

	return app
}

func inititalize(url, apiKey string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(home+"/.config/cbctl", 0755); err != nil {
		return "", err
	}

	configBytes := []byte(fmt.Sprintf("url = \"%s\"\napi_key = \"%s\"\n", url, apiKey))

	path := home + "/.config/cbctl/config.toml"

	if err := os.WriteFile(path, configBytes, 0644); err != nil {
		return "", err
	}

	return path, nil
}
