package huffman

import (
	"bytes"
	"encoding/binary"
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

func Compress(input []byte) (*bytes.Buffer, error) {
	// Build encoding table
	f, err := FindFrequencies(input)
	if err != nil {
		return nil, fmt.Errorf("error finding character frequencies")
	}

	tree := NewHuffmanTree(f)
	codes := tree.Codes()
	// TODO: Write tree into encoded output too

	// Encode input
	buf := new(bytes.Buffer)

	var currentWord uint32
	var currentBitsUsed uint8
	for _, b := range input {
		code := codes[b]

		// Check if the code fits into the current word.
		// If not, part of it will have to go the next word
		if currentBitsUsed+code.BitsUsed > 32 {
			// Find out how many bits go into the current word
			nextWordBits := (currentBitsUsed + code.BitsUsed) - 32
			currentWordBits := code.BitsUsed - nextWordBits

			// Select only the bits you are interested
			// and put them into the current word.
			result := code.Code & findMask(currentWordBits)
			currentWord |= result << currentBitsUsed

			// Save the current word, and put the remaining bits in the next one.
			binary.Write(buf, binary.BigEndian, currentWord)
			currentWord = code.Code >> nextWordBits
			currentBitsUsed = nextWordBits
			continue
		}

		// It fits into the current word,
		// shift the code {currentBitsUsed} times to the left
		// and OR the resulting mask with the current word.
		// Increment the number of bits used.
		mask := code.Code << currentBitsUsed
		currentWord |= mask
		currentBitsUsed += code.BitsUsed
	}

	// Make sure to write the current word out
	binary.Write(buf, binary.BigEndian, currentWord)

	return buf, nil
}

// TODO: document this properly
func findMask(leftBits uint8) uint32 {
	var mask uint32

	for leftBits > 0 {
		mask <<= 1
		mask += 1
		leftBits--
	}

	return mask
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

type HuffmanCode struct {
	Code     uint32
	BitsUsed uint8
}

func (c *HuffmanCode) String() string {
	return fmt.Sprintf("%0"+strconv.Itoa(int(c.BitsUsed))+"b", c.Code)
}

func (t *HuffmanTree) Codes() map[byte]HuffmanCode {
	codes := map[byte]HuffmanCode{}

	var traverse func(uint32, uint8, HuffmanNode)
	traverse = func(code uint32, bitsUsed uint8, node HuffmanNode) {
		// This is a leaf node, just save the code
		if node.Left == nil && node.Right == nil {
			codes[node.Value] = HuffmanCode{Code: code, BitsUsed: bitsUsed}
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
