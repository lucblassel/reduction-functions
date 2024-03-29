package reductions

import (
	"fmt"
	"gonum.org/v1/gonum/stat/combin"
	"sort"
	"strings"
	"sync"
)

// DistanceRecord keeps the distance between 2 sequences with the given keys
type DistanceRecord struct {
	Key1, Key2                   string
	RawDistance, ReducedDistance float64
}

// String implements the Stringer interface for DistanceRecord structs
func (record DistanceRecord) String() string {
	return fmt.Sprintf(
		"{%v,%v: %0.3f, %0.3f}",
		record.Key1, record.Key2, record.RawDistance, record.ReducedDistance,
	)
}

// sortRecordKeys puts the keys in lexicographical order
func sortRecordKeys(record DistanceRecord) DistanceRecord {
	if record.Key1 < record.Key2 {
		return record
	}
	return DistanceRecord{
		Key1:            record.Key2,
		Key2:            record.Key1,
		RawDistance:     record.RawDistance,
		ReducedDistance: record.ReducedDistance,
	}
}

// AreDistanceRecordSlicesEqual checks that 2 slices contain the same DistanceRecords
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
		return r1.RawDistance < r2.RawDistance
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

// IsEqual checks if 2 distance records are equal
func (record DistanceRecord) IsEqual(other DistanceRecord) bool {
	if record.RawDistance != other.RawDistance ||
		record.ReducedDistance != other.ReducedDistance {
		return false
	}
	return (record.Key1 == other.Key1 && record.Key2 == other.Key2) ||
		(record.Key1 == other.Key2 && record.Key2 == other.Key1)
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
func GetDistances(seqRecords map[string]string, k int, reduction func(string) string) []DistanceRecord {
	seqKeys := make([]string, 0, len(seqRecords))
	for key := range seqRecords {
		seqKeys = append(seqKeys, key)
	}

	nRecords := combin.Binomial(len(seqKeys), 2)
	distances := make([]DistanceRecord, 0, nRecords)

	queue := make(chan DistanceRecord, 1)
	var wg sync.WaitGroup

	for i, key1 := range seqKeys {
		for _, key2 := range seqKeys[i+1:] {
			wg.Add(1)
			go func(key1, key2 string) {
				rawDist, _ := KmerizedJaccardDistance(seqRecords[key1], seqRecords[key2], k)
				redDist, _ := KmerizedJaccardDistance(reduction(seqRecords[key1]), reduction(seqRecords[key2]), k)
				queue <- DistanceRecord{
					Key1:            key1,
					Key2:            key2,
					RawDistance:     rawDist,
					ReducedDistance: redDist,
				}
			}(key1, key2)
		}
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

func distanceWorker(id int, jobs <-chan []string, results chan<- DistanceRecord, seqRecords map[string]string, k int, reduction func(string) string) {
	for job := range jobs {
		rawDist, _ := KmerizedJaccardDistance(seqRecords[job[0]], seqRecords[job[1]], k)
		redDist, _ := KmerizedJaccardDistance(reduction(seqRecords[job[0]]), reduction(seqRecords[job[1]]), k)
		results <- DistanceRecord{
			Key1:            job[0],
			Key2:            job[1],
			RawDistance:     rawDist,
			ReducedDistance: redDist,
		}
	}
}

// GetDistancesMultiThread computes distances between all pairs of strings in a list
func GetDistancesMultiThread(seqRecords map[string]string, k int, reduction func(string) string, threads int) []DistanceRecord {
	seqKeys := make([]string, 0, len(seqRecords))
	for key := range seqRecords {
		seqKeys = append(seqKeys, key)
	}

	nRecords := combin.Binomial(len(seqKeys), 2)
	distances := make([]DistanceRecord, 0, nRecords)

	if nRecords < threads {
		threads = nRecords
	}

	results := make(chan DistanceRecord, nRecords)
	keyChannel := make(chan []string, nRecords)

	for w := 1; w <= nRecords; w++ {
		go distanceWorker(w, keyChannel, results, seqRecords, k, reduction)
	}

	for i, key1 := range seqKeys {
		for _, key2 := range seqKeys[i+1:] {
			keyChannel <- []string{key1, key2}
		}
	}

	close(keyChannel)

	for r := 1; r <= nRecords; r++ {
		distances = append(distances, <-results)
	}

	return distances
}

// MakeSequenceSets separates a set of sequences distances into the set of close sequences (dist <= radius)
// and the set of far sequences (dist > radius) and returns them both
func MakeSequenceSets(distances []DistanceRecord, radius float64) ([]DistanceRecord, []DistanceRecord) {
	closeSet := []DistanceRecord{}
	farSet := []DistanceRecord{}

	for _, record := range distances {
		if record.RawDistance <= radius {
			closeSet = append(closeSet, record)
		} else {
			farSet = append(farSet, record)
		}
	}
	return closeSet, farSet
}

// MakeWFASequenceSets generates a "close" and a "far" sequence set from the WFA generate_dataset sequences
func MakeWFASequenceSets(distances []DistanceRecord) ([]DistanceRecord, []DistanceRecord) {
	closeSet := []DistanceRecord{}
	farSet := []DistanceRecord{}

	for _, record := range distances {
		if strings.Replace(record.Key1, "_err", "", -1) ==
			strings.Replace(record.Key2, "_err", "", -1) {
			closeSet = append(closeSet, record)
		} else {
			farSet = append(farSet, record)
		}
	}
	return closeSet, farSet
}
