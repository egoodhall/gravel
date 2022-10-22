package cache

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/emm035/gravel/internal/gravel"
	"github.com/emm035/gravel/pkg/resolve"
	"github.com/emm035/gravel/pkg/types"
)

func NewHashes(graph types.Graph[resolve.Pkg], paths gravel.Paths, ignoreOld bool) (resolve.Hashes, error) {
	oldHashes, err := loadHashes(paths, ignoreOld)
	if err != nil {
		return resolve.Hashes{}, err
	}

	newHashes, err := computeHashes(graph, paths)
	if err != nil {
		return resolve.Hashes{}, err
	}

	return resolve.Hashes{
		Old: oldHashes,
		New: newHashes,
	}, nil
}

func loadHashes(paths gravel.Paths, fakeLoad bool) (map[string]string, error) {
	if fakeLoad {
		return make(map[string]string), nil
	}

	data, err := os.ReadFile(paths.HashesFile)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	cache := make(map[string]string)
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}

	return cache, nil
}

var (
	ErrPkgDirNotFound = errors.New("package directory not found")
)

func computeHashes(graph types.Graph[resolve.Pkg], paths gravel.Paths) (map[string]string, error) {
	hashes := make(map[string]string)
	for pkg := range graph.Nodes() {
		hash, err := pkg.Hash()
		if err != nil {
			return nil, err
		}
		hashes[pkg.PkgPath] = hash
	}
	return hashes, nil
}
