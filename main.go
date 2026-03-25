package main

import (
	_ "embed"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops. An error while executing command '%s'\n", err)
		os.Exit(1)
	}
}
