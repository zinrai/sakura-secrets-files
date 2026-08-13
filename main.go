package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "present":
		err = runPresent(os.Args[2:])
	case "absent":
		err = runAbsent(os.Args[2:])
	case "version":
		runVersion()
	default:
		fmt.Fprintf(os.Stderr, "Unknown subcommand: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Usage: sakura-secrets-files <subcommand> [options]

Subcommands:
  present  Write the declared secrets to their destination paths
  absent   Remove the declared destination paths
  version  Print version

Use "sakura-secrets-files <subcommand> -h" for more information about a subcommand.
`)
}
