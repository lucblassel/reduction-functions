package main

import (
	"math"
	"strings"
	"sync"
)

// HomopolymerCompression returns the homopolymer compressed version of a read
func HomopolymerCompression(read string) string {
	var lastChar rune
	var builder strings.Builder

	for _, char := range read {
		if char == lastChar {
			continue
		}
		builder.WriteRune(char)
		lastChar = char
	}

	return builder.String()
}

// Identity returns the input read unchanged
func Identity(read string) string {
	return read
}

func CloseDistanceRoutine(r DistanceRecord, seq1, seq2 string, k int, queue chan DistanceRecord) {
	dist, _ := KmerizedJaccardDistance(seq1, seq2, k)
	queue <- DistanceRecord{
		key1:     r.key1,
		key2:     r.key2,
		distance: dist,
	}
}

func FarDistanceRoutine(r DistanceRecord, seq1, seq2 string, k int, queue chan DistanceRecord) {
	dist, _ := KmerizedJaccardDistance(seq1, seq2, k)
	queue <- DistanceRecord{
		key1:     r.key1,
		key2:     r.key2,
		distance: dist / r.distance,
	}
}

// ComputePhiTerm computes either one of the terms of the objective function depending
func ComputePhiTerm(sequences map[string]string, set []DistanceRecord, k int, routine func(DistanceRecord, string, string, int, chan DistanceRecord)) []DistanceRecord {
	queue := make(chan DistanceRecord, len(sequences))
	var wg sync.WaitGroup
	distances := []DistanceRecord{}

	for _, record := range set {
		seq1, seq2 := sequences[record.key1], sequences[record.key2]
		wg.Add(1)
		go routine(record, seq1, seq2, k, queue)
	}
	go func() {
		for record := range queue {
			distances = append(distances, record)
			wg.Done()
		}
		close(queue)
	}()

	wg.Wait()

	return distances
}

// SumDistance returns the total distance of a slice of DistanceRecords
func SumDistance(records []DistanceRecord) float64 {
	var totalDistance float64
	for _, record := range records {
		totalDistance += record.distance
	}
	return totalDistance
}

// PhiRecord
type PhiRecord struct {
	phi, C, F, mu float64
}

// ObjectivePhi is the objective function to evaluate the reduction function
func ObjectivePhi(sequences map[string]string, closeSet, farSet []DistanceRecord, k int) PhiRecord {
	closeDists := ComputePhiTerm(sequences, closeSet, k, CloseDistanceRoutine)
	farRatios := ComputePhiTerm(sequences, farSet, k, FarDistanceRoutine)
	closeDist := SumDistance(closeDists) / float64(len(closeDists))
	mu := SumDistance(farRatios) / float64(len(farRatios))
	farVariance := 0.0

	for _, record := range farRatios {
		farVariance += math.Pow(record.distance-mu, 2)
	}

	farVariance /= float64(len(farRatios))
	return PhiRecord{
		phi: closeDist + farVariance/mu,
		C:   closeDist,
		F:   farVariance,
		mu:  mu,
	}
}
