package cli

import "strings"

var Description = strings.Trim(`
Gravel is a build tool for small go monorepos.

It will build/test only packages that have been
changed since the last run.
`, " \n\t")

type Cli struct {
	Build   buildCmd   `name:"build" cmd:"" help:"Run a test/build cycle for any packages that have changed since the last run"`
	Version versionCmd `name:"version" cmd:"" help:"Print out the version of the gravel binary and exit"`
	Plan    planCmd    `name:"plan" cmd:""`
	Install installCmd `name:"install" cmd:""`
}
