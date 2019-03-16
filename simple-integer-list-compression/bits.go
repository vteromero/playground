package simple

import "math/bits"

func bitsLen(x uint32) int {
	if x == 0 {
		return 1
	}
	return bits.Len32(x)
}

func sizeInBytes(bits int) int {
	sz := bits / 8
	if bits%8 > 0 {
		sz++
	}
	return sz
}
