package semver

import (
	"errors"
	"fmt"
	"strings"
)

var ErrUnknownMod = errors.New("unknown version mod")

//go:generate go run golang.org/x/tools/cmd/stringer -type=Mod -linecomment
type Mod byte

const (
	ModPatch Mod = iota // patch
	ModMinor            // minor
	ModMajor            // major
	ModDate             // date
)

func (seg *Mod) UnmarshalText(p []byte) error {
	switch strings.ToLower(string(p)) {
	case ModMajor.String():
		*seg = ModMajor
	case ModMinor.String():
		*seg = ModMinor
	case ModPatch.String():
		*seg = ModPatch
	case ModDate.String():
		*seg = ModDate
	default:
		return fmt.Errorf("%w: %s", ErrUnknownMod, string(p))
	}
	return nil
}

func (seg Mod) MarshalText() ([]byte, error) {
	return []byte(seg.String()), nil
}
