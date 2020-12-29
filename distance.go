package main

import (
	comb "github.com/mxschmitt/golang-combinations"
	"gonum.org/v1/gonum/stat/combin"
	"sync"
)

// DistanceRecord keeps the distance between 2 sequences with the given keys
type DistanceRecord struct {
	key1, key2 string
	distance   float64
}

// JaccardSimilarity returns the Jaccard index between two
// sets of strings.
func JaccardSimilarity(set1, set2 StringSet) float64 {
	interLen := len(set1.Intersection(set2))
	unionLen := len(set1.Union(set2))

	return float64(interLen / unionLen)
}

// KmerizedJaccardDistance returns the Jaccard distance between the kmers of 2 sequences
// for a given k.
func KmerizedJaccardDistance(seq1, seq2 string, k int) (float64, error) {
	kmers1, err := Kmerize(seq1, k)
	if err != nil {
		return 0, err
	}
	kmers2, err := Kmerize(seq2, k)
	if err != nil {
		return 0, err
	}

	return 1. - JaccardSimilarity(kmers1, kmers2), nil
}

// GetDistances computes distances between all pairs of strings in a list
func GetDistances(seqRecords map[string]string, k int) []DistanceRecord {
	seqKeys := make([]string, 0, len(seqRecords))
	for key := range seqRecords {
		seqKeys = append(seqKeys, key)
	}

	distances := make([]DistanceRecord, 0, combin.Binomial(len(seqKeys), 2))

	queue := make(chan DistanceRecord, 1)
	var wg sync.WaitGroup

	for _, elements := range comb.Combinations(seqKeys, 2) {
		wg.Add(1)
		go func(key1, key2 string) {
			dist, _ := KmerizedJaccardDistance(seqRecords[key1], seqRecords[key2], k)
			queue <- DistanceRecord{
				key1:     key1,
				key2:     key2,
				distance: dist,
			}
		}(elements[0], elements[1])
	}

	go func() {
		for record := range queue {
			distances = append(distances, record)
			wg.Done()
		}
	}()

	wg.Wait()

	return distances
}
