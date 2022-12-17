package build

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"golang.org/x/sync/errgroup"
)

func execBinaryBuilds(parent context.Context, cfg Config) error {
	if !cfg.Options.Binary.Enabled || len(cfg.Plan.Build) == 0 {
		return nil
	}

	eg, ctx := errgroup.WithContext(parent)
	for _, tgt := range cfg.Plan.Build {
		args := generateBuildArgs(cfg.Action, cfg.Paths, tgt)

		c := exec.CommandContext(ctx, "go", args...)
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout

		fmt.Println(tgt.PkgPath)
		eg.Go(c.Run)
	}
	return eg.Wait()
}
