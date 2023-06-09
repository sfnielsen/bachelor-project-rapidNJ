package main

import (
	"bytes"
	"encoding/binary"
)

const digit = 8

func radixsort(data []Tuple) {
	buf := bytes.NewBuffer(nil)
	ds := make([][]byte, len(data))

	for i, e := range data {
		binary.Write(buf, binary.LittleEndian, e.value)
		b := make([]byte, digit)
		buf.Read(b)
		ds[i] = b
	}

	bucket := make([][][]byte, 256)
	tuples := make([][]Tuple, 256)

	for i := 0; i < digit; i++ {
		for asdf, b := range ds {
			bucket[b[i]] = append(bucket[b[i]], b)
			tuples[b[i]] = append(tuples[b[i]], data[asdf])
		}
		j := 0
		for k, bs := range bucket {
			copy(ds[j:], bs)
			copy(data[j:], tuples[k])

			j += len(bs)

			bucket[k] = bs[:0]
			tuples[k] = tuples[k][:0]
		}
	}
	/*
		var w float64

		for i, b := range ds {
			buf.Write(b)
			binary.Read(buf, binary.LittleEndian, &w)
			data[i] = placementMap[w]
		}
	*/
}
