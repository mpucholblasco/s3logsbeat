package main

import (
	"os"

	"github.com/sequra/s3logsbeat/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
