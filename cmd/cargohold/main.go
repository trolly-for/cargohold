package main

import (
	"fmt"
	"os"

	"cargohold/internal/cli"
	"cargohold/internal/output"
	"cargohold/internal/passphrase"
	"cargohold/internal/store"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	st, err := store.Default()
	if err != nil {
		return fmt.Errorf("initialising store: %w", err)
	}

	out := output.Default()

	readPass := func(confirm bool) (string, error) {
		if confirm {
			return passphrase.ReadWithConfirm()
		}
		return passphrase.Read()
	}

	runner := cli.New(st, out, readPass)

	args := os.Args[1:]
	if len(args) == 0 {
		printUsage(out)
		return nil
	}

	switch args[0] {
	case "init":
		if len(args) < 2 {
			return fmt.Errorf("usage: cargohold init <env>")
		}
		return runner.Init(args[1])

	case "set":
		if len(args) < 4 {
			return fmt.Errorf("usage: cargohold set <env> <key> <value>")
		}
		return runner.Set(args[1], args[2], args[3])

	case "get":
		if len(args) < 3 {
			return fmt.Errorf("usage: cargohold get <env> <key>")
		}
		return runner.Get(args[1], args[2])

	case "delete":
		if len(args) < 3 {
			return fmt.Errorf("usage: cargohold delete <env> <key>")
		}
		return runner.Delete(args[1], args[2])

	case "list":
		if len(args) < 2 {
			return fmt.Errorf("usage: cargohold list <env>")
		}
		return runner.List(args[1])

	case "help", "--help", "-h":
		printUsage(out)
		return nil

	default:
		return fmt.Errorf("unknown command %q — run 'cargohold help' for usage", args[0])
	}
}

func printUsage(out *output.Formatter) {
	out.Success("cargohold — environment-specific secret bundle manager")
	out.KeyValue("init <env>", "create a new encrypted bundle for an environment")
	out.KeyValue("set <env> <key> <value>", "store a secret in a bundle")
	out.KeyValue("get <env> <key>", "retrieve a secret from a bundle")
	out.KeyValue("delete <env> <key>", "remove a secret from a bundle")
	out.KeyValue("list <env>", "list all keys in a bundle")
}
