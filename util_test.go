package main

import (
	"testing"
)

//
// go test -bench=.
//
func BenchmarkCA(b *testing.B) {

	ca := &cntArr{}

	for n := 0; n < b.N; n++ {
		incCounter(ca)
		ca.a++
	}

}

//
//
//
func BenchmarkInt2ByteArr(b *testing.B) {

	//	t := NewIntToByteArr()

	for n := 0; n < b.N; n++ {
		// t.in = time.Now().UnixNano()
		// intToByteArr(t)
		b := make([]byte, 256, 512)
		b = b
		b = nil
	}

}
