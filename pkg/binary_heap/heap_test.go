package binary_heap

import (
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test__MaxHeap(t *testing.T) {
	rand.Seed(time.Now().Unix())

	list := createRandomList(20)
	heap := NewBinaryHeap(list, func(a, b int) bool { return a > b })

	ns := []int{}
	for heap.Len() > 0 {
		n := heap.Pop()
		ns = append(ns, *n)
	}

	// Note: i and j are indexes and j is i-1
	assert.True(t, sort.SliceIsSorted(ns, func(i, j int) bool {
		return ns[i] > ns[j]
	}))
}

func Test__MinHeap(t *testing.T) {
	rand.Seed(time.Now().Unix())

	list := createRandomList(20)
	heap := NewBinaryHeap(list, func(a, b int) bool { return a < b })

	ns := []int{}
	for heap.Len() > 0 {
		n := heap.Pop()
		ns = append(ns, *n)
	}

	// Note: i and j are indexes and j is i-1
	assert.True(t, sort.SliceIsSorted(ns, func(i, j int) bool {
		return ns[i] < ns[j]
	}))
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
