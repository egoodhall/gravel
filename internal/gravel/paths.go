package gravel

import (
	"path/filepath"
)

const (
	BinDirName     = "bin"
	IgnoreFileName = ".gravelignore"
	PlanFileName   = ".gravel_plan.json.gz"
)

type Paths struct {
	RootDir  string `json:"rootDir"`
	PlanFile string `json:"planFile"`
	BinDir   string `json:"binDir"`
	// SkipDir is a directory to exclude from watch (the output bin dir).
	SkipDir string `json:"skipDir"`
}

func NewPaths(root string) (Paths, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return Paths{}, err
	}

	binDir := filepath.Join(absRoot, BinDirName)
	return Paths{
		RootDir:  absRoot,
		PlanFile: filepath.Join(absRoot, PlanFileName),
		BinDir:   binDir,
		SkipDir:  binDir,
	}, nil
}
