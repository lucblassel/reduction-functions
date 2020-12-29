package main

import (
	"testing"
)

func TestReverseComplement(t *testing.T) {
	var tests = map[string]string{
		"ATGC":          "GCAT",
		"AATAAGTCGGCCA": "TGGCCGACTTATT",
		"AAAAAAAATTTTTTTGGGGTTTTAAATTTAGGGATTAAGGATTAGAGCCCATAGAC": "GTCTATGGGCTCTAATCCTTAATCCCTAAATTTAAAACCCCAAAAAAATTTTTTTT",
	}

	for seq, wanted := range tests {
		t.Run(seq, func(t *testing.T) {
			ans, _ := ReverseComplement(seq)
			if ans != wanted {
				t.Errorf("%s (got)\n %s (wanted)", ans, wanted)
			}
		})
	}
}

func TestReverseComplementError(t *testing.T) {

	var tests = [2]string{"", "ATGCH"}

	for _, seq := range tests {
		t.Run(seq, func(t *testing.T) {
			_, err := ReverseComplement(seq)
			if err == nil {
				t.Errorf("Expected error when applying reverse complement to %s", seq)
			}
		})
	}
}

func TestCanonize(t *testing.T) {
	var tests = map[string]string{
		"TCG": "CGA",
		"CGA": "CGA",
		"GAT": "ATC",
		"ATC": "ATC",
		"TCA": "TCA",
		"CAC": "CAC",
		"GTG": "CAC",
	}
	for kmer, wanted := range tests {
		t.Run(kmer, func(t *testing.T) {
			ans, _ := Canonize(kmer)
			if ans != wanted {
				t.Errorf("%s (got)\n %s (wanted)", ans, wanted)
			}
		})
	}
}

func TestKmerize(t *testing.T) {
	var tests = []struct {
		name, seq string
		k         int
		wanted    StringSet
	}{
		{
			name:   "K3",
			seq:    "ATCGATCAC",
			k:      3,
			wanted: MakeSet([]string{"ATC", "CGA", "TCA", "CAC"}),
		},
		{
			name:   "K3OneKmer",
			seq:    "AAAAAAAAA",
			k:      3,
			wanted: MakeSet([]string{"AAA"}),
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ans, _ := Kmerize(testCase.seq, testCase.k)
			if !ans.IsEqual(testCase.wanted) {
				t.Errorf("%s (got)\n%s (wanted)", ans.ToString(), testCase.wanted.ToString())
			}
		})
	}
}

func TestKmerizeErrors(t *testing.T) {
	var tests = []struct {
		name, seq string
		k         int
	}{
		{
			name: "KTooSmall",
			seq:  "ATGCTGAC",
			k:    1,
		},
		{
			name: "KNegative",
			seq:  "ATGCTGAC",
			k:    -1,
		},
		{
			name: "KTooBig",
			seq:  "ATGCTGAC",
			k:    10,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := Kmerize(testCase.seq, testCase.k)
			if err == nil {
				t.Errorf("Was expecting error when kmerizing %s with k=%d", testCase.seq, testCase.k)
			}
		})
	}
}
