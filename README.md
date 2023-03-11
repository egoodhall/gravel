# gravel

Build tool for Golang monorepos. To speed up builds, gravel keeps track of what files have
changed since the last build, and only tests/builds their packages. Gravel uses Go modules
for managing dependencies, which makes it more feasible to use than something like Bazel for
small projects.

### Installation

```
$ go install github.com/emm035/gravel@latest
```

### Usage

```
$ gravel build
```

### Output

All built binaries are placed into `$root/gravel/bin`. For example, if run on this repository:

```
$ gravel build
$ ls ./gravel/bin
gravel
```

### Embedding Version Information

Binaries built by gravel can have their version embedded as shown below:
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
