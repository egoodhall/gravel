package main

import (
	"strings"

	"github.com/alecthomas/kong"
	"github.com/emm035/gravel/cmd/gravel/internal/cli"
)

var description = strings.Trim(`
Gravel monorepo build tool
`, " \n\t")

var GravelCli struct {
	Build   cli.BuildCmd   `name:"build" cmd:"" help:"Run a test/build cycle for any packages that have changed since the last run"`
	Version cli.VersionCmd `name:"version" cmd:"" help:"Print out the version of the gravel binary and exit"`
}

func main() {
	ktx := kong.Parse(&GravelCli,
		kong.Description(description),
	)
	ktx.FatalIfErrorf(ktx.Run())
}
