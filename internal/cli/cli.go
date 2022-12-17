package cli

import (
	"strings"
)

var Description = strings.Trim(`
Gravel is a build tool for small go monorepos.

It will build/test only packages that have been
changed since the last run.
`, " \n\t")

type Cli struct {
	Build   buildCmd   `name:"build" cmd:"" help:"Run a test/build cycle."`
	Version versionCmd `name:"version" cmd:"" help:"Print version information for the gravel binary and exit."`
	Plan    planCmd    `name:"plan" cmd:"" help:"Print the build plan and exit."`
	Install installCmd `name:"install" cmd:"" help:"Install binaries into the $PATH."`
}
