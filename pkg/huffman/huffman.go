package huffman

import (
	"fmt"
	"strconv"

	"github.com/lucaspin/computer-science/pkg/binary_heap"
)

type HuffmanNode struct {
	Value  byte
	Weight int
	Left   *HuffmanNode
	Right  *HuffmanNode
}

type HuffmanTree struct {
	root *HuffmanNode
}

// TODO: document this function properly
func NewHuffmanTree(nodes binary_heap.BinaryHeap[HuffmanNode]) HuffmanTree {
	// We iterate until only a single element is left at our heap
	for nodes.Len() > 1 {
		first := nodes.Pop()
		second := nodes.Pop()
		nodes.Push(HuffmanNode{
			Left:   first,
			Right:  second,
			Weight: first.Weight + second.Weight,
		})
	}

	// The element left in our heap is the root of our huffman tree.
	return HuffmanTree{root: nodes.Pop()}
}

func (t *HuffmanTree) Codes() map[byte]string {
	codes := map[byte]string{}

	var traverse func(int32, int, HuffmanNode)
	traverse = func(code int32, bitsUsed int, node HuffmanNode) {
		// This is a leaf node, just save the code
		if node.Left == nil && node.Right == nil {
			codes[node.Value] = fmt.Sprintf("%0"+strconv.Itoa(int(bitsUsed))+"b", code)
			return
		}

		traverse(code<<1, bitsUsed+1, *node.Left)
		traverse(code<<1+1, bitsUsed+1, *node.Right)
	}

	traverse(0, 0, *t.root)
	return codes
}

// Find the frequencies of the characters in a stream of bytes
// and return a binary heap where each element represent a byte and its weight.
// TODO: receive file name as argument instead of []byte
func FindFrequencies(content []byte) (binary_heap.BinaryHeap[HuffmanNode], error) {
	freqs := map[byte]int{}

	// TODO: doing byte per byte but we should consider bigger chunks for a real implementation
	for _, b := range content {
		freqs[b]++
	}

	// Build list of huffman nodes
	list := []HuffmanNode{}
	for k, v := range freqs {
		n := HuffmanNode{Value: k, Weight: v}
		list = append(list, n)
	}

	// we want the characters with lower frequencies to be first, so we use a min-heap
	return binary_heap.NewBinaryHeap(
		list,
		func(a, b HuffmanNode) bool {
			return a.Weight < b.Weight
		},
	), nil
}
