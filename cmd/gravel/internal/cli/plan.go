package cli

type PlanCmd struct {
	BuildFlags
}

func (cmd *PlanCmd) Run() error {
	build := new(BuildCmd)
	build.BuildFlags = cmd.BuildFlags
	build.printPlanAndExit = true
	return build.Run()
}
