package reductions

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"sort"
	"testing"
	"github.com/hillbig/rsdic"
)

func isSurjection(nEnd int, mapping []int) bool {
	outputs := make(map[int]bool)
	for _, elem := range mapping {
		outputs[elem] = true
	}
	for i := 0; i < nEnd; i++ {
		if !outputs[i] {
			return false
		}
	}
	return true
}

func TestRandomMapping(t *testing.T) {
	for i := 0; i < 100; i++ {
		start, end := rand.Intn(99)+2, rand.Intn(99)+2
		mapping1, mapping2 := RandomMapping(start, end), RandomMapping(start, end)
		equal := true
		for j, elem1 := range mapping1 {
			if elem1 != mapping2[j] {
				equal = false
				break
			}
		}
		if equal {
			t.Errorf("Random mappings should be different:\n%v %v", mapping1, mapping2)
		}
	}
}

func TestSurjection(t *testing.T) {
	for i := 0; i < 100; i++ {
		start, end := rand.Intn(99)+1, rand.Intn(99)+1
		mapping, err := Surjection(start, end)
		if end > start {
			if err == nil {
				t.Errorf("Surjection of %v to %v is impossible, an error should have been returned", start, end)
			}
			continue
		} else {
			if !isSurjection(end, mapping) {
				t.Errorf("mapping %v is not a surjectiong of %v to %v", mapping, start, end)
			}
		}
	}
}

func TestCountSurjections(t *testing.T) {
	cases := []struct {
		nStart, nEnd, wanted int
	}{
		{nStart: 5, nEnd: 5, wanted: 120},
		{nStart: 6, nEnd: 5, wanted: 1800},
		{nStart: 6, nEnd: 6, wanted: 720},
		{nStart: 7, nEnd: 5, wanted: 16800},
		{nStart: 7, nEnd: 6, wanted: 15120},
		{nStart: 7, nEnd: 7, wanted: 5040},
		{nStart: 8, nEnd: 5, wanted: 126000},
		{nStart: 8, nEnd: 6, wanted: 191520},
		{nStart: 8, nEnd: 7, wanted: 141120},
		{nStart: 8, nEnd: 8, wanted: 40320},
		{nStart: 9, nEnd: 5, wanted: 834120},
		{nStart: 9, nEnd: 6, wanted: 1905120},
		{nStart: 9, nEnd: 7, wanted: 2328480},
		{nStart: 9, nEnd: 8, wanted: 1451520},
		{nStart: 9, nEnd: 9, wanted: 362880},
	}
	for _, testCase := range cases {
		t.Run(fmt.Sprintf("%v->%v", testCase.nStart, testCase.nEnd), func(t *testing.T) {
			if count := CountSurjections(testCase.nStart, testCase.nEnd); count != testCase.wanted {
				t.Errorf("Got %v expected %v for (%v->%v)", count, testCase.wanted, testCase.nStart, testCase.nEnd)
			}
		})
	}
}

func TestGetTuples(t *testing.T) {
	alphabet := "ATGC"
	cases := []struct {
		name, universe string
		n              int
		wanted         []string
	}{
		{name: "simplest", universe: alphabet, n: 1, wanted: []string{"A", "C", "G", "T"}},
		{
			name:     "length2",
			universe: alphabet,
			n:        2,
			wanted: []string{
				"AA", "AC", "AG", "AT",
				"CA", "CC", "CG", "CT",
				"GA", "GC", "GG", "GT",
				"TA", "TC", "TG", "TT",
			},
		},
		{
			name:     "length3",
			universe: alphabet,
			n:        3,
			wanted: []string{
				"AAA", "AAC", "AAG", "AAT",
				"ACA", "ACC", "ACG", "ACT",
				"AGA", "AGC", "AGG", "AGT",
				"ATA", "ATC", "ATG", "ATT",
				"CAA", "CAC", "CAG", "CAT",
				"CCA", "CCC", "CCG", "CCT",
				"CGA", "CGC", "CGG", "CGT",
				"CTA", "CTC", "CTG", "CTT",
				"GAA", "GAC", "GAG", "GAT",
				"GCA", "GCC", "GCG", "GCT",
				"GGA", "GGC", "GGG", "GGT",
				"GTA", "GTC", "GTG", "GTT",
				"TAA", "TAC", "TAG", "TAT",
				"TCA", "TCC", "TCG", "TCT",
				"TGA", "TGC", "TGG", "TGT",
				"TTA", "TTC", "TTG", "TTT",
			},
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			perms := GetTuples(testCase.n, 0, testCase.universe, []string{})
			sort.Strings(perms)
			if len(perms) != len(testCase.wanted) {
				t.Errorf("Mismatched lengths: got %v, wanted %v", perms, testCase.wanted)
			}
			for i, perm := range perms {
				if perm != testCase.wanted[i] {
					t.Errorf("Not the same permutations. Got %v wanted %v", perms, testCase.wanted)
				}
			}

		})
	}
}

