package cli

import (
	"strings"

	"github.com/alecthomas/kong"
	"github.com/egoodhall/gravel/internal/build"
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
	Watch   watchCmd   `name:"watch" cmd:"" help:"Watch the repository, rebuilding parts as files get changed."`
	Install installCmd `name:"install" cmd:"" help:"Install binaries into the $PATH."`
	Release releaseCmd `name:"release" cmd:"" help:"Release the specified binaries, using the strategies specified. Available mods are major,minor,version,date"`
}

type installFlags struct {
	ConfigFile kong.ConfigFlag `name:"config" short:"c" help:"A config file to load default flags from."`
	Root       string          `name:"root"  default:"." required:"" help:"The root directory to build. All other paths are relative to the root."`
}

type buildFlags struct {
	// Build targets filter
	Targets []string `name:"targets" arg:"" optional:"" help:"Targets to build. If none are specified, all eligible targets will be built."`

	// Build features
	Cache  bool `name:"cache" negatable:"" default:"true" help:"Use a build cache so only changed packages (and downstream dependents) are tested/built."`
	Test   bool `name:"tests" negatable:"" default:"true" help:"Run tests for changed packages during the build process."`
	Binary bool `name:"binary" negatable:"" default:"true" help:"Build a binary in the $root/gravel/bin directory"`
	Docker bool `name:"docker" negatable:"" default:"false" help:"Build a docker image containing the output binary"`

	// Docker configuration
	DockerRegistry string `name:"docker.registry" default:"" help:"The docker registry to use when building image tags."`
	DockerOrg      string `name:"docker.org" default:"" help:"The docker organization to use when building image tags."`
	DockerPush     bool   `name:"docker.push" default:"false" help:"Push images to the remote docker registry."`

	// Internal flags used by other commands
	// to configure behavior when invoking
	// the build command.
	planOnly    bool
	buildAction build.Action
}
