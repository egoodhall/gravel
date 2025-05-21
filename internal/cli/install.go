package cli

import (
	"context"

	"github.com/egoodhall/gravel/internal/build"
)

type installCmd struct {
	installFlags
}

func (cmd *installCmd) Run(ctx context.Context) error {
	bcmd := new(buildCmd)
	bcmd.buildAction = build.Install
	bcmd.installFlags = cmd.installFlags
	bcmd.Binary = true
	return bcmd.Run(ctx)
}
