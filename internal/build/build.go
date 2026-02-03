package build

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/egoodhall/gravel/internal/resolve"
)

func Exec(ctx context.Context, cfg Config) error {
	if err := execTests(ctx, cfg); err != nil {
		return err
	}

	if len(cfg.Options.CrossTargets) > 0 {
		for _, t := range cfg.Options.CrossTargets {
			goos, goarch, err := parseCrossTarget(t)
			if err != nil {
				return err
			}
			cfg.Options.GOOS = goos
			cfg.Options.GOARCH = goarch
			if err := execBinaryBuilds(ctx, cfg); err != nil {
				return err
			}
		}
		return nil
	}

	return execBinaryBuilds(ctx, cfg)
}

// parseCrossTarget splits "GOOS/GOARCH" and returns (goos, goarch, nil) or ("", "", error).
func parseCrossTarget(s string) (goos, goarch string, err error) {
	parts := strings.Split(s, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid cross target %q: expected GOOS/GOARCH (e.g. linux/amd64)", s)
	}
	goos, goarch = strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	if goos == "" || goarch == "" {
		return "", "", fmt.Errorf("invalid cross target %q: GOOS and GOARCH must be non-empty", s)
	}
	return goos, goarch, nil
}

func generateBuildArgs(cfg Config, tgt Target) []string {
	commit, _ := resolve.GitCommit()

	version := cfg.Options.Version

	effectiveGOOS := cfg.Options.GOOS
	if effectiveGOOS == "" {
		effectiveGOOS = runtime.GOOS
	}
	effectiveGOARCH := cfg.Options.GOARCH
	if effectiveGOARCH == "" {
		effectiveGOARCH = runtime.GOARCH
	}

	binName := tgt.Binary
	crossCompiling := effectiveGOOS != runtime.GOOS || effectiveGOARCH != runtime.GOARCH
	alwaysSuffix := cfg.Options.TargetSuffix || len(cfg.Options.CrossTargets) > 0
	if alwaysSuffix || crossCompiling {
		binName = tgt.Binary + "_" + effectiveGOOS + "_" + effectiveGOARCH
	}

	args := []string{cfg.Action.String()}
	if cfg.Action != Install {
		outPath := filepath.Join(cfg.Paths.BinDir, binName)
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
