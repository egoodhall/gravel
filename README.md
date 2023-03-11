# gravel

Gravel is a build tool for Go monorepos. To speed up builds, gravel caches a hash of the files in a package that have changed since the last build, and only tests/builds the packages that depend on them. Gravel uses Go modules for managing dependencies, which makes it more feasible to use than something like Bazel for small/medium-sized projects.

### Installation

```
$ go install github.com/emm035/gravel@latest
```

### Usage

```
$ gravel build
```

### Embedding Version Information

Binaries built by gravel can have version embedded as shown below:
```go
package example

import (
  "fmt"

  "github.com/emm035/gravel/pkg/buildinfo"
)

func PrintBuildInfo() {
  // These values will be set by gravel during the build
  fmt.Println("version:", buildinfo.GetVersion())
  fmt.Println("commit: ", buildinfo.GetCommit())
}
```


### Output

All binaries are placed into `$root/gravel/bin`.
