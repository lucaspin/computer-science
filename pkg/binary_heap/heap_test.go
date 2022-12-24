package binary_heap

import (
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test__SuboptimalToHeap(t *testing.T) {
	rand.Seed(time.Now().Unix())

	list := createRandomList(500)
	heap := SuboptimalToHeap(list)

	ns := []int{}
	for len(heap) > 0 {
		h, n := Pop(heap)
		ns = append(ns, n)
		heap = h
	}

	// Note: j is i-1 here
	assert.True(t, sort.SliceIsSorted(ns, func(i, j int) bool { return i < j }))
}

func createRandomList(size int) []int {
	list := []int{}
	for size > 0 {
		n := rand.Intn(1000)
		list = append(list, n)
		size--
	}

	return list
}
