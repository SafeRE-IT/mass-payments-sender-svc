package main

import (
	"os"

	"github.com/SafeRE-IT/mass-payments-sender-svc/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
