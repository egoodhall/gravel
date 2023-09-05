package cli

import (
	"github.com/alecthomas/kong"
	"github.com/egoodhall/gravel/internal/build"
)

type InstallFlags struct {
	ConfigFile kong.ConfigFlag `name:"config" short:"c" help:"A config file to load default flags from."`
	Root       string          `name:"root"  default:"." required:"" help:"The root directory to build. All other paths are relative to the root."`
}

type installCmd struct {
	InstallFlags
}

func (cmd *installCmd) Run() error {
	bcmd := new(buildCmd)
	bcmd.buildAction = build.Install
	bcmd.InstallFlags = cmd.InstallFlags
	bcmd.Binary = true
	return bcmd.Run()
}
