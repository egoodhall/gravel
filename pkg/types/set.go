package types

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](ts ...T) Set[T] {
	s := make(Set[T])
	switch len(ts) {
	case 0:
	case 1:
		s.Add(ts[0])
	default:
		s.Add(ts[0], ts[1:]...)
	}
	return s
}

func (s Set[T]) Add(t T, ts ...T) {
	s[t] = struct{}{}
	for _, t := range ts {
		s[t] = struct{}{}
	}
}

func (s Set[T]) AddSet(ts Set[T]) {
	for t := range ts {
		s.Add(t)
	}
}

func (s Set[T]) Has(t T) bool {
	_, ok := s[t]
	return ok
}

func (s Set[T]) Filter(filter func(T) bool) Set[T] {
	s2 := make(Set[T])
	for t := range s {
		if filter(t) {
			s2.Add(t)
		}
	}
	return s2
}

func (s Set[T]) Slice() []T {
	s2 := make([]T, len(s))

	i := 0
	for t := range s {
		s2[i] = t
		i++
	}

	return s2
}
