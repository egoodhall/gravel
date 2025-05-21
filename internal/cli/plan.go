package cli

import "context"

type planCmd struct {
	buildFlags
}

func (cmd *planCmd) Run(ctx context.Context) error {
	build := new(buildCmd)
	build.buildFlags = cmd.buildFlags
	build.planOnly = true
	return build.Run(ctx)
}
