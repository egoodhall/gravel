package gravel

import (
	"path/filepath"
)

const (
	GravelDirName  = "gravel"
	PlanFileName   = "plan.json.gz"
	HashesFileName = "cache.json.gz"
	BinDirName     = "bin"
	IgnoreFileName = ".gravelignore"
)

type Paths struct {
	RootDir    string `json:"rootDir"`
	GravelDir  string `json:"gravelDir"`
	PlanFile   string `json:"planFile"`
	HashesFile string `json:"hashesFile"`
	BinDir     string `json:"binDir"`
}

func NewPaths(root string) (Paths, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return Paths{}, err
	}

	return Paths{
		RootDir:    absRoot,
		GravelDir:  filepath.Join(absRoot, GravelDirName),
		PlanFile:   filepath.Join(absRoot, GravelDirName, PlanFileName),
		HashesFile: filepath.Join(absRoot, GravelDirName, HashesFileName),
		BinDir:     filepath.Join(absRoot, GravelDirName, BinDirName),
	}, nil
}
