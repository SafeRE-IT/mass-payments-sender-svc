package main

import (
	"os"

	"gitlab.com/tokend/mass-payments-sender-svc/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
