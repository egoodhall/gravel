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

	Test      bool   `name:"tests" negatable:"" default:"true" help:"Run tests for changed packages during the build process."`
	BuildType string `name:"build" enum:"binary,docker" default:"binary" help:"The type of build to run. Must be one of: 'binary','docker'."`

	Cache bool `name:"cache" negatable:"" default:"true" help:"Use a build cache so only changed packages (and downstream dependents) are tested/built."`

	// Docker configuration
	DockerRegistry string `name:"docker.registry" default:"" help:"The docker registry to use when building image tags."`
	DockerOrg      string `name:"docker.org" default:"" help:"The docker organization to use when building image tags."`
	DockerPush     bool   `name:"docker.push" default:"false" help:"Push images to the remote docker registry."`

	// Versioning configuration
	Strategy *semver.Strategy `name:"version.strategy" xor:"type" help:"A version update strategy to use when a build occurs."`
	Segment  *semver.Segment  `name:"version.segment" xor:"type" help:"A semantic version segment to update."`
	Extra    string           `name:"version.extra" default:"" help:"An extra string that will be added to the updated version. For a value of 'test', the new version would be: 'v##.##.##-test'."`
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

	if os.Getenv("CI") != "" && cmd.buildAction == build.Build {
		// Now that the build is finished, we can update any version
		// files and regenerate the hashes for the built packages.
		if err := updateVersionFiles(plan, hashes); err != nil {
			return err
		}
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
