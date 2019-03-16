package main

import (
	"bytes"
	"compress/zlib"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/vteromero/playground/simple-integer-list-compression"
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

type compressor struct {
	name     string
	compress func([]int32) int
}

func simpleCompress(cardHeaderSize int) func([]int32) int {
	return func(data []int32) int {
		compressor := simple.NewCompressor(simple.OrderAscending, cardHeaderSize)
		in := slice.Int32ToUint32Slice(data)
		out := make([]byte, compressor.MaxCompressedLen(len(in)))
		numBits, err := compressor.Compress(in, out)
		if err != nil {
			panic(err)
		}
		numBytes := numBits / 8
		if numBits%8 > 0 {
			numBytes++
		}
		return numBytes
	}
}

func zlibCompress(data []int32) int {
	var out bytes.Buffer
	in := slice.Int32ToByteSlice(data)
	writer, err := zlib.NewWriterLevel(&out, zlib.BestCompression)
	if err != nil {
		panic(err)
	}
	writer.Write(in)
	writer.Close()
	return out.Len()
}

func compressWith(codec encoding.Integer, data []int32) int {
	inlen := len(data)
	out := make([]int32, inlen*2)
	inpos := cursor.New()
	outpos := cursor.New()
	codec.Compress(data, inpos, inlen, out, outpos)
	return outpos.Get() * 4
}

func bp32Compress(data []int32) int {
	return compressWith(composition.New(bp32.New(), variablebyte.New()), data)
}

func deltaBp32Compress(data []int32) int {
	return compressWith(composition.New(deltabp32.New(), deltavb.New()), data)
}

func fastpforCompress(data []int32) int {
	return compressWith(composition.New(fastpfor.New(), variablebyte.New()), data)
}

func deltaFastpforCompress(data []int32) int {
	return compressWith(composition.New(deltafastpfor.New(), deltavb.New()), data)
}

func sizeFunc(data []int32) int {
	return len(data)
}

func sizeInBytesFunc(data []int32) int {
	return 4 * len(data)
}

func parseSizes(str string) []int {
	strSizes := strings.Split(str, ",")
	sizes := make([]int, 0, len(strSizes))
	for _, size := range strSizes {
		if i, err := strconv.Atoi(size); err == nil && i >= 0 {
			sizes = append(sizes, i)
		}
	}
	return sizes
}

func isCardinalityHeaderSizeValid(sz int) bool {
	return sz >= 1 && sz <= 32
}

func makeRandomSlices(sizes []int) [][]int32 {
	slices := make([][]int32, len(sizes))
	for i, size := range sizes {
		slices[i] = slice.SortAscInt32Slice(slice.RandomInt31Slice(size))
	}
	return slices
}

func ratioString(in, out int) string {
	if out == 0 {
		return "inf"
	}
	return fmt.Sprintf("%.2f", float64(in)/float64(out))
}

func outputSizeString(out int) string {
	return strconv.Itoa(out)
}

func makeTable(compressors []compressor, inputData [][]int32, headerRows int, showRatio bool) [][]string {
	tableRows := len(compressors)
	tableColumns := 1 + len(inputData)
	table := make([][]string, tableRows)

	for i := 0; i < tableRows; i++ {
		table[i] = make([]string, tableColumns)
		table[i][0] = compressors[i].name

		for j := 1; j < tableColumns; j++ {
			inputLen := sizeInBytesFunc(inputData[j-1])
			outputLen := compressors[i].compress(inputData[j-1])

			switch {
			case i < headerRows:
				table[i][j] = outputSizeString(outputLen)
			case showRatio:
				table[i][j] = ratioString(inputLen, outputLen)
			default:
				table[i][j] = outputSizeString(outputLen)
			}
		}
	}

	return table
}

func printTable(table [][]string) {
	rows := len(table)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight)

	for i := 0; i < rows; i++ {
		var line strings.Builder
		columns := len(table[i])

		for j := 0; j < columns; j++ {
			line.WriteString(table[i][j])
			line.WriteString("\t")
		}

		fmt.Fprintln(w, line.String())
	}

	w.Flush()
}

func usage() {
	fmt.Println(`
usage: compare-compression-ratio [-help] [-sizes=LIST] [-cardinality-header-size=SIZE] [-ratio]

options:`)
	flag.PrintDefaults()
	fmt.Println()
}

func main() {
	flag.Usage = usage

	helpPtr := flag.Bool("help", false, "print this message")
	sizesPtr := flag.String("sizes", "", "comma-separated sizes, e.g.: 10,100,1000")
	cardHeaderSize := flag.Int("cardinality-header-size", 32, "cardinality header size")
	ratioPtr := flag.Bool("ratio", false, "show compression ratio rather than output size")

	flag.Parse()

	if *helpPtr {
		flag.Usage()
		os.Exit(0)
	}

	if !isCardinalityHeaderSizeValid(*cardHeaderSize) {
		log.Fatalln("invalid -cardinality-header-size value, must be between 1 and 32, both inclusive")
	}

	sizes := parseSizes(*sizesPtr)
	if len(sizes) == 0 {
		log.Fatalln("missing or empty -sizes option")
	}

	randomSlices := makeRandomSlices(sizes)
	compressors := []compressor{
		{name: "integers", compress: sizeFunc},
		{name: "bytes", compress: sizeInBytesFunc},
		{name: "simple", compress: simpleCompress(*cardHeaderSize)},
		{name: "zlib", compress: zlibCompress},
		{name: "bp32", compress: bp32Compress},
		{name: "delta bp32", compress: deltaBp32Compress},
		{name: "fastpfor", compress: fastpforCompress},
		{name: "delta fastpfor", compress: deltaFastpforCompress},
	}

	table := makeTable(compressors, randomSlices, 1, *ratioPtr)

	printTable(table)
}
