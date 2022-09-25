package cli

import (
	"fmt"

	"github.com/emm035/gravel/internal/gravel"
)

type VersionCmd struct {
}

func (cmd *VersionCmd) Run() error {
	fmt.Println(gravel.Version)
	return nil
}
