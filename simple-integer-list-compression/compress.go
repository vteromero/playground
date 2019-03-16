package simple

import (
	"github.com/vteromero/bitstream"
)

const (
	OrderAscending = iota
	OrderDescending
)

type Compressor struct {
	ListOrder             int
	CardinalityHeaderSize int
	input                 []uint32
	writer                *bitstream.Writer
}

func NewCompressor(order int, cardHeaderSize int) *Compressor {
	return &Compressor{
		ListOrder:             order,
		CardinalityHeaderSize: cardHeaderSize,
		input:  nil,
		writer: nil,
	}
}

func (c *Compressor) isCardinalityHeaderSizeValid() bool {
	return c.CardinalityHeaderSize >= 1 && c.CardinalityHeaderSize <= 32
}

func (c *Compressor) isInputSizeValid(size int) bool {
	return size >= 0 && size < (1<<uint(c.CardinalityHeaderSize))
}

func (c *Compressor) writeCardinality() error {
	return c.writer.Write(uint64(len(c.input)), c.CardinalityHeaderSize)
}

func (c *Compressor) writeValuesAsc() error {
	w := 32
	for i := len(c.input) - 1; i >= 0; i-- {
		value := c.input[i]
		if err := c.writer.Write(uint64(value), w); err != nil {
			return err
		}
		w = bitsLen(value)
	}
	return nil
}

func (c *Compressor) writeValuesDesc() error {
	w := 32
	for _, value := range c.input {
		if err := c.writer.Write(uint64(value), w); err != nil {
			return err
		}
		w = bitsLen(value)
	}
	return nil
}

func (c *Compressor) writeValues() error {
	if c.ListOrder == OrderAscending {
		return c.writeValuesAsc()
	}
	return c.writeValuesDesc()
}

func (c *Compressor) Compress(input []uint32, output []byte) (int, error) {
	if !c.isCardinalityHeaderSizeValid() {
		return 0, ErrCardinalityHeaderSizeOutOfBound
	}

	if !c.isInputSizeValid(len(input)) {
		return 0, ErrInputTooLong
	}

	c.input = input
	c.writer = bitstream.NewWriter(output)

	if err := c.writeCardinality(); err != nil {
		return 0, err
	}

	if err := c.writeValues(); err != nil {
		return 0, err
	}

	return c.writer.Offset(), nil
}

func (c *Compressor) MaxCompressedLen(n int) int {
	if !c.isCardinalityHeaderSizeValid() || !c.isInputSizeValid(n) {
		return 0
	}
	return sizeInBytes(c.CardinalityHeaderSize + 32*n)
}
