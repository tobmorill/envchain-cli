package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/envchain-cli/envchain/internal/chain"
	"github.com/envchain-cli/envchain/internal/validate"
	"github.com/urfave/cli/v2"
)

func runValidate(c *cli.Context) error {
	project := c.String("project")
	if project == "" {
		return fmt.Errorf("--project is required")
	}

	passphrase, err := resolvePassphrase(c)
	if err != nil {
		return err
	}

	mgr, err := chain.New(defaultStorePath(c))
	if err != nil {
		return fmt.Errorf("open store: %w", err)
	}

	entries, err := mgr.Load(project, passphrase)
	if err != nil {
		return fmt.Errorf("load chain: %w", err)
	}

	rule := validate.Rule{
		KeyPattern:  c.String("key-pattern"),
		ForbidEmpty: c.Bool("forbid-empty"),
	}
	if req := c.String("required"); req != "" {
		for _, k := range strings.Split(req, ",") {
			if k = strings.TrimSpace(k); k != "" {
				rule.Required = append(rule.Required, k)
			}
		}
	}

	violations := validate.Validate(entries, rule)

	if c.Bool("json") {
		type jsonViolation struct {
			Key     string `json:"key"`
			Message string `json:"message"`
		}
		out := make([]jsonViolation, len(violations))
		for i, v := range violations {
			out[i] = jsonViolation{Key: v.Key, Message: v.Message}
		}
		return json.NewEncoder(os.Stdout).Encode(out)
	}

	if len(violations) == 0 {
		fmt.Println("✓ all entries are valid")
		return nil
	}

	fmt.Fprintf(os.Stderr, "validation failed (%d violation(s)):\n", len(violations))
	for _, v := range violations {
		fmt.Fprintf(os.Stderr, "  • %s\n", v.Error())
	}
	return fmt.Errorf("validation failed")
}

// validateCommand returns the CLI sub-command definition for validate.
func validateCommand() *cli.Command {
	return &cli.Command{
		Name:  "validate",
		Usage: "validate entries in a chain against a set of rules",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "project", Aliases: []string{"p"}, Usage: "project name"},
			&cli.StringFlag{Name: "passphrase", Aliases: []string{"P"}, EnvVars: []string{"ENVCHAIN_PASSPHRASE"}},
			&cli.StringFlag{Name: "key-pattern", Usage: "regex that every key must match"},
			&cli.StringFlag{Name: "required", Usage: "comma-separated list of required keys"},
			&cli.BoolFlag{Name: "forbid-empty", Usage: "fail when any value is empty"},
			&cli.BoolFlag{Name: "json", Usage: "output violations as JSON"},
		},
		Action: runValidate,
	}
}
