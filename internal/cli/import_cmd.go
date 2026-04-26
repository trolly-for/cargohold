package cli

import (
	"fmt"
	"os"

	importer "cargohold/internal/import"
	"cargohold/internal/passphrase"
)

// runImport handles the `cargohold import <env> <file> [--format dotenv|json]` command.
func (r *Runner) runImport(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: cargohold import <env> <file> [--format dotenv|json]")
	}

	env := args[0]
	filePath := args[1]

	formatStr := "dotenv"
	for i, a := range args[2:] {
		if a == "--format" && i+1 < len(args[2:]) {
			formatStr = args[2:][i+1]
		}
	}

	fmt_, err := importer.ParseFormat(formatStr)
	if err != nil {
		return err
	}

	pass, err := passphrase.Read(env)
	if err != nil {
		return fmt.Errorf("read passphrase: %w", err)
	}

	v, err := r.vault.Open(env, pass)
	if err != nil {
		return fmt.Errorf("open vault: %w", err)
	}

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open file %q: %w", filePath, err)
	}
	defer f.Close()

	n, err := importer.Import(v, f, fmt_)
	if err != nil {
		return fmt.Errorf("import: %w", err)
	}

	if err := r.vault.Save(env, pass, v); err != nil {
		return fmt.Errorf("save vault: %w", err)
	}

	r.out.Success(fmt.Sprintf("imported %d key(s) into %q from %s", n, env, filePath))
	return nil
}
