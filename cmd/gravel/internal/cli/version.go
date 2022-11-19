package cli

import (
	_ "embed"
	"fmt"

	"github.com/emm035/gravel/pkg/buildinfo"
)

type VersionCmd struct {
	Version bool `name:"version" short:"v"`
	Commit  bool `name:"commit" short:"c"`
}

func (cmd *VersionCmd) Run() error {
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

func (cmd *VersionCmd) printValue(name, value string) {
	if cmd.Version != cmd.Commit {
		fmt.Println(value)
	} else {
		fmt.Printf("%-7s = %s\n", name, value)
	}
}
