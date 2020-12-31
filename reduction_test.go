package reductions

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

func TestSumArray(t *testing.T) {
	tests := []struct {
		name   string
		slice  []float64
		wanted float64
	}{
		{
			name:   "Empty",
			slice:  []float64{},
			wanted: 0,
		},
		{
			name:   "Full",
			slice:  []float64{1, 2, 3},
			wanted: 6,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			if ans := SumArray(testCase.slice); ans != testCase.wanted {
				t.Errorf("wanted %v, got %v", testCase.wanted, ans)
			}
		})
	}
}

func areSlicesEqual(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != a[i] {
			return false
		}
	}
	return true
}

func TestComputeCloseTerms(t *testing.T) {
	records := []DistanceRecord{
		{Key1: "k1", Key2: "k2", ReducedDistance: 1},
		{Key1: "k1", Key2: "k2", ReducedDistance: 2},
		{Key1: "k1", Key2: "k2", ReducedDistance: 3},
	}
	wanted := []float64{1, 2, 3}
	ans := ComputeCloseTerms(records)
	if !areSlicesEqual(wanted, ans) {
		t.Errorf("wanted %v got %v", wanted, ans)
	}
}

func TestComputeFarRatios(t *testing.T) {
	records := []DistanceRecord{
		{Key1: "k1", Key2: "k2", ReducedDistance: 1, RawDistance: 2},
		{Key1: "k1", Key2: "k2", ReducedDistance: 2, RawDistance: 2},
		{Key1: "k1", Key2: "k2", ReducedDistance: 3, RawDistance: 2},
	}
	wanted := []float64{0.5, 1, 1.5}

	ans := ComputeFarRatios(records)
	if !areSlicesEqual(wanted, ans) {
		t.Errorf("wanted %v got %v", wanted, ans)
	}
}

func TestComputeFarTerms(t *testing.T) {
	records := []DistanceRecord{
		{Key1: "k1", Key2: "k2", ReducedDistance: 1, RawDistance: 2},
		{Key1: "k1", Key2: "k2", ReducedDistance: 2, RawDistance: 2},
		{Key1: "k1", Key2: "k2", ReducedDistance: 3, RawDistance: 2},
	}

	wanted := []float64{0, 1, 4}

	farRatios := ComputeFarRatios(records)
	ans := ComputeFarTerms(farRatios, 1)
	if !areSlicesEqual(wanted, ans) {
		t.Errorf("wanted %v got %v", wanted, ans)
	}
}
