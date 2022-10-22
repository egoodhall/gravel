# gravel

Build tool for Golang monorepos. To speed up builds, gravel keeps track of what files have
changed since the last build, and only tests/builds their packages. Gravel uses Go modules
for managing dependencies, which makes it more feasible to use than something like Bazel for
small projects.

All built binaries are placed into `$root/gravel/bin`. For example, if run on this repository:

```
$ gravel build
$ ls ./gravel/bin
gravel
```

### Installation

```
$ go install github.com/emm035/gravel@latest
```

### Usage

```
$ gravel build
```
