package types

import "sync"

type Graph[N comparable] interface {
	Has(N) bool
	PutNode(N)
	PutEdge(N, N)
	Nodes() Set[N]
	Transpose() Graph[N]
	Descendants(N) Set[N]
}

type graph[N comparable] Multimap[N, N]

func NewGraph[N comparable]() Graph[N] {
	return graph[N](make(Multimap[N, N]))
}

func (g graph[N]) Has(n N) bool {
	return Multimap[N, N](g).Has(n)
}

func (g graph[N]) PutNode(n N) {
	Multimap[N, N](g).Add(n)
}

func (g graph[N]) PutEdge(u N, v N) {
	Multimap[N, N](g).Put(u, v)
}

func (g graph[N]) Nodes() Set[N] {
	ns := make(Set[N])
	for k, vs := range Multimap[N, N](g) {
		ns.Add(k)
		ns.AddSet(vs)
	}
	return ns
}

func (g graph[N]) Transpose() Graph[N] {
	mm := make(Multimap[N, N])
	for k, vs := range Multimap[N, N](g) {
		if len(vs) == 0 {
			mm.Add(k)
		}

		for v := range vs {
			mm.Put(v, k)
		}
	}
	return graph[N](mm)
}

func (g graph[N]) Descendants(node N) Set[N] {
	ns := make(Set[N])
	q := make(chan N, 100)
	wg := newWaitGroup(q)

	q <- node
	for n := range q {
		if ns.Has(n) {
			wg.Done()
			continue
		}

		if n != node {
			ns.Add(n)
		}

		children := Multimap[N, N](g)[n]
		wg.Add(len(children))
		go enqueue(q, children)

		wg.Done()
	}

	return ns
}

func enqueue[N comparable](q chan<- N, ns Set[N]) {
	for n := range ns {
		q <- n
	}
}

func newWaitGroup[N any](q chan<- N) *sync.WaitGroup {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		wg.Wait()
		close(q)
	}()
	return wg
}
