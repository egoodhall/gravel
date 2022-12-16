package semver

import (
	"errors"
	"fmt"
	"strings"
)

var ErrUnknownStrategy = errors.New("unknown strategy")

//go:generate go run golang.org/x/tools/cmd/stringer -type=Strategy -linecomment
type Strategy byte

const (
	StrategyDate Strategy = iota // date
)

func (strat *Strategy) UnmarshalText(p []byte) error {
	switch strings.ToLower(string(p)) {
	case StrategyDate.String():
		*strat = StrategyDate
	default:
		return fmt.Errorf("%w: %s", ErrUnknownSegment, string(p))
	}
	return nil
}

func (strat Strategy) MarshalText() ([]byte, error) {
	return []byte(strat.String()), nil
}
