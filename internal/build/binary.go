package build

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sync/errgroup"
)

func execBinaryBuilds(parent context.Context, cfg Config) error {
	if !cfg.Options.Binary.Enabled || len(cfg.Plan.Build) == 0 {
		return nil
	}

	if err := os.MkdirAll(cfg.Paths.BinDir, 0755); err != nil {
		return err
	}

	env := buildEnv(cfg.Options.GOOS, cfg.Options.GOARCH)

	eg, ctx := errgroup.WithContext(parent)
	for _, tgt := range cfg.Plan.Build {
		args := generateBuildArgs(cfg, tgt)

		c := exec.CommandContext(ctx, "go", args...)
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout
		if env != nil {
			c.Env = env
		}

		fmt.Println(tgt.PkgPath)
		eg.Go(c.Run)
	}
	return eg.Wait()
}

// buildEnv returns an environment slice with GOOS and/or GOARCH set for cross-compilation.
// Returns nil if neither is set (use default process env).
func buildEnv(goos, goarch string) []string {
	if goos == "" && goarch == "" {
		return nil
	}
	env := os.Environ()
	set := func(key, val string) {
		prefix := key + "="
		for i, e := range env {
			if strings.HasPrefix(e, prefix) {
				env[i] = prefix + val
				return
			}
		}
		env = append(env, key+"="+val)
	}
	if goos != "" {
		set("GOOS", goos)
	}
	if goarch != "" {
		set("GOARCH", goarch)
	}
	return env
}
