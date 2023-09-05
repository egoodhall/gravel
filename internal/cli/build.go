package cli

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"

	"github.com/egoodhall/gravel/internal/build"
	"github.com/egoodhall/gravel/internal/cache"
	"github.com/egoodhall/gravel/internal/gravel"
	"github.com/egoodhall/gravel/internal/resolve"
)

type BuildFlags struct {
	InstallFlags

	// Internal flags used by other commands
	// to configure behavior when invoking
	// the build command.
	planOnly    bool
	buildAction build.Action

	// Build targets filter
	Targets []string `name:"targets" arg:"" optional:"" help:"Targets to build. If none are specified, all eligible targets will be built."`

	// Build features
	Cache  bool `name:"cache" negatable:"" default:"true" help:"Use a build cache so only changed packages (and downstream dependents) are tested/built."`
	Test   bool `name:"tests" negatable:"" default:"true" help:"Run tests for changed packages during the build process."`
	Binary bool `name:"binary" negatable:"" default:"true" help:"Build a binary in the $root/gravel/bin directory"`
	Docker bool `name:"docker" negatable:"" default:"false" help:"Build a docker image containing the output binary"`

	// Docker configuration
	DockerRegistry string `name:"docker.registry" default:"" help:"The docker registry to use when building image tags."`
	DockerOrg      string `name:"docker.org" default:"" help:"The docker organization to use when building image tags."`
	DockerPush     bool   `name:"docker.push" default:"false" help:"Push images to the remote docker registry."`
}

type buildCmd struct {
	BuildFlags
}

func (cmd *buildCmd) Run() error {
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
		return cache.Store(paths, hashes)
	} else {
		return nil
	}
}

func printJson(obj any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(obj)
}
