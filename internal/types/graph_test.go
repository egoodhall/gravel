package types_test

import (
	"reflect"
	"testing"

	"github.com/emm035/gravel/internal/types"
)

func TestGraph(t *testing.T) {
	g := types.NewGraph[string]()
	g.PutEdge("a", "b")
	g.PutEdge("a", "d")
	g.PutEdge("b", "c")
	g.PutEdge("b", "e")

	if !reflect.DeepEqual(g.Descendants("a"), types.NewSet("b", "c", "d", "e")) {
		t.Fatal()
	}
	if !reflect.DeepEqual(g.Descendants("b"), types.NewSet("c", "e")) {
		t.Fatal()
	}
	if !reflect.DeepEqual(g.Descendants("c"), types.NewSet[string]()) {
		t.Fatal()
	}
}
