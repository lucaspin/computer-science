package huffman

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/lucaspin/computer-science/pkg/binary_heap"
	"github.com/lucaspin/computer-science/pkg/stack"
)

type HuffmanNode struct {
	Value  byte
	Weight int
	Left   *HuffmanNode
	Right  *HuffmanNode
}

func (n *HuffmanNode) IsLeaf() bool {
	return n.Left == nil && n.Right == nil
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
	buffer := new(bytes.Buffer)
	encodeTree(buffer, tree)
	encodeInput(buffer, tree, input)

	return buffer, nil
}

func Decompress(input []byte) (*bytes.Buffer, error) {
	header, data, err := splitInput(input)
	if err != nil {
		return nil, err
	}

	tree := NewHuffmanTreeFromEncodedHeader(header)
	buf := new(bytes.Buffer)
	if err := decode(buf, data, tree); err != nil {
		return nil, fmt.Errorf("error decoding input: %v", err)
	}

	return buf, nil
}

func decode(buf *bytes.Buffer, data []byte, tree HuffmanTree) error {
	currentTreeNode := tree.root

	for _, b := range data {

		// Start at the most significant bit and look at each bit
		for i := 7; i >= 0; i-- {
			// Use bit to traverse tree
			if (b & (1 << i)) == 0 {
				currentTreeNode = currentTreeNode.Left
			} else {
				currentTreeNode = currentTreeNode.Right
			}

			// If current node is a leaf one, print its value, and go back to root.
			if currentTreeNode.IsLeaf() {
				if err := binary.Write(buf, binary.BigEndian, currentTreeNode.Value); err != nil {
					return fmt.Errorf("error writing value: %v", err)
				}

				currentTreeNode = tree.root
			}
		}
	}

	return nil
}

func splitInput(input []byte) ([]byte, []byte, error) {
	var headerSize int32
	err := binary.Read(bytes.NewReader(input[:4]), binary.BigEndian, &headerSize)
	if err != nil {
		return nil, nil, fmt.Errorf("error finding header size: %v", err)
	}

	fmt.Printf("Header size from 4-byte piece: %d\n", headerSize)

	header := input[4 : headerSize+4]
	encoded := input[headerSize+4:]

	fmt.Printf("Header size: %d\n", len(header))
	fmt.Printf("Encoded data size: %d\n", len(encoded))

	return header, encoded, nil
}

func encodeTree(buffer *bytes.Buffer, huffmanTree HuffmanTree) error {
	header := generateHeader(huffmanTree)
	headerLen := int32(len(header))

	fmt.Printf("Encoding - header size: %d\n", headerLen)

	// Write the header size (4 bytes)
	if err := binary.Write(buffer, binary.BigEndian, headerLen); err != nil {
		return fmt.Errorf("error writing header size: %v", err)
	}

	fmt.Printf("Header: %v\n", header)

	// Write the actual header.
	if err := binary.Write(buffer, binary.BigEndian, header); err != nil {
		return fmt.Errorf("error writing header: %v", err)
	}

	return nil
}

func generateHeader(huffmanTree HuffmanTree) []byte {
	b := []byte{}
	postOrderWalk(*huffmanTree.root, func(bytes []byte) {
		b = append(b, bytes...)
	})

	return b
}

func postOrderWalk(node HuffmanNode, callback func(bytes []byte)) {
	isNonLeaf := false

	if node.Left != nil {
		isNonLeaf = true
		postOrderWalk(*node.Left, callback)
	}

	if node.Right != nil {
		isNonLeaf = true
		postOrderWalk(*node.Right, callback)
	}

	// TODO: this should not use a whole byte to write 0x00
	if isNonLeaf {
		callback([]byte{0x00})
		return
	}

	// TODO: this should not use a whole byte to write 0xff
	callback([]byte{0xff, node.Value})
}

func encodeInput(buffer *bytes.Buffer, huffmanTree HuffmanTree, rawInput []byte) error {
	codes := huffmanTree.Codes()

	var currentWord uint32
	var currentBitsUsed uint8
	for _, b := range rawInput {
		code := codes[b]

		// Check if the code fits into the current word.
		// If not, part of it will have to go the next word
		if currentBitsUsed+code.BitsUsed > 32 {
			// Find out how many bits go into the current word
			nextWordBits := (currentBitsUsed + code.BitsUsed) - 32
			currentWordBits := code.BitsUsed - nextWordBits

			fmt.Printf(
				"Code '%s' does not fit into current word - word='%032b' next=%d, current=%d\n",
				code.String(), currentWord, nextWordBits, currentWordBits,
			)

			// Put only the bits that fit into the current word, and save it.
			currentWord = (currentWord << uint32(currentWordBits)) + (code.Code >> uint32(nextWordBits))
			fmt.Printf("Saving word '%032b'...\n", currentWord)
			err := binary.Write(buffer, binary.BigEndian, currentWord)
			if err != nil {
				return err
			}

			// Put the remaining bits in the next one.
			currentWord = code.Code >> currentWordBits
			fmt.Printf("Current word updated: '%032b'...\n", currentWord)
			currentBitsUsed = nextWordBits
			continue
		}

		fmt.Printf("Code '%s' fits into current word '%032b'\n", code.String(), currentWord)

		// It fits into the current word,
		// shift the current word to the left,
		// to open space for the new code, and add the code to the word.
		currentWord = (currentWord << uint32(code.BitsUsed)) + code.Code
		currentBitsUsed += code.BitsUsed

		fmt.Printf("Word updated: '%032b'\n", currentWord)
	}

	// Make sure to write the current word out
	return binary.Write(buffer, binary.BigEndian, currentWord)
}

// TODO
func NewHuffmanTreeFromEncodedHeader(header []byte) HuffmanTree {
	stack := stack.NewStack[HuffmanNode]()

	for i := 0; i < len(header); {

		// If we read a 0xff, we reached a leaf node.
		// Read the next byte (its value), and push it into the stack.
		if header[i] == 0xff {
			n := HuffmanNode{Value: header[i+1]}
			stack.Push(&n)
			i += 2
			continue
		}

		// If we read a 0x00, we check the stack.
		// If only one element is present in it, we are done.
		if stack.Len() == 1 {
			break
		}

		// If more elements are present, we combine the first two into a new one.
		// TODO: explain why the first goes to the right and second goes to the left.
		first, _ := stack.Pop()
		second, _ := stack.Pop()
		stack.Push(&HuffmanNode{Left: second, Right: first})
		i++
	}

	root, _ := stack.Pop()
	return HuffmanTree{root: root}
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

func (t *HuffmanTree) PrintCodes() {
	codes := t.Codes()
	for b, code := range codes {
		fmt.Printf("'%c': %s\n", b, code.String())
	}
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
