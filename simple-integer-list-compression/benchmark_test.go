package simple

import (
	"bytes"
	"compress/zlib"
	"io"
	"testing"

	"github.com/vteromero/playground/simple-integer-list-compression/slice"
	"github.com/dataence/encoding"
	"github.com/dataence/encoding/bp32"
	"github.com/dataence/encoding/composition"
	"github.com/dataence/encoding/cursor"
	deltabp32 "github.com/dataence/encoding/delta/bp32"
	deltafastpfor "github.com/dataence/encoding/delta/fastpfor"
	deltavb "github.com/dataence/encoding/delta/variablebyte"
	"github.com/dataence/encoding/fastpfor"
	"github.com/dataence/encoding/variablebyte"
)

var (
	sliceLen          = 10000000
	sortedInt32Slice  = slice.SortAscInt32Slice(slice.RandomInt31Slice(sliceLen))
	sortedUint32Slice = slice.Int32ToUint32Slice(sortedInt32Slice)
	byteSlice         = slice.Uint32ToByteSlice(sortedUint32Slice)
)

func benchmarkCompressEncodingLibWithCodec(b *testing.B, codec encoding.Integer) {
	out := make([]int32, sliceLen*2)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		inpos := cursor.New()
		outpos := cursor.New()
		codec.Compress(sortedInt32Slice, inpos, sliceLen, out, outpos)
	}
}

func benchmarkDecompressEncodingLibWithCodec(b *testing.B, codec encoding.Integer) {
	compInLen := sliceLen
	compOut := make([]int32, compInLen*2)
	compInpos := cursor.New()
	compOutpos := cursor.New()
	codec.Compress(sortedInt32Slice, compInpos, compInLen, compOut, compOutpos)
	compOutLen := compOutpos.Get()

	uncompInLen := compOutLen
	uncompOut := make([]int32, sliceLen)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		uncompInpos := cursor.New()
		uncompOutpos := cursor.New()
		codec.Uncompress(compOut, uncompInpos, uncompInLen, uncompOut, uncompOutpos)
	}
}

func BenchmarkCompressSimple(b *testing.B) {
	c := NewCompressor(OrderAscending, 32)
	out := make([]byte, c.MaxCompressedLen(sliceLen))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Compress(sortedUint32Slice, out)
	}
}

func BenchmarkCompressZlib(b *testing.B) {
	var output bytes.Buffer
	writer, _ := zlib.NewWriterLevel(&output, zlib.BestSpeed)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		writer.Write(byteSlice)
		writer.Close()
	}
}

func BenchmarkCompressBP32(b *testing.B) {
	benchmarkCompressEncodingLibWithCodec(b, composition.New(bp32.New(), variablebyte.New()))
}

func BenchmarkCompressFastpfor(b *testing.B) {
	benchmarkCompressEncodingLibWithCodec(b, composition.New(fastpfor.New(), variablebyte.New()))
}

func BenchmarkCompressDeltaBP32(b *testing.B) {
	benchmarkCompressEncodingLibWithCodec(b, composition.New(deltabp32.New(), deltavb.New()))
}

func BenchmarkCompressDeltaFastpfor(b *testing.B) {
	benchmarkCompressEncodingLibWithCodec(b, composition.New(deltafastpfor.New(), deltavb.New()))
}

func BenchmarkDecompressSimple(b *testing.B) {
	c := NewCompressor(OrderAscending, 32)
	data := make([]byte, c.MaxCompressedLen(sliceLen))
	c.Compress(sortedUint32Slice, data)

	d := NewDecompressor(OrderDescending, 32)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d.Decompress(data)
	}
}

func BenchmarkDecompressZlib(b *testing.B) {
	var buff bytes.Buffer

	writer, _ := zlib.NewWriterLevel(&buff, zlib.BestSpeed)
	writer.Write(byteSlice)
	writer.Close()
	compressed := buff.Bytes()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bytesReader := bytes.NewReader(compressed)
		var bytesBuff bytes.Buffer
		reader, _ := zlib.NewReader(bytesReader)
		io.Copy(&bytesBuff, reader)
		reader.Close()
	}
}

func BenchmarkDecompressBP32(b *testing.B) {
	benchmarkDecompressEncodingLibWithCodec(b, composition.New(bp32.New(), variablebyte.New()))
}

func BenchmarkDecompressFastpfor(b *testing.B) {
	benchmarkDecompressEncodingLibWithCodec(b, composition.New(fastpfor.New(), variablebyte.New()))
}

func BenchmarkDecompressDeltaBP32(b *testing.B) {
	benchmarkDecompressEncodingLibWithCodec(b, composition.New(deltabp32.New(), deltavb.New()))
}

func BenchmarkDecompressDeltaFastpfor(b *testing.B) {
	benchmarkDecompressEncodingLibWithCodec(b, composition.New(deltafastpfor.New(), deltavb.New()))
}
