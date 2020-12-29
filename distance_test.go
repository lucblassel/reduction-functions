package main

import "testing"

func BenchmarkGetDistances(b *testing.B) {
	sequences, _ := ParseFasta("test.fasta")
	for i := 0; i < b.N; i++ {
		GetDistances(sequences, 5)
	}
}
