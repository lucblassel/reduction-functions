package reductions

import (
	"testing"
)

func BenchmarkGetDistancesIdentity(b *testing.B) {
	sequences, _, _ := ParseFasta("test_data/seqs.fasta")
	for i := 0; i < b.N; i++ {
		GetDistances(sequences, 5, Identity)
	}
}

func BenchmarkGetDistancesHomopolymerCompression(b *testing.B) {
	sequences, _, _ := ParseFasta("test_data/seqs.fasta")
	for i := 0; i < b.N; i++ {
		GetDistances(sequences, 5, HomopolymerCompression)
	}
}

func BenchmarkGetDistancesMultithreadIdentity(b *testing.B) {
	sequences, _, _ := ParseFasta("test_data/seqs.fasta")
	for i := 0; i < b.N; i++ {
		GetDistancesMultiThread(sequences, 5, Identity, 4)
	}
}

func BenchmarkGetDistancesMultithreadHomopolymerCompression(b *testing.B) {
	sequences, _, _ := ParseFasta("test_data/seqs.fasta")
	for i := 0; i < b.N; i++ {
		GetDistancesMultiThread(sequences, 5, HomopolymerCompression, 4)
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
		{Key1: "k1", Key2: "k2", RawDistance: 0.1},
		{Key1: "k3", Key2: "k4", RawDistance: 0.2},
		{Key1: "k5", Key2: "k6", RawDistance: 0.3},
		{Key1: "k7", Key2: "k8", RawDistance: 0.6},
		{Key1: "k9", Key2: "k10", RawDistance: 0.7},
		{Key1: "k11", Key2: "k12", RawDistance: 0.9},
	}
	c, f := MakeSequenceSets(distances, 0.4)

	if !AreDistanceRecordSlicesEqual(c, distances[:3]) {
		t.Errorf("wanted %v got %v", distances[:2], c)
	}
	if !AreDistanceRecordSlicesEqual(f, distances[3:]) {
		t.Errorf("wanted %v got %v", distances[3:], f)
	}
}

func TestGetDistancesMultiThread(t *testing.T) {

	seqs := map[string]string{
		"seq1": "ATTGCATCAT",
		"seq2": "AGTCAGGCAG",
		"seq3": "GTCAGGCATA",
		"seq4": "CGATGGCATA",
	}
	tests := []struct {
		name      string
		wanted    []DistanceRecord
		reduction func(string) string
	}{
		{
			name: "Identity",
			wanted: []DistanceRecord{
				{Key1: "seq2", Key2: "seq3", RawDistance: 0.33333333333333337, ReducedDistance: 0.33333333333333337},
				{Key1: "seq2", Key2: "seq4", RawDistance: 0.8333333333333334, ReducedDistance: 0.8333333333333334},
				{Key1: "seq3", Key2: "seq4", RawDistance: 0.6363636363636364, ReducedDistance: 0.6363636363636364},
				{Key1: "seq1", Key2: "seq2", RawDistance: 0.8181818181818181, ReducedDistance: 0.8181818181818181},
				{Key1: "seq1", Key2: "seq3", RawDistance: 0.7272727272727273, ReducedDistance: 0.7272727272727273},
				{Key1: "seq1", Key2: "seq4", RawDistance: 0.7, ReducedDistance: 0.7},
			},
			reduction: Identity,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			distances := GetDistancesMultiThread(seqs, 3, testCase.reduction, 8)
			if !AreDistanceRecordSlicesEqual(testCase.wanted, distances) {
				t.Errorf("Wanted %v got %v\n(wanted length %v, got %v)", testCase.wanted, distances, len(testCase.wanted), len(distances))
			}
		})
	}
}

func TestGetDistances(t *testing.T) {

	seqs := map[string]string{
		"seq1": "ATTGCATCAT",
		"seq2": "AGTCAGGCAG",
		"seq3": "GTCAGGCATA",
		"seq4": "CGATGGCATA",
	}
	tests := []struct {
		name      string
		wanted    []DistanceRecord
		reduction func(string) string
	}{
		{
			name: "Identity",
			wanted: []DistanceRecord{
				{Key1: "seq2", Key2: "seq3", RawDistance: 0.33333333333333337, ReducedDistance: 0.33333333333333337},
				{Key1: "seq2", Key2: "seq4", RawDistance: 0.8333333333333334, ReducedDistance: 0.8333333333333334},
				{Key1: "seq3", Key2: "seq4", RawDistance: 0.6363636363636364, ReducedDistance: 0.6363636363636364},
				{Key1: "seq1", Key2: "seq2", RawDistance: 0.8181818181818181, ReducedDistance: 0.8181818181818181},
				{Key1: "seq1", Key2: "seq3", RawDistance: 0.7272727272727273, ReducedDistance: 0.7272727272727273},
				{Key1: "seq1", Key2: "seq4", RawDistance: 0.7, ReducedDistance: 0.7},
			},
			reduction: Identity,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			distances := GetDistances(seqs, 3, testCase.reduction)
			if !AreDistanceRecordSlicesEqual(testCase.wanted, distances) {
				t.Errorf("Wanted %v got %v\n(wanted length %v, got %v)", testCase.wanted, distances, len(testCase.wanted), len(distances))
			}
		})
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
			r1:     DistanceRecord{Key1: "k1", Key2: "k2", RawDistance: 1.0},
			r2:     DistanceRecord{Key1: "k1", Key2: "k2", RawDistance: 1.0},
			wanted: true,
		},
		{
			name:   "SameSetSwitchedKeys",
			r1:     DistanceRecord{Key1: "k1", Key2: "k2", RawDistance: 1.0},
			r2:     DistanceRecord{Key1: "k2", Key2: "k1", RawDistance: 1.0},
			wanted: true,
		},
		{
			name:   "DifferentKeys",
			r1:     DistanceRecord{Key1: "k1", Key2: "k2", RawDistance: 1.0},
			r2:     DistanceRecord{Key1: "k3", Key2: "k5", RawDistance: 1.0},
			wanted: false,
		},
		{
			name:   "DifferentDistance",
			r1:     DistanceRecord{Key1: "k1", Key2: "k2", RawDistance: 1.0},
			r2:     DistanceRecord{Key1: "k1", Key2: "k2", RawDistance: 0.5},
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

func TestMakeWFASequenceSets(t *testing.T) {
	distances := []DistanceRecord{
		{Key1: "key1", Key2: "key1_err"},
		{Key1: "key2", Key2: "key2_err"},
		{Key1: "key2", Key2: "key1_err"},
		{Key1: "key4", Key2: "key3_err"},
		{Key1: "key3", Key2: "key4_err"},
	}
	closeSet, farSet := MakeWFASequenceSets(distances)

	if !AreDistanceRecordSlicesEqual(closeSet, distances[:2]) {
		t.Errorf("%v and %v not equal", closeSet, distances[:2])
	}
	if !AreDistanceRecordSlicesEqual(farSet, distances[2:]) {
		t.Errorf("%v and %v not equal", farSet, distances[:2])
	}
}
