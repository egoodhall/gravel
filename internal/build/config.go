package build

import (
	"github.com/emm035/gravel/internal/gravel"
	"github.com/emm035/gravel/internal/resolve"
	"github.com/emm035/gravel/internal/types"
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
	Test   TestOptions
	Binary BinaryOptions
	Docker DockerOptions
}

type TestOptions struct {
	Enabled bool
}

type BinaryOptions struct {
	Enabled bool
}

type DockerOptions struct {
	Enabled  bool
	Registry string
	Org      string
}
