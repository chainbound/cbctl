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
					if key == "" {
						return fmt.Errorf("API key not provided")
					}

					url := c.String("url")
					if url == "" {
						return fmt.Errorf("URL not provided")
					}

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
