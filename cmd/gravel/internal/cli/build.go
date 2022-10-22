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
	Root       string `name:"root" default:"." required:"" help:"The root directory to build. All other paths are relative to the root"`
	ForceBuild bool   `name:"force" short:"f" help:"Force all packages to be built/tested, regardless of whether they have changed"`
	PlanOnly   bool   `name:"plan-only" short:"p" help:"Generate a build/test plan (without running it), and write it to gravel/plan.json"`
}

func (cmd *BuildCmd) Run() error {
	ctx, cancel := signal.NotifyContext(context.TODO(), os.Interrupt, os.Kill)
	defer cancel()

	paths, err := gravel.NewPaths(cmd.Root)
	if err != nil {
		return err
	}

	graph, err := resolve.DependencyGraph(ctx, paths)
	if err != nil {
		return err
	}

	hashes, err := cache.NewHashes(graph, paths, cmd.ForceBuild)
	if err != nil {
		return err
	}

	bld, err := build.NewPlan(paths, graph, hashes)
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
