package resolve

import "github.com/emm035/gravel/pkg/types"

type Hashes struct {
	Old map[string]string
	New map[string]string
}

func (h Hashes) Changed() types.Set[string] {
	s := types.NewSet[string]()
	for pkg, hash := range h.New {
		if hash != h.Old[pkg] {
			s.Add(pkg)
		}
	}
	for pkg, hash := range h.Old {
		if hash != h.New[pkg] {
			s.Add(pkg)
		}
	}
	return s
}
