package build

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/emm035/gravel/internal/gravel"
	"github.com/emm035/gravel/internal/resolve"
	"golang.org/x/sync/errgroup"
)

var funcs = template.FuncMap{
	"rel": filepath.Rel,
}

var dockerfile = template.Must(template.New("Dockerfile").Funcs(funcs).Parse(`
FROM golang:1.19 as build

WORKDIR /

# Copy Dependency packages
{{- range .Dependencies}}
COPY ./{{ rel $.Paths.RootDir .DirPath }} /{{ rel $.Paths.RootDir .DirPath }}
{{- end}}

# Copy build files
COPY {{ rel $.Paths.RootDir .Package.Module.DirPath }}/go.mod /go.mod
COPY {{ rel $.Paths.RootDir .Package.Module.DirPath }}/go.sum /go.sum

# Copy {{ .Package.Binary }} main package
COPY ./{{ rel $.Paths.RootDir .Package.DirPath }} /{{ rel $.Paths.RootDir .Package.DirPath }}

RUN {{ .BuildCommand }}

FROM gcr.io/distroless/base:debug

WORKDIR /

ENTRYPOINT ["{{.Package.Binary}}"]

COPY --from=build /gravel/bin/{{ .Package.Binary }} /bin/{{ .Package.Binary }}

`))

type dockerContext struct {
	Paths        gravel.Paths
	Package      resolve.Pkg
	Dependencies []resolve.Pkg
	BuildCommand string
}

func execDockerBuilds(ctx context.Context, cfg Config) error {
	if !cfg.Options.Docker.Enabled || len(cfg.Plan.Build) == 0 {
		return nil
	}

	egrp, gctx := errgroup.WithContext(ctx)
	for _, tgt := range cfg.Plan.Build {
		egrp.Go(execDockerBuild(gctx, cfg, tgt))
	}
	return egrp.Wait()
}

func execDockerBuild(ctx context.Context, cfg Config, tgt Target) func() error {
	tag := buildDockerTag(cfg.Options.Docker, tgt)
	cmd := exec.Command("docker", "build", "--tag", tag, "-f", "-", cfg.Paths.RootDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return func() error { return err }
	}

	if err := cmd.Start(); err != nil {
		return func() error { return err }
	}

	if err := dockerfile.Execute(stdin, dockerContext{
		Paths:        cfg.Paths,
		Package:      tgt.Pkg,
		Dependencies: cfg.Graph.Descendants(tgt.Pkg).Slice(),
		BuildCommand: "go " + strings.Join(quote(generateBuildArgs(Build, cfg.Paths, tgt)), " "),
	}); err != nil {
		_ = cmd.Process.Kill()
		return func() error { return err }
	}

	stdin.Close()
	return cmd.Wait
}

func buildDockerTag(opts DockerOptions, tgt Target) string {
	tag := tgt.Binary + ":" + tgt.Version.String()
	if opts.Org != "" {
		tag = opts.Org + "/" + tag
	}
	if opts.Registry != "" {
		tag = opts.Registry + "/" + tag
	}
	return tag
}

func quote(strs []string) []string {
	for i, str := range strs {
		if strings.ContainsAny(str, " \t") {
			if strings.ContainsRune(str, '\'') {
				strs[i] = fmt.Sprintf(`"%s"`, str)
			} else {
				strs[i] = fmt.Sprintf("'%s'", str)
			}
		}
	}
	return strs
}
