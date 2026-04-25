package main

import (
	"fmt"
	"os"

	"github.com/user/envchain-cli/internal/chain"
	"github.com/user/envchain-cli/internal/export"
	"github.com/user/envchain-cli/internal/passphrase"
	"github.com/user/envchain-cli/internal/project"
	"github.com/urfave/cli/v2"
)

func runExport(c *cli.Context) error {
	projectName := c.String("project")
	if projectName == "" {
		var err error
		projectName, err = project.ResolveWD()
		if err != nil {
			return fmt.Errorf("export: resolve project: %w", err)
		}
	}

	pass, err := resolvePassphrase(c)
	if err != nil {
		return err
	}

	cm := chain.New(defaultChainDir())
	entries, err := cm.Load(projectName, pass)
	if err != nil {
		return fmt.Errorf("export: load chain: %w", err)
	}

	out := os.Stdout
	if outPath := c.String("output"); outPath != "" {
		f, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("export: create file: %w", err)
		}
		defer f.Close()
		out = f
	}

	ex := export.New()
	if err := ex.Write(out, projectName, entries); err != nil {
		return fmt.Errorf("export: write: %w", err)
	}
	return nil
}

func runImport(c *cli.Context) error {
	inputPath := c.Args().First()
	if inputPath == "" {
		return fmt.Errorf("import: input file path required")
	}

	f, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("import: open file: %w", err)
	}
	defer f.Close()

	ex := export.New()
	bundle, err := ex.Read(f)
	if err != nil {
		return fmt.Errorf("import: read bundle: %w", err)
	}

	projectName := c.String("project")
	if projectName == "" {
		projectName = bundle.Project
	}

	pass, err := passphrase.PromptConfirm("New passphrase for " + projectName + ": ")
	if err != nil {
		return fmt.Errorf("import: passphrase: %w", err)
	}

	cm := chain.New(defaultChainDir())
	if err := cm.Save(projectName, bundle.Entries, pass); err != nil {
		return fmt.Errorf("import: save chain: %w", err)
	}

	fmt.Fprintf(c.App.Writer, "Imported %d variable(s) into project %q.\n", len(bundle.Entries), projectName)
	return nil
}
