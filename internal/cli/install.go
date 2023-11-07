package cli

import (
	"github.com/egoodhall/gravel/internal/build"
)

type installCmd struct {
	installFlags
}

func (cmd *installCmd) Run() error {
	bcmd := new(buildCmd)
	bcmd.buildAction = build.Install
	bcmd.installFlags = cmd.installFlags
	bcmd.Binary = true
	return bcmd.Run()
}
