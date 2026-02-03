package build

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/egoodhall/gravel/internal/resolve"
)

func Exec(ctx context.Context, cfg Config) error {
	if err := execTests(ctx, cfg); err != nil {
		return err
	}

	if err := execBinaryBuilds(ctx, cfg); err != nil {
		return err
	}

	return nil
}

func generateBuildArgs(cfg Config, tgt Target) []string {
	commit, _ := resolve.GitCommit()

	version := cfg.Options.Version

	args := []string{cfg.Action.String()}
	if cfg.Action != Install {
		outPath := filepath.Join(cfg.Paths.BinDir, tgt.Binary)
		if p, err := filepath.Rel(tgt.Module.DirPath, outPath); err == nil {
			outPath = "./" + p
		}
		args = append(args, "-o", outPath)
	}

	args = append(args,
		"-ldflags", buildLdFlags(version, commit),
		tgt.PkgPath,
	)
	return args
}

func buildLdFlags(version, commit string) string {
	return strings.Join([]string{
		"-s -w",
		fmt.Sprintf("-X github.com/egoodhall/gravel/pkg/buildinfo.version=%s", version),
		fmt.Sprintf("-X github.com/egoodhall/gravel/pkg/buildinfo.commit=%s", commit),
	}, " ")
}
