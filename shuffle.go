// Package fisheryates shuffles collections of data.
package fisheryates

import (
	"math/rand"
	"sort"
)

// Interface models a type that can be shuffled by this package's functions.
// Any type that fulfills sort.Interface may be used.
type Interface interface {
	// Swap exchanges elements indexed by i and j.  See
	// sort.Interface.Swap.
	Swap(i, j int)
	// Len reports the number of elements in the collection.  See
	// sort.Interface.Len.
	Len() int
}

var _ Interface = sort.Interface(nil)

// Rand emits a random value along the interval 0 <= v < n.  See rand.Intn.
type Rand func(n int) int

var defaultRand = rand.Intn

// Shuffle permutates the data randomly according to the randomizer r.  If r is
// nil, the function defaults to rand.Intn, which while convenient puts
// concurrent callers of shuffle or package sort's default functions under the
// same lock.
func Shuffle(data Interface, r Rand) {
	if r == nil {
		r = defaultRand
	}
	for i := 0; i < data.Len()-2; i++ {
		data.Swap(i, r(data.Len()-i))
	}
}
