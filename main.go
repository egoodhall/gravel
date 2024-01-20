package main

import (
	"github.com/alecthomas/kong"
	"github.com/egoodhall/gravel/internal/cli"
)

func main() {
	ctx := kong.Parse(new(cli.Cli),
		kong.Description(cli.Description),
	)
	ctx.FatalIfErrorf(ctx.Run())
}
