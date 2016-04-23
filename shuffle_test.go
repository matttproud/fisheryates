package fisheryates

import (
	"flag"
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"testing/quick"
)

func newRand() Rand { return rand.New(rand.NewSource(42)).Intn }

func TestEquivalence(t *testing.T) {
	t.Parallel()
	f := func(d sort.IntSlice) sort.IntSlice {
		cpy := make(sort.IntSlice, len(d))
		copy(cpy, d)
		Shuffle(cpy, newRand())
		return cpy
	}
	if err := quick.CheckEqual(f, f, nil); err != nil {
		t.Errorf("expected f() == f(): %v", err)
	}
}

var iterationDepth int

func TestSatisfactoryShuffle(t *testing.T) {
	t.Parallel()
	const shuffleLevel = 0.5
	rnd := newRand()
	chk := func(d sort.IntSlice) bool {
		if len(d) < 3 {
			return true
		}
		mutated := make(sort.IntSlice, len(d))
		for i := 0; i < iterationDepth; i++ {
			equal := 0
			copy(mutated, d)
			Shuffle(mutated, rnd)
			for j, mVal := range mutated {
				if d[j] == mVal {
					equal++
				}
			}
			if float64(equal)/float64(d.Len()) < shuffleLevel {
				return true
			}
		}
		t.Errorf("expected to achieve shuffle rate of %f after %d runs", shuffleLevel, iterationDepth)
		return false
	}
	if err := quick.Check(chk, nil); err != nil {
		t.Error(err)
	}
}

func TestFullSpectrumMutated(t *testing.T) {
	t.Parallel()
	data := sort.IntSlice{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	cpy := make(sort.IntSlice, len(data))
	rnd := newRand()
	notMutated := map[int]bool{
		0: true,
		1: true,
		2: true,
		3: true,
		4: true,
		5: true,
		6: true,
		7: true,
		8: true,
		9: true,
	}
	for i := 0; i < iterationDepth; i++ {
		copy(cpy, data)
		Shuffle(cpy, rnd)
		for shfIdx, shfVal := range cpy {
			if shfVal != data[shfIdx] {
				delete(notMutated, shfIdx)
			}
		}
	}
	if len(notMutated) != 0 {
		t.Errorf("expected %d shuffles to have removed all items; %#v remain", iterationDepth, notMutated)

	}
}

func BenchmarkStandard(b *testing.B) {
	data := sort.IntSlice{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	rnd := newRand()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Shuffle(data, rnd)
	}
}

func ExampleShuffle() {
	data := sort.IntSlice{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	rnd := rand.New(rand.NewSource(4)) // Not required; want determinism.
	Shuffle(data, rnd.Intn)
	fmt.Println(data)
	// Output: [6 7 2 4 1 5 9 3 8 0]
}

func init() {
	flag.IntVar(&iterationDepth, "iteration_depth", 10, "no. of runs for test that await satisfactory condition.")
}
