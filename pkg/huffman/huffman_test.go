package huffman

import (
	"bytes"
	"compress/flate"
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
	tree.PrintCodes()
}

func Test__Compress(t *testing.T) {
	text := "let us write a whole lot of characters here"
	output, err := Compress([]byte(text))
	assert.NoError(t, err)

	raw, err := Decompress(output.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, text, raw.String())
}

func Test__Flate(t *testing.T) {
	b := new(bytes.Buffer)
	writer, err := flate.NewWriter(b, flate.HuffmanOnly)
	assert.NoError(t, err)

	n, err := writer.Write([]byte("go go gophersgo go gophersgo go gophersgo go gophersgo go gophersgo go gophers"))
	assert.NoError(t, err)
	fmt.Printf("Wrote %d bytes.\n", n)
	err = writer.Close()
	assert.NoError(t, err)

	fmt.Printf("Length: %d.\n", b.Len())
	reader := flate.NewReader(bytes.NewReader(b.Bytes()))
	output := make([]byte, 256)
	n, _ = reader.Read(output)
	fmt.Printf("Read %d bytes.\n", n)
	fmt.Printf("Output: '%s'\n", string(output))
}
