package main

import (
	"os"

	"ordo/internal/cli"
)

func main() {
	os.Exit(cli.Execute())
}
