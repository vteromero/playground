package simple

import (
	"github.com/vteromero/bitstream"
)

type Decompressor struct {
	ListOrder             int
	CardinalityHeaderSize int
	cardinality           int
	input                 []byte
	reader                *bitstream.Reader
}

func NewDecompressor(order int, cardHeaderSize int) *Decompressor {
	return &Decompressor{
		ListOrder:             order,
		CardinalityHeaderSize: cardHeaderSize,
		cardinality:           0,
		input:                 nil,
		reader:                nil,
	}
}

func (d *Decompressor) readCardinality() error {
	v, err := d.reader.Read(d.CardinalityHeaderSize)
	d.cardinality = int(v)
	return err
}

func (d *Decompressor) readValuesAsc(output []uint32) ([]uint32, error) {
	w := 32

	for i := d.cardinality - 1; i >= 0; i-- {
		v, err := d.reader.Read(w)
		if err != nil {
			return nil, err
		}

		output[i] = uint32(v)

		w = bitsLen(uint32(v))
	}

	return output, nil
}

func (d *Decompressor) readValuesDesc(output []uint32) ([]uint32, error) {
	w := 32

	for i := 0; i < d.cardinality; i++ {
		v, err := d.reader.Read(w)
		if err != nil {
			return nil, err
		}

		output[i] = uint32(v)

		w = bitsLen(uint32(v))
	}

	return output, nil
}

func (d *Decompressor) readValues() ([]uint32, error) {
	values := make([]uint32, d.cardinality)
	if d.ListOrder == OrderAscending {
		return d.readValuesAsc(values)
	}
	return d.readValuesDesc(values)
}

func (d *Decompressor) Decompress(input []byte) ([]uint32, error) {
	d.input = input
	d.reader = bitstream.NewReader(input)

	if err := d.readCardinality(); err != nil {
		return nil, err
	}

	return d.readValues()
}
