package huffman

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test__FindFrequencies(t *testing.T) {
	h, err := FindFrequencies([]byte("go go gophers"))
	assert.NoError(t, err)
	assert.Equal(t, h.Len(), 8)
}

func Test__Codes(t *testing.T) {
	h, _ := FindFrequencies([]byte("go go gophers"))
	tree := NewHuffmanTree(h)
	codes := tree.Codes()

	for b, code := range codes {
		fmt.Printf("'%c': %s\n", b, code)
	}
}
