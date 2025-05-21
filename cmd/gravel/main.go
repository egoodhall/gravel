package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/alecthomas/kong"
	"github.com/egoodhall/gravel/internal/cli"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	ktx := kong.Parse(new(cli.Cli),
		kong.Description(cli.Description),
		kong.BindTo(ctx, new(context.Context)),
	)
	ktx.FatalIfErrorf(ktx.Run(ctx))
}
