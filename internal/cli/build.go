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
	printPlanAndExit bool
	skipTests        bool
	skipSaveCache    bool
	skipSaveVersion  bool
	buildAction      build.Action

	// Flags for setting versioning behavior
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

	hashes, err := cache.NewHashes(graph, paths, cmd.ForceBuild)
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

	if cmd.printPlanAndExit {
		return printJson(plan)
	}

	if cmd.skipTests {
		plan.Test = make([]resolve.Pkg, 0)
	}

	if err := build.Exec(semver.BumperContext(ctx, vbump), cmd.buildAction, plan); err != nil {
		return err
	}

	if cmd.Extra == "" && !cmd.skipSaveVersion {
		// Now that the build is finished, we can update any version
		// files and regenerate the hashes for the built packages.
		if err := updateVersionFiles(plan, hashes); err != nil {
			return err
		}
	}

	if cmd.skipSaveCache {
		return nil
	} else {
		return cache.Store(paths, hashes)
	}
}

func printJson(obj any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(obj)
}

func updateVersionFiles(plan build.Plan, hashes resolve.Hashes) error {
	for _, tgt := range plan.Build {
		bf, err := resolve.BuildFile(tgt.Pkg)
		if err != nil {
			// No need to update a hash if we don't
			// have a build file to write to
			continue
		}

		// Because we're writing to the version file,
		// we need to rehash the package that was built
		bf.Version = tgt.Version
		if err := bf.Save(); err != nil {
			return err
		}

		if err := hashes.New.ReHash(tgt.Pkg, tgt.Version); err != nil {
			return err
		}
	}
	return nil
}
