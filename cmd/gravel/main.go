package main

import (
	"github.com/alecthomas/kong"
	"github.com/emm035/gravel/internal/cli"
)

func main() {
	ktx := kong.Parse(new(cli.Cli),
		kong.Description(cli.Description),
	)
	ktx.FatalIfErrorf(ktx.Run())
}
