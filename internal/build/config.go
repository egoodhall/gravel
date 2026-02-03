package build

import (
	"github.com/egoodhall/gravel/internal/gravel"
	"github.com/egoodhall/gravel/internal/resolve"
	"github.com/egoodhall/gravel/internal/types"
)

type Config struct {
	Action  Action
	Options Options
	Paths   gravel.Paths
	Graph   types.Graph[resolve.Pkg]
	Plan    Plan
}

//go:generate go run golang.org/x/tools/cmd/stringer -type=Action -linecomment
type Action byte

const (
	Build   Action = iota // build
	Install               // install
)

type Options struct {
	Test          TestOptions
	Binary        BinaryOptions
	Version       string   // Optional version to embed in binaries
	GOOS          string   // Target OS for single cross-compile (empty = current)
	GOARCH        string   // Target architecture for single cross-compile (empty = current)
	CrossTargets  []string // Multiple platforms as "GOOS/GOARCH"; when set, build for each (single GOOS/GOARCH ignored)
	TargetSuffix  bool     // If true, always add _GOOS_GOARCH to binary name; if false, only when cross-compiling
}

type TestOptions struct {
	Enabled bool
}

type BinaryOptions struct {
	Enabled bool
}
