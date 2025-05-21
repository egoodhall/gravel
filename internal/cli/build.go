package cli

import (
	"context"
	"encoding/json"
	"os"

	"github.com/egoodhall/gravel/internal/build"
	"github.com/egoodhall/gravel/internal/cache"
	"github.com/egoodhall/gravel/internal/gravel"
	"github.com/egoodhall/gravel/internal/resolve"
)

type buildCmd struct {
	installFlags
	buildFlags
}

func (cmd *buildCmd) Run(ctx context.Context) error {
	paths, err := gravel.NewPaths(cmd.Root)
	if err != nil {
		return err
	}

	graph, err := resolve.DependencyGraph(ctx, paths)
	if err != nil {
		return err
	}

	hashes, err := cache.NewHashes(graph, paths, !cmd.Cache)
	if err != nil {
		return err
	}

	plan, err := build.NewPlan(paths, graph, hashes, cmd.Targets)
	if err != nil {
		return err
	}

	if cmd.planOnly {
		return printJson(plan)
	}

	buildCfg := build.Config{
		Action: cmd.buildAction,
		Paths:  paths,
		Plan:   plan,
		Graph:  graph,
		Options: build.Options{
			Test: build.TestOptions{
				Enabled: cmd.Test,
			},
			Binary: build.BinaryOptions{
				Enabled: cmd.Binary,
			},
			Docker: build.DockerOptions{
				Enabled:  cmd.Docker,
				Org:      cmd.DockerOrg,
				Registry: cmd.DockerRegistry,
				Push:     cmd.DockerPush,
			},
		},
	}

	if err := build.Exec(ctx, buildCfg); err != nil {
		return err
	}

	if cmd.Cache {
		return cache.Write(paths, hashes)
	} else {
		return nil
	}
}

func printJson(obj any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(obj)
}
