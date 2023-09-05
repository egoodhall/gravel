package cli

import (
	_ "embed"
	"fmt"

	"github.com/egoodhall/gravel/pkg/buildinfo"
)

type versionCmd struct {
	Version bool `name:"version" short:"v" help:"Include the current binary's version."`
	Commit  bool `name:"commit" short:"c" help:"Include the git commit at the time the current binary was built."`
}

func (cmd *versionCmd) Run() error {
	if !cmd.Version && !cmd.Commit {
		fmt.Println(buildinfo.GetVersion())
		return nil
	}
	if cmd.Version {
		cmd.printValue("version", buildinfo.GetVersion())
	}
	if cmd.Commit {
		cmd.printValue("commit", buildinfo.GetCommit())
	}
	return nil
}

func (cmd *versionCmd) printValue(name, value string) {
	if cmd.Version != cmd.Commit {
		fmt.Println(value)
	} else {
		fmt.Printf("%-7s = %s\n", name, value)
	}
}
