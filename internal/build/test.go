package build

import (
	"context"
	"os"
	"os/exec"
)

func execTests(ctx context.Context, cfg Config) error {
	if !cfg.Options.Test.Enabled || len(cfg.Plan.Test) == 0 {
		return nil
	}

	args := []string{"test", "-count", "1"}
	for _, pkg := range cfg.Plan.Test {
		args = append(args, pkg.PkgPath)
	}

	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
