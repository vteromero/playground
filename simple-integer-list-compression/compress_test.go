package simple

import (
	"testing"

	"github.com/vteromero/playground/simple-integer-list-compression/slice"
	"github.com/stretchr/testify/assert"
)

func TestCompressor_NewCompressor(t *testing.T) {
	c := NewCompressor(OrderAscending, 10)
	assert.Equal(t, c.ListOrder, OrderAscending)
	assert.Equal(t, c.CardinalityHeaderSize, 10)
}

func TestCompressor_Compress(t *testing.T) {
	params := []struct {
		compressor     *Compressor
		input          []uint32
		expectedN      int
		expectedErr    error
		expectedOutput []byte
	}{
		{NewCompressor(OrderDescending, 0), []uint32{}, 0, ErrCardinalityHeaderSizeOutOfBound, []byte{}},
		{NewCompressor(OrderDescending, -10), []uint32{}, 0, ErrCardinalityHeaderSizeOutOfBound, []byte{}},
		{NewCompressor(OrderDescending, 33), []uint32{}, 0, ErrCardinalityHeaderSizeOutOfBound, []byte{}},
		{NewCompressor(OrderAscending, 2), []uint32{1, 2, 3, 4}, 0, ErrInputTooLong, []byte{}},
		{NewCompressor(OrderAscending, 3), []uint32{1, 2, 3, 4, 5, 6, 7, 8}, 0, ErrInputTooLong, []byte{}},
		{NewCompressor(OrderDescending, 8), []uint32{}, 8, nil, []byte{0x00}},
		{NewCompressor(OrderDescending, 8), []uint32{8888}, 40, nil, []byte{0x01, 0xb8, 0x22, 0x00, 0x00}},
		{NewCompressor(OrderAscending, 8), []uint32{8888}, 40, nil, []byte{0x01, 0xb8, 0x22, 0x00, 0x00}},
		{NewCompressor(OrderDescending, 8), []uint32{8888, 111, 5}, 61, nil, []byte{0x03, 0xb8, 0x22, 0x00, 0x00, 0x6f, 0x40, 0x01}},
		{NewCompressor(OrderAscending, 8), []uint32{5, 111, 8888}, 61, nil, []byte{0x03, 0xb8, 0x22, 0x00, 0x00, 0x6f, 0x40, 0x01}},
	}

	for _, testCase := range params {
		n := len(testCase.input)
		output := make([]byte, testCase.compressor.MaxCompressedLen(n))
		m, err := testCase.compressor.Compress(testCase.input, output)
		assert.Equal(t, testCase.expectedN, m)
		assert.Equal(t, testCase.expectedErr, err)
		assert.Equal(t, testCase.expectedOutput, output[:sizeInBytes(m)])
	}
}

func TestCompressor_MaxCompressedLen(t *testing.T) {
	params := []struct {
		compressor *Compressor
		n          int
		expectedN  int
	}{
		{NewCompressor(OrderAscending, 1), 0, 1},
		{NewCompressor(OrderAscending, 1), 1, 5},
		{NewCompressor(OrderDescending, 7), 1, 5},
		{NewCompressor(OrderDescending, 7), 100, 401},
		{NewCompressor(OrderAscending, 31), 10, 44},
		{NewCompressor(OrderAscending, 31), 200, 804},
		{NewCompressor(OrderDescending, 31), 8000, 32004},
		{NewCompressor(OrderDescending, 31), 15000, 60004},
		{NewCompressor(OrderAscending, -10), 100, 0},
		{NewCompressor(OrderAscending, 33), 100, 0},
		{NewCompressor(OrderAscending, 2), 4, 0},
		{NewCompressor(OrderAscending, 2), 100, 0},
		{NewCompressor(OrderAscending, 8), 256, 0},
	}

	for _, testCase := range params {
		compLen := testCase.compressor.MaxCompressedLen(testCase.n)
		assert.Equal(t, testCase.expectedN, compLen)
	}
}

func TestCompressAndDecompress(t *testing.T) {
	params := []struct {
		order                 int
		cardinalityHeaderSize int
		inputSize             int
	}{
		{OrderAscending, 8, 0},
		{OrderAscending, 8, 10},
		{OrderAscending, 8, 100},
		{OrderAscending, 16, 0},
		{OrderAscending, 16, 10},
		{OrderAscending, 16, 100},
		{OrderAscending, 16, 1000},
		{OrderAscending, 32, 0},
		{OrderAscending, 32, 10},
		{OrderAscending, 32, 100},
		{OrderAscending, 32, 1000},
		{OrderDescending, 8, 0},
		{OrderDescending, 8, 10},
		{OrderDescending, 8, 100},
		{OrderDescending, 16, 0},
		{OrderDescending, 16, 10},
		{OrderDescending, 16, 100},
		{OrderDescending, 16, 1000},
		{OrderDescending, 32, 0},
		{OrderDescending, 32, 10},
		{OrderDescending, 32, 100},
		{OrderDescending, 32, 1000},
	}

	for _, testCase := range params {
		var input []uint32
		if testCase.order == OrderAscending {
			input = slice.SortAscUint32Slice(slice.RandomUint32Slice(testCase.inputSize))
		} else {
			input = slice.SortDescUint32Slice(slice.RandomUint32Slice(testCase.inputSize))
		}

		c := NewCompressor(testCase.order, testCase.cardinalityHeaderSize)
		compOutput := make([]byte, c.MaxCompressedLen(testCase.inputSize))
		_, err := c.Compress(input, compOutput)
		assert.Nil(t, err)

		d := NewDecompressor(testCase.order, testCase.cardinalityHeaderSize)
		output, err := d.Decompress(compOutput)
		assert.Nil(t, err)

		assert.Equal(t, input, output)
	}
}
