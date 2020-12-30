package main

import (
	"testing"
)

func BenchmarkGetDistances(b *testing.B) {
	sequences, _ := ParseFasta("test.fasta")
	for i := 0; i < b.N; i++ {
		GetDistances(sequences, 5)
	}
}

func TestJaccardSimilarity(t *testing.T) {
	tests := []struct {
		name       string
		set1, set2 StringSet
		wanted     float64
	}{
		{
			name:   "EmptySets",
			set1:   StringSet{},
			set2:   StringSet{},
			wanted: 0.0,
		},
		{
			name:   "SameSet",
			set1:   MakeSet([]string{"ATG", "CTT", "GTA"}),
			set2:   MakeSet([]string{"ATG", "CTT", "GTA"}),
			wanted: 1.0,
		},
		{
			name:   "DifferentSets",
			set1:   MakeSet([]string{"ATG", "CTT", "GTA"}),
			set2:   MakeSet([]string{"ATG", "CTT", "GTT"}),
			wanted: 0.5,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			if ans := JaccardSimilarity(testCase.set1, testCase.set2); ans != testCase.wanted {
				t.Errorf("Wanted %v, got %v instead", testCase.wanted, ans)
			}
		})
	}
}

func TestKmerizedJaccardDistance(t *testing.T) {
	tests := []struct {
		name, seq1, seq2 string
		k                int
		wanted           float64
	}{
		{
			name:   "SameString",
			seq1:   "ATGCATGCATCAGCATTGCA",
			seq2:   "ATGCATGCATCAGCATTGCA",
			k:      3,
			wanted: 0.0,
		},
		{
			name:   "DifferentStrings",
			seq1:   "ATGCATGCATCAGCATTGCA",
			seq2:   "ATGCAGGCATCAGGCATCAG",
			k:      3,
			wanted: 0.5,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			if ans, _ := KmerizedJaccardDistance(testCase.seq1, testCase.seq2, testCase.k); ans != testCase.wanted {
				t.Errorf("Wanted %v, got %v", testCase.wanted, ans)
			}
		})
	}
}

func TestMakeSequenceSets(t *testing.T) {
	distances := []DistanceRecord{
		{key1: "k1", key2: "k2", distance: 0.1},
		{key1: "k3", key2: "k4", distance: 0.2},
		{key1: "k5", key2: "k6", distance: 0.3},
		{key1: "k7", key2: "k8", distance: 0.6},
		{key1: "k9", key2: "k10", distance: 0.7},
		{key1: "k11", key2: "k12", distance: 0.9},
	}
	c, f := MakeSequenceSets(distances, 0.4)

	if !testEq(c, distances[:3]) {
		t.Errorf("wanted %v got %v", distances[:2], c)
	}
	if !testEq(f, distances[3:]) {
		t.Errorf("wanted %v got %v", distances[3:], f)
	}
}

func TestGetDistances(t *testing.T) {

	seqs := map[string]string{
		"seq1": "ATTGCATCAT",
		"seq2": "AGTCAGGCAG",
		"seq3": "GTCAGGCATA",
		"seq4": "CGATGGCATA",
	}
	wanted := []DistanceRecord{
		{"seq2", "seq3", 0.33333333333333337},
		{"seq2", "seq4", 0.8333333333333334},
		{"seq3", "seq4", 0.6363636363636364},
		{"seq1", "seq2", 0.8181818181818181},
		{"seq1", "seq3", 0.7272727272727273},
		{"seq1", "seq4", 0.7},
	}
	distances := GetDistances(seqs, 3)

	if !AreDistanceRecordSlicesEqual(wanted, distances) {
		t.Errorf("Wanted %v got %v", wanted, distances)
	}
}

func TestDistanceRecord_IsEqual(t *testing.T) {
	tests := []struct {
		name   string
		r1, r2 DistanceRecord
		wanted bool
	}{
		{
			name:   "SameSet",
			r1:     DistanceRecord{"k1", "k2", 1.0},
			r2:     DistanceRecord{"k1", "k2", 1.0},
			wanted: true,
		},
		{
			name:   "SameSetSwitchedKeys",
			r1:     DistanceRecord{"k1", "k2", 1.0},
			r2:     DistanceRecord{"k2", "k1", 1.0},
			wanted: true,
		},
		{
			name:   "DifferentKeys",
			r1:     DistanceRecord{"k1", "k2", 1.0},
			r2:     DistanceRecord{"k3", "k5", 1.0},
			wanted: false,
		},
		{
			name:   "DifferentDistance",
			r1:     DistanceRecord{"k1", "k2", 1.0},
			r2:     DistanceRecord{"k1", "k2", 0.5},
			wanted: false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			if ans := testCase.r1.IsEqual(testCase.r2); ans != testCase.wanted {
				t.Errorf("got equal=%v for records %v and %v", ans, testCase.r1, testCase.r2)
			}
		})
	}
}