func TestGetRandomReduction(t *testing.T) {
	inAlph, outAlph := "ATGC", "ATGC."
	cases := []struct {
		name, inAlph, outAlph string
		inSize, outSize       int
	}{
		{
			name: "2->1", inAlph: inAlph, outAlph: outAlph, inSize: 2, outSize: 1,
		},
		{
			name: "3->2", inAlph: inAlph, outAlph: outAlph, inSize: 3, outSize: 2,
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			function, err := GetRandomReduction(testCase.inAlph, testCase.outAlph, testCase.inSize, testCase.outSize)
			if err != nil {
				t.Errorf("error while getting reduction function: %v", err)
			}
			nCases := math.Pow(float64(len(testCase.inAlph)), float64(testCase.inSize))
			if len(function) != int(nCases) {
				t.Errorf("Mismatched lengths expected %v got %v", nCases, len(function))
			}
			outputs := map[string]bool{}
			for _, v := range function {
				outputs[v] = true
			}
			nOuts := math.Pow(float64(len(testCase.outAlph)), float64(testCase.outSize))
			if len(outputs) != int(nOuts) {
				t.Errorf("Not a surjection, got %v outupts wanted %v", nOuts, len(outputs))
			}
		})
	}
}

func TestMakeReductionFunction(t *testing.T) {
	cases := []struct {
		name      string
		reduction map[string]string
		input     string
		wanted    string
	}{
		{
			name: "empty",
			reduction: map[string]string{
				"AA": ".", "AC": ".", "AG": ".", "AT": ".",
				"CA": ".", "CC": ".", "CG": ".", "CT": ".",
				"GA": ".", "GC": ".", "GG": ".", "GT": ".",
				"TA": ".", "TC": ".", "TG": ".", "TT": ".",
			},
			input:  "GATGCCAGTGCAGCATGTCAGGACGAC",
			wanted: "G",
		},
		{
			name: "homopolymer",
			reduction: map[string]string{
				"AA": ".", "AC": "C", "AG": "G", "AT": "T",
				"CA": "A", "CC": ".", "CG": "G", "CT": "T",
				"GA": "A", "GC": "C", "GG": ".", "GT": "T",
				"TA": "A", "TC": "C", "TG": "G", "TT": ".",
			},
			input:  "AAAAATTTTTTTGGGGGGCCCCCAAAAAATTTTTTGGGGGGGCCCCCCC",
			wanted: "ATGCATGC",
		},
		{
			name: "first letter",
			reduction: map[string]string{
				"AA": "A", "AC": "A", "AG": "A", "AT": "A",
				"CA": "C", "CC": "C", "CG": "C", "CT": "C",
				"GA": "G", "GC": "G", "GG": "G", "GT": "G",
				"TA": "T", "TC": "T", "TG": "T", "TT": "T",
			},
			input:  "AGTCAGA",
			wanted: "AAGTCAG",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			reduction := MakeReductionFunction(testCase.reduction)
			output := reduction(testCase.input)
			if output != testCase.wanted {
				t.Errorf("Error in making reduction function.\nWanted %v\nGot: %v", testCase.wanted, output)
			}
		})
	}
}

func TestMakeReductionFunctionKeepOffsets(t *testing.T) {
	cases := []struct {
		name      string
		reduction map[string]string
		input     string
		wanted    string
		encoded   string
	}{
		{
			name: "empty",
			reduction: map[string]string{
				"AA": ".", "AC": ".", "AG": ".", "AT": ".",
				"CA": ".", "CC": ".", "CG": ".", "CT": ".",
				"GA": ".", "GC": ".", "GG": ".", "GT": ".",
				"TA": ".", "TC": ".", "TG": ".", "TT": ".",
			},
			input:   "GATGC",
			wanted:  "G",
			encoded: "M1D4",
		},
		{
			name: "homopolymer",
			reduction: map[string]string{
				"AA": ".", "AC": "C", "AG": "G", "AT": "T",
				"CA": "A", "CC": ".", "CG": "G", "CT": "T",
				"GA": "A", "GC": "C", "GG": ".", "GT": "T",
				"TA": "A", "TC": "C", "TG": "G", "TT": ".",
			},
			input:   "AATTGGCC",
			wanted:  "ATGC",
			encoded: "M1D1M1D1M1D1M1D1",
		},
		{
			name: "first letter",
			reduction: map[string]string{
				"AA": "A", "AC": "A", "AG": "A", "AT": "A",
				"CA": "C", "CC": "C", "CG": "C", "CT": "C",
				"GA": "G", "GC": "G", "GG": "G", "GT": "G",
				"TA": "T", "TC": "T", "TG": "T", "TT": "T",
			},
			input:   "AGTCAGA",
			wanted:  "AAGTCAG",
			encoded: "M7",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			reduction := MakeReductionFunctionKeepOffsets(testCase.reduction)
			output, encoded := reduction(testCase.input)
			if output != testCase.wanted {
				t.Errorf("Error in making reduction function.\nWanted: %v\nGot: %v", testCase.wanted, output)
			}
			if !reflect.DeepEqual(encoded, testCase.encoded) {
				t.Errorf("Error, offsets are not the same.\nWanted: %v\nGot:%v", testCase.encoded, encoded)
			}
		})
	}
}

