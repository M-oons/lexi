package random

import "testing"

func TestNewGenerator_SeededDeterministic(t *testing.T) {
	const (
		seed = int64(42)
		n    = 100
	)

	g1 := NewGenerator(seed)
	g2 := NewGenerator(seed)

	for i := range 20 {
		v1 := g1.Intn(n)
		v2 := g2.Intn(n)

		if v1 != v2 {
			t.Fatalf("determinism mismatch at iteration %d: %d != %d", i, v1, v2)
		}
		if v1 < 0 || v1 >= n {
			t.Fatalf("value out of range at iteration %d: %d", i, v1)
		}
	}
}

func TestNewGenerator_NegativeSeed(t *testing.T) {
	g := NewGenerator(-1)

	for _, n := range []int{1, 2, 10, 100} {
		for range 20 {
			v := g.Intn(n)
			if v < 0 || v >= n {
				t.Fatalf("Intn(%d) returned out of range value %d", n, v)
			}
		}
	}
}
