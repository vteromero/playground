package slice

import (
	"encoding/binary"
	"math/rand"
	"sort"
)

func RandomInt31Slice(n int) []int32 {
	s := make([]int32, n)
	for i := 0; i < n; i++ {
		s[i] = rand.Int31()
	}
	return s
}

func RandomUint32Slice(n int) []uint32 {
	s := make([]uint32, n)
	for i := 0; i < n; i++ {
		s[i] = rand.Uint32()
	}
	return s
}

func Int32ToUint32Slice(s []int32) []uint32 {
	n := len(s)
	out := make([]uint32, n)
	for i := 0; i < n; i++ {
		out[i] = uint32(s[i])
	}
	return out
}

func Uint32ToInt32Slice(s []uint32) []int32 {
	n := len(s)
	out := make([]int32, n)
	for i := 0; i < n; i++ {
		out[i] = int32(s[i])
	}
	return out
}

func Int32ToByteSlice(s []int32) []byte {
	n := len(s)
	out := make([]byte, n*4)
	for i := 0; i < n; i++ {
		binary.LittleEndian.PutUint32(out[i*4:], uint32(s[i]))
	}
	return out
}

func Uint32ToByteSlice(s []uint32) []byte {
	n := len(s)
	out := make([]byte, n*4)
	for i := 0; i < n; i++ {
		binary.LittleEndian.PutUint32(out[i*4:], s[i])
	}
	return out
}

func SortAscUint32Slice(s []uint32) []uint32 {
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})
	return s
}

func SortDescUint32Slice(s []uint32) []uint32 {
	sort.Slice(s, func(i, j int) bool {
		return s[i] > s[j]
	})
	return s
}

func SortAscInt32Slice(s []int32) []int32 {
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})
	return s
}

func SortDescInt32Slice(s []int32) []int32 {
	sort.Slice(s, func(i, j int) bool {
		return s[i] > s[j]
	})
	return s
}
