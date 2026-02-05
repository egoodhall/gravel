package cli

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/egoodhall/gravel/internal/build"
	"github.com/egoodhall/gravel/internal/gravel"
	"github.com/egoodhall/gravel/internal/resolve"
)

type watchCmd struct {
	installFlags
	buildFlags
}

func (cmd *watchCmd) Run(ctx context.Context) error {
	paths, err := gravel.NewPaths(cmd.Root)
	if err != nil {
		return err
	}

	// Build the initial graph and plan to determine which packages to watch
	graph, err := resolve.DependencyGraph(ctx, paths)
	if err != nil {
		return err
	}

	platforms, err := build.PlatformsFromOptions(cmd.Cross)
	if err != nil {
		return err
	}

	plan, err := build.NewPlan(paths, graph, cmd.Targets, platforms)
	if err != nil {
		return err
	}

	// Watch only directories for packages in the dependency graph for our targets
	// (all packages that could affect the build: test packages = build targets + their dependencies)
	events, err := build.Watch(ctx, paths, plan.Test)
	if err != nil {
		return err
	}

	fmt.Printf("Watching files in %s for changes\n", paths.RootDir)

	// Use the exact same build execution logic as buildCmd.
	runBuild := func() error {
		// Rebuild the graph and plan to pick up any changes
		graph, err := resolve.DependencyGraph(ctx, paths)
		if err != nil {
			return err
		}

		plan, err := build.NewPlan(paths, graph, cmd.Targets, platforms)
		if err != nil {
			return err
		}

		buildCfg := build.Config{
			Action: build.Build,
			Paths:  paths,
			Plan:   plan,
			Graph:  graph,
			Options: build.Options{
				Test:         build.TestOptions{Enabled: cmd.Test},
				Binary:       build.BinaryOptions{Enabled: cmd.Binary},
				Version:      cmd.Version,
				CrossTargets: cmd.Cross,
			},
		}

		return build.Exec(ctx, buildCfg)
	}

	// Run an initial build so something is built before the first file change.
	if err := runBuild(); err != nil {
		fmt.Printf("Initial build failed: %v\n", err)
	}

	for file := range events {
		if relfile, err := filepath.Rel(paths.RootDir, file); err == nil {
			file = relfile
		}

		fmt.Printf("%s changed, rebuilding.\n", file)

		if err := runBuild(); err != nil {
			fmt.Printf("Build failed: %v\n", err)
		}
	}
	return nil
}
