package cli

import (
	"context"
	"os"
	"os/signal"

	"github.com/emm035/gravel/internal/build"
	"github.com/emm035/gravel/internal/cache"
	"github.com/emm035/gravel/internal/gravel"
	"github.com/emm035/gravel/pkg/resolve"
)

type BuildCmd struct {
	Root     string `name:"root" default:"." required:""`
	PlanOnly bool   `name:"plan" short:"p"`
}

func (cmd *BuildCmd) Run() error {
	ctx, cancel := signal.NotifyContext(context.TODO(), os.Interrupt, os.Kill)
	defer cancel()

	paths, err := gravel.NewPaths(cmd.Root)
	if err != nil {
		return err
	}

	graph, err := resolve.DependencyGraph(ctx, paths.RootDir, "./...")
	if err != nil {
		return err
	}

	hashes, err := cache.NewHashes(graph, paths)
	if err != nil {
		return err
	}

	bld, err := cache.NewBuild(paths, graph, hashes)
	if err != nil {
		return err
	}

	if !cmd.PlanOnly {
		if err := build.Exec(ctx, bld); err != nil {
			return err
		}
	}

	return cache.Store(bld, hashes, cmd.PlanOnly)
}
