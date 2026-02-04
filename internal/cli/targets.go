package cli

import (
	"context"
	"fmt"

	"github.com/egoodhall/gravel/internal/gravel"
	"github.com/egoodhall/gravel/internal/resolve"
)

type targetsCmd struct {
	installFlags
}

func (cmd *targetsCmd) Run(ctx context.Context) error {
	paths, err := gravel.NewPaths(cmd.Root)
	if err != nil {
		return err
	}

	graph, err := resolve.DependencyGraph(ctx, paths)
	if err != nil {
		return err
	}

	for pkg := range graph.Nodes() {
		if pkg.PkgName == "main" {
			fmt.Println(pkg.PkgPath)
		}
	}

	return nil
}
