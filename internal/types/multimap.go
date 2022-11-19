package types

type Multimap[K, V comparable] map[K]Set[V]

func (mm Multimap[K, V]) Add(k K) {
	if _, ok := mm[k]; !ok {
		mm[k] = make(Set[V])
	}
}

func (mm Multimap[K, V]) Put(k K, v V, vs ...V) {
	if _, ok := mm[k]; !ok {
		mm[k] = make(Set[V])
	}
	mm[k].Add(v, vs...)
}

func (mm Multimap[K, V]) PutSet(k K, vs Set[V]) {
	if _, ok := mm[k]; !ok {
		mm[k] = make(Set[V])
	}
	mm[k].AddSet(vs)
}

func (mm Multimap[K, V]) Has(i any) bool {
	for k, vs := range mm {
		if k == i {
			return true
		}
		for v := range vs {
			if v == i {
				return true
			}
		}
	}
	return false
}

func (mm Multimap[K, V]) Keys() Set[K] {
	ks := make(Set[K])
	for k := range mm {
		ks.Add(k)
	}
	return ks
}

func (mm Multimap[K, V]) Values() Set[V] {
	vs := make(Set[V])
	for _, vs := range mm {
		vs.AddSet(vs)
	}
	return vs
}
