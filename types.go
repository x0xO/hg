package hg

import (
	"bufio"
	"os"
)

type (
	hfiter struct {
		scanner *bufio.Scanner
		hfile   *HFile
	}

	HFile struct {
		err    error
		file   *os.File
		hfiter *hfiter
		name   HString
	}

	HDir struct {
		err  error
		path HString
	}

	HString                   string
	HInt                      int
	HFloat                    float64
	HBytes                    []byte
	HSlice[T any]             []T
	HMap[K comparable, V any] map[K]V
	HSet[T comparable]        map[T]struct{}

	hMapPair[K comparable, V any] struct {
		Key   K
		Value V
	}

	HMapOrd[K comparable, V any] HSlice[hMapPair[K, V]]
)
