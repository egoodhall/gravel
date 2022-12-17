package cli

type planCmd struct {
	BuildFlags
}

func (cmd *planCmd) Run() error {
	build := new(buildCmd)
	build.BuildFlags = cmd.BuildFlags
	build.planOnly = true
	return build.Run()
}
