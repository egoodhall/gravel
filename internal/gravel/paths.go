package gravel

import (
	"path/filepath"
)

const (
	GravelDirName  = ".gravel"
	BuildFileName  = "build.json"
	HashesFileName = "cache.json"
	BinDirName     = "bin"
)

type Paths struct {
	RootDir    string `json:"rootDir"`
	GravelDir  string `json:"gravelDir"`
	BuildFile  string `json:"buildFile"`
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
		BuildFile:  filepath.Join(absRoot, GravelDirName, BuildFileName),
		HashesFile: filepath.Join(absRoot, GravelDirName, HashesFileName),
		BinDir:     filepath.Join(absRoot, GravelDirName, BinDirName),
	}, nil
}
