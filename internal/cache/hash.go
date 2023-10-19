package cache

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"os"

	"github.com/egoodhall/gravel/internal/gravel"
	"github.com/egoodhall/gravel/internal/resolve"
	"github.com/egoodhall/gravel/internal/types"
)

var emptyCache = resolve.CacheFile{
	Packages: nil,
}

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

func loadHashes(paths gravel.Paths, fakeLoad bool) (*resolve.CacheFile, error) {
	if fakeLoad {
		return &emptyCache, nil
	}

	file, err := os.Open(paths.HashesFile)
	if os.IsNotExist(err) {
		return &emptyCache, nil
	} else if err != nil {
		return nil, err
	}

	decfile, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}

	cache := new(resolve.CacheFile)
	if err := json.NewDecoder(decfile).Decode(cache); err != nil {
		return nil, err
	}
	return cache, nil
}

var (
	ErrPkgDirNotFound = errors.New("package directory not found")
)

func computeHashes(graph types.Graph[resolve.Pkg], paths gravel.Paths) (*resolve.CacheFile, error) {
	cacheFile := &resolve.CacheFile{
		Packages: make(map[string]string),
	}

	for pkg := range graph.Nodes() {
		hash, err := pkg.Hash()
		if err != nil {
			return nil, err
		} else if hash == "" {
			continue
		}
		cacheFile.Packages[pkg.PkgPath] = hash
	}
	return cacheFile, nil
}
