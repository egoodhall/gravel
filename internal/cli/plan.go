package cli

type planCmd struct {
	buildFlags
}

func (cmd *planCmd) Run() error {
	build := new(buildCmd)
	build.buildFlags = cmd.buildFlags
	build.planOnly = true
	return build.Run()
}