func makeBitArray(bits []bool)(*rsdic.RSDic) {
	rsd := rsdic.New()
	for _, bit := range bits {
		rsd.PushBack(bit)
	}
	return rsd
}

func bToi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func areBitArraysEqual(a1, a2 *rsdic.RSDic) (bool, string, string) {

	a1S, a2S := "", ""
	equal := true
	if a1.Num() != a2.Num() {
		return false, a1S, a2S
	}
	for i:=uint64(0); i<a1.Num(); i++ {
		b1, b2 := a1.Bit(i), a2.Bit(i)
		if b1 != b2 {
			equal = false
		}
		a1S += fmt.Sprintf("%d", bToi(b1))
		a2S += fmt.Sprintf("%d", bToi(b2))
	}
	return equal, a1S, a2S
}

func TestMakeReductionFunctionBitVector(t *testing.T) {
	cases := []struct {
		name      string
		reduction map[string]string
		input     string
		wanted    string
		bits []bool
	}{
		{
			name: "empty",
			reduction: map[string]string{
				"AA": ".", "AC": ".", "AG": ".", "AT": ".",
				"CA": ".", "CC": ".", "CG": ".", "CT": ".",
				"GA": ".", "GC": ".", "GG": ".", "GT": ".",
				"TA": ".", "TC": ".", "TG": ".", "TT": ".",
			},
			input:   "GATGCCAGTGCAGCATGTCAGGACGAC",
			wanted:  "G",
			bits: []bool{ //M1D26
				true, false, false, //M1D2
				false, false, false, false, //D4
				false, false, false, false, //D4
				false, false, false, false, //D4
				false, false, false, false, //D4
				false, false, false, false, //D4
				false, false, false, false, //D4
			},
		},
		{
			name: "homopolymer",
			reduction: map[string]string{
				"AA": ".", "AC": "C", "AG": "G", "AT": "T",
				"CA": "A", "CC": ".", "CG": "G", "CT": "T",
				"GA": "A", "GC": "C", "GG": ".", "GT": "T",
				"TA": "A", "TC": "C", "TG": "G", "TT": ".",
			},
			input:  "AAAAATTTTTTTGGGGGGCCCCCAAAAAATTTTTTGGGGGGGCCCCCCC",
			wanted: "ATGCATGC",
			bits: []bool{
				true, false, false, false, false, //A M1D4
				true, false, false, false, false, false, false, //T M1D6
				true, false, false, false, false, false, //G M1D5
				true, false, false, false, false, //C M1D4
				true, false, false, false, false, false, //A M1D5
				true, false, false, false, false, false, //T M1D5
				true, false, false, false, false, false, false, //G M1D6
				true, false, false, false, false, false, false, //C M1D6
			},
		},
		{
			name: "first letter",
			reduction: map[string]string{
				"AA": "A", "AC": "A", "AG": "A", "AT": "A",
				"CA": "C", "CC": "C", "CG": "C", "CT": "C",
				"GA": "G", "GC": "G", "GG": "G", "GT": "G",
				"TA": "T", "TC": "T", "TG": "T", "TT": "T",
			},
			input:  "AGTCAGA",
			wanted: "AAGTCAG",
			bits: []bool{true, true, true, true, true, true, true}, //M7
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			reduction := MakeReductionFunctionBitVector(testCase.reduction)
			output, offsets := reduction(testCase.input)
			if output != testCase.wanted {
				t.Errorf("Error in making reduction function.\nWanted %v\nGot: %v", testCase.wanted, output)
			}
			if int(offsets.OneNum()) != len(testCase.wanted) {
				fmt.Println(len(testCase.wanted))
				t.Errorf("Error in match number.\nWanted %d got %d", len(testCase.wanted), offsets.OneNum())
			}
			diff := len(testCase.input) - len(testCase.wanted)
			if int(offsets.ZeroNum()) != diff {
				t.Errorf("Error in deletion number.\nWanted %d got %d", diff, offsets.ZeroNum())
			}
			vector := makeBitArray(testCase.bits)
			equal, s1, s2 := areBitArraysEqual(vector, offsets)
			if !equal {
				t.Errorf("Offsets are not as we wish.\nW: [%v]\nG: [%v]", s1, s2)
			}
		})
	}
}
