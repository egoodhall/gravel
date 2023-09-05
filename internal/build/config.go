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
	Push     bool
	Registry string
	Org      string
}
