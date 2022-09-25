package gravel

import _ "embed"

//go:generate sh version.sh
//go:embed version.txt
var Version string
