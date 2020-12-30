package main

import (
	"testing"
)

func TestIdentity(t *testing.T) {
	tests := []struct {
		name, read, wanted string
	}{
		{name: "EmptyString", read: "", wanted: ""},
		{name: "NonEmptyString", read: "AATGCCAGTCA", wanted: "AATGCCAGTCA"},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ans := Identity(testCase.read)
			if ans != testCase.wanted {
				t.Errorf("Wanted %s, but got %s", testCase.wanted, ans)
			}
		})
	}
}

func TestHomopolymerCompression(t *testing.T) {
	tests := []struct {
		name, read, wanted string
	}{
		{name: "EmptyString", read: "", wanted: ""},
		{name: "NonEmptyString", read: "AATGCCAGTCA", wanted: "ATGCAGTCA"},
		{name: "NoRepeats", read: "ATGCATGCATGTGCA", wanted: "ATGCATGCATGTGCA"},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ans := HomopolymerCompression(testCase.read)
			if ans != testCase.wanted {
				t.Errorf("Wanted %s, but got %s", testCase.wanted, ans)
			}
		})
	}
}

func TestSumDistance(t *testing.T) {
	tests := []struct {
		name    string
		records []DistanceRecord
		wanted  float64
	}{
		{
			name:    "EmptySet",
			records: []DistanceRecord{},
			wanted:  0,
		},
		{
			name: "WithRecords",
			records: []DistanceRecord{
				{key1: "k1", key2: "k2", distance: 1.0},
				{key1: "k1", key2: "k2", distance: 2.0},
				{key1: "k1", key2: "k2", distance: 3.0},
			},
			wanted: 6.0,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			if ans := SumDistance(testCase.records); ans != testCase.wanted {
				t.Errorf("Wanted %v but got %v", testCase.wanted, ans)
			}
		})
	}
}

func TestComputePhiTerm(t *testing.T) {
	k := 3

	seqs := map[string]string{
		"seq1": "ATTTGAGCA",
		"seq2": "ATTTGAGTC",
		"seq3": "ATTTGCAGT",
		"seq4": "ATTTGCCAG",
	}
	set := []DistanceRecord{
		{key1: "seq1", key2: "seq2", distance: 0.8},
		{key1: "seq3", key2: "seq4", distance: 0.4},
	}

	tests := []struct {
		name    string
		wanted  []DistanceRecord
		routine func(DistanceRecord, string, string, int, chan DistanceRecord)
	}{
		{
			name: "Close",
			wanted: []DistanceRecord{
				{"seq1", "seq2", 0.4444444444444444},
				{"seq3", "seq4", 0.375},
			},
			routine: CloseDistanceRoutine,
		},
		{
			name: "Close",
			wanted: []DistanceRecord{
				{"seq1", "seq2", 0.5555555555555555},
				{"seq3", "seq4", 0.9375},
			},
			routine: FarDistanceRoutine,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {

			distances := ComputePhiTerm(seqs, set, k, testCase.routine)

			if !AreDistanceRecordSlicesEqual(testCase.wanted, distances) {
				t.Errorf("wanted: %v\ngot: %v", testCase.wanted, distances)
			}

		})
	}
}
