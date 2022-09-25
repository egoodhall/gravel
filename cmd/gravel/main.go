package main

import (
	"github.com/alecthomas/kong"
	"github.com/emm035/gravel/internal/cli"
)

var GravelCli struct {
	Build   cli.BuildCmd   `name:"build" cmd:""`
	Version cli.VersionCmd `name:"version" cmd:""`
}

func main() {
	ktx := kong.Parse(&GravelCli)
	ktx.FatalIfErrorf(ktx.Run())
}
