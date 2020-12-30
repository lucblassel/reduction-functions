package main

import (
	comb "github.com/mxschmitt/golang-combinations"
	"gonum.org/v1/gonum/stat/combin"
	"sort"
	"sync"
)

// DistanceRecord keeps the distance between 2 sequences with the given keys
type DistanceRecord struct {
	key1, key2 string
	distance   float64
}

// tests if 2 slices of distance records are equal
func testEq(a, b []DistanceRecord) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func sortRecordKeys(record DistanceRecord) DistanceRecord {
	if record.key1 < record.key2 {
		return record
	}
	return DistanceRecord{record.key2, record.key1, record.distance}
}

func AreDistanceRecordSlicesEqual(a, b []DistanceRecord) bool {
	aReord := []DistanceRecord{}
	bReord := []DistanceRecord{}

	for _, record := range a {
		aReord = append(aReord, sortRecordKeys(record))
	}
	for _, record := range b {
		bReord = append(bReord, sortRecordKeys(record))
	}

	comp := func(i, j int, recs []DistanceRecord) bool {
		r1, r2 := recs[i], recs[j]
		return r1.distance < r2.distance
	}

	sort.SliceStable(aReord, func(i, j int) bool { return comp(i, j, aReord) })
	sort.SliceStable(bReord, func(i, j int) bool { return comp(i, j, bReord) })

	if len(aReord) != len(bReord) {
		return false
	}

	for i := range aReord {
		if aReord[i] != bReord[i] {
			return false
		}
	}
	return true
}

// Check if 2 distance records are equal
func (r1 DistanceRecord) IsEqual(r2 DistanceRecord) bool {
	if r1.distance != r2.distance {
		return false
	}
	return (r1.key1 == r2.key1 && r1.key2 == r2.key2) ||
		(r1.key1 == r2.key2 && r1.key2 == r2.key1)
}

// JaccardSimilarity returns the Jaccard index between two
// sets of strings.
func JaccardSimilarity(set1, set2 StringSet) float64 {
	if len(set1) == 0 || len(set2) == 0 {
		return 0.0
	}
	interLen := len(set1.Intersection(set2))
	unionLen := len(set1.Union(set2))

	return float64(interLen) / float64(unionLen)
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

// MakeSequenceSets separates a set of sequences distances into the set of close sequences (dist <= radius)
// and the set of far sequences (dist > radius) and returns them both
func MakeSequenceSets(distances []DistanceRecord, radius float64) ([]DistanceRecord, []DistanceRecord) {
	closeSet := []DistanceRecord{}
	farSet := []DistanceRecord{}

	for _, record := range distances {
		if record.distance <= radius {
			closeSet = append(closeSet, record)
		} else {
			farSet = append(farSet, record)
		}
	}
	return closeSet, farSet
}
