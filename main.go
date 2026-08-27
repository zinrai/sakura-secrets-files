package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Usage = printUsage

	configPath := flag.String("config", "", "Path to the manifest (required)")
	baseDir := flag.String("base-dir", "", "Directory the dest paths are resolved against (required)")
	showVersion := flag.Bool("version", false, "Print version and exit")

	flag.Parse()

	if *showVersion {
		printVersion()
		return
	}

	if *configPath == "" {
		flag.Usage()
		fail(fmt.Errorf("-config is required"))
	}
	// No default and no cwd fallback: without a directory chosen for this run,
	// a mistaken invocation writes plaintext wherever it happens to be standing
	if *baseDir == "" {
		flag.Usage()
		fail(fmt.Errorf("-base-dir is required"))
	}

	if err := run(*configPath, *baseDir); err != nil {
		fail(err)
	}
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Usage: sakura-secrets-files -config <path> -base-dir <dir>

Writes the secrets declared in the manifest to their dest paths under base-dir.
The lifetime of the written files is the caller's responsibility.

Options:
`)
	flag.PrintDefaults()
}
