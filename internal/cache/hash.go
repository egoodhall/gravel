package cache

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/emm035/gravel/internal/gravel"
	"github.com/emm035/gravel/pkg/resolve"
	"github.com/emm035/gravel/pkg/types"
)

type File struct {
	Hashes map[string]string `json:"hashes"`
}

type Hashes struct {
	Old map[string]string
	New map[string]string
}

func NewHashes(graph types.Graph[resolve.Pkg], paths gravel.Paths) (Hashes, error) {
	oldHashes, err := loadHashes(paths)
	if err != nil {
		return Hashes{}, err
	}

	newHashes, err := computeHashes(graph, paths)
	if err != nil {
		return Hashes{}, err
	}

	return Hashes{
		Old: oldHashes,
		New: newHashes,
	}, nil
}

func (h Hashes) Changed() types.Set[string] {
	s := types.NewSet[string]()
	for pkg, hash := range h.New {
		if hash != h.Old[pkg] {
			s.Add(pkg)
		}
	}
	for pkg, hash := range h.Old {
		if hash != h.New[pkg] {
			s.Add(pkg)
		}
	}
	return s
}

func loadHashes(paths gravel.Paths) (map[string]string, error) {
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
