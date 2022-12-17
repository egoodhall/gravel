package cli

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"

	"github.com/emm035/gravel/internal/build"
	"github.com/emm035/gravel/internal/cache"
	"github.com/emm035/gravel/internal/gravel"
	"github.com/emm035/gravel/internal/resolve"
	"github.com/emm035/gravel/internal/semver"
)

type BuildFlags struct {
	InstallFlags

	// Internal flags used by other commands
	// to configure behavior when invoking
	// the build command.
	planOnly    bool
	buildAction build.Action

	Test      bool   `name:"tests" negatable:"" default:"true"`
	BuildType string `name:"build" enum:"binary,docker" default:"binary"`

	Cache bool `name:"cache" negatable:"" default:"true"`

	// Docker configuration
	DockerRegistry string `name:"docker.registry" default:""`
	DockerOrg      string `name:"docker.org" default:""`

	// Versioning configuration
	Strategy *semver.Strategy `name:"version.strategy" xor:"type"`
	Segment  *semver.Segment  `name:"version.segment" xor:"type"`
	Extra    string           `name:"version.extra" default:""`
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

	vbump, err := semver.NewBumper(cmd.Segment, cmd.Strategy, cmd.Extra)
	if err != nil {
		return err
	}

	plan, err := build.NewPlan(paths, vbump, graph, hashes)
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
			Test:   build.TestOptions{Enabled: cmd.Test},
			Binary: build.BinaryOptions{Enabled: cmd.BuildType == "binary"},
		},
	}
	if cmd.BuildType == "docker" {
		buildCfg.Options.Docker = build.DockerOptions{
			Enabled:  true,
			Org:      cmd.DockerOrg,
			Registry: cmd.DockerRegistry,
		}
	}

	if err := build.Exec(semver.BumperContext(ctx, vbump), buildCfg); err != nil {
		return err
	}

	// Now that the build is finished, we can update any version
	// files and regenerate the hashes for the built packages.
	if err := updateVersionFiles(plan, hashes); err != nil {
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

func updateVersionFiles(plan build.Plan, hashes resolve.Hashes) error {
	for _, tgt := range plan.Build {
		vf := resolve.Version(tgt.Pkg)

		// Because we're writing to the version file,
		// we need to rehash the package that was built
		vf.Version = tgt.Version
		if err := vf.Save(); err != nil {
			return err
		}

		if err := hashes.New.ReHash(tgt.Pkg, tgt.Version); err != nil {
			return err
		}
	}
	return nil
}
