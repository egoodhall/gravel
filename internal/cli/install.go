package cli

import (
	"github.com/alecthomas/kong"
	"github.com/emm035/gravel/internal/build"
)

type InstallFlags struct {
	ConfigFile kong.ConfigFlag `name:"config" short:"c"`
	Root       string          `name:"root" arg:""  default:"." required:"" help:"The root directory to build. All other paths are relative to the root"`
	ForceBuild bool            `name:"force" short:"f" help:"Force all packages to be built/tested, regardless of whether they have changed"`
}

type installCmd struct {
	InstallFlags
}

func (cmd *installCmd) Run() error {
	bcmd := new(buildCmd)
	bcmd.ForceBuild = true
	bcmd.buildAction = build.Install
	bcmd.InstallFlags = cmd.InstallFlags
	bcmd.BuildType = "binary"
	return bcmd.Run()
}
