package main

import (
	"context"
	"os"

	"github.com/JerrySabor/ngts-warden/internal/cli"
)

var version = "dev"

func main() {
	os.Exit(cli.New(version, os.Stdout, os.Stderr).ExecuteContext(context.Background()))
}
