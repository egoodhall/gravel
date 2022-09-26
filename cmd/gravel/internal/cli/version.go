package cli

import (
	_ "embed"
	"fmt"
)

//go:generate sh version.sh
//go:embed version.txt
var Version string

type VersionCmd struct {
}

func (cmd *VersionCmd) Run() error {
	fmt.Println(Version)
	return nil
}
