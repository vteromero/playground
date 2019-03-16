package simple

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecompressor_NewDecompressor(t *testing.T) {
	c := NewDecompressor(OrderAscending, 10)
	assert.Equal(t, c.ListOrder, OrderAscending)
	assert.Equal(t, c.CardinalityHeaderSize, 10)
}

func TestDecompressor_Decompress(t *testing.T) {
	params := []struct {
		decompressor   *Decompressor
		input          []byte
		expectedErr    error
		expectedOutput []uint32
	}{
		{NewDecompressor(OrderDescending, 8), []byte{0x01, 0xb8, 0x22, 0x00, 0x00}, nil, []uint32{8888}},
		{NewDecompressor(OrderAscending, 8), []byte{0x01, 0xb8, 0x22, 0x00, 0x00}, nil, []uint32{8888}},
		{NewDecompressor(OrderDescending, 8), []byte{0x03, 0xb8, 0x22, 0x00, 0x00, 0x6f, 0x40, 0x01}, nil, []uint32{8888, 111, 5}},
		{NewDecompressor(OrderAscending, 8), []byte{0x03, 0xb8, 0x22, 0x00, 0x00, 0x6f, 0x40, 0x01}, nil, []uint32{5, 111, 8888}},
	}

	for _, testCase := range params {
		output, err := testCase.decompressor.Decompress(testCase.input)
		assert.Equal(t, testCase.expectedErr, err)
		assert.Equal(t, testCase.expectedOutput, output)
	}
}
