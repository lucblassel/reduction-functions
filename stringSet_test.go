package reductions

import (
	"fmt"
	"testing"
)

func TestMakeSet(t *testing.T) {
	var tests = []struct {
		name     string
		elements []string
		wanted   StringSet
	}{
		{
			name: "MakeSetUniques", elements: []string{"ATG", "GAC", "GCC"},
			wanted: StringSet{"ATG": true, "GAC": true, "GCC": true},
		},
		{
			name: "MakeSetDUplicated", elements: []string{"ATG", "GAC", "GCC", "ATG", "GAC", "GCC", "ATG", "GAC", "GCC"},
			wanted: StringSet{"ATG": true, "GAC": true, "GCC": true},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ans := MakeSet(testCase.elements)
			if !ans.IsEqual(testCase.wanted) {
				t.Errorf("%s not equal to %s", ans.String(), testCase.wanted.String())
			}
		})
	}
}

func TestStringSet_IsEqual(t *testing.T) {
	set := MakeSet([]string{"ATG", "GAC", "GCC"})
	var tests = []struct {
		name       string
		set1, set2 StringSet
		wanted     bool
	}{
		{
			name: "EqualFull", set1: set,
			set2: set, wanted: true,
		},
		{
			name: "EqualEmpty", set1: StringSet{}, set2: StringSet{}, wanted: true,
		},
		{
			name: "DifferentElems", set1: set,
			set2: MakeSet([]string{"ATG", "GAC", "GCA"}), wanted: false,
		},
		{
			name: "DifferentLengths", set1: set,
			set2: MakeSet([]string{"ATG", "GAC", "GCA", "GCC"}), wanted: false,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ans := testCase.set1.IsEqual(testCase.set2)
			if ans != testCase.wanted {
				t.Errorf(
					"got equal=%v, wanted %v for sets:%s and %s",
					ans, testCase.wanted, testCase.set1.String(), testCase.set2.String(),
				)
			}
		})
	}
}

func TestStringSet_Intersection(t *testing.T) {
	set := MakeSet([]string{"ATG", "GAC", "GCC"})
	var tests = []struct {
		name               string
		set1, set2, wanted StringSet
	}{
		{
			name:   "EmptySet",
			set1:   set,
			set2:   StringSet{},
			wanted: StringSet{},
		},
		{
			name:   "Subsets",
			set1:   set,
			set2:   MakeSet([]string{"ATG", "GAC", "GCA", "GGG"}),
			wanted: MakeSet([]string{"ATG", "GAC"}),
		},
		{
			name:   "DisjointSets",
			set1:   set,
			set2:   MakeSet([]string{"GCA", "GGG", "TTT"}),
			wanted: StringSet{},
		},
		{
			name:   "SameSet",
			set1:   set,
			set2:   set,
			wanted: set,
		},
		{
			name:   "EmptySets",
			set1:   StringSet{},
			set2:   StringSet{},
			wanted: StringSet{},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ans := testCase.set1.Intersection(testCase.set2)
			if !ans.IsEqual(testCase.wanted) {
				t.Errorf("%s not equal to %s", ans.String(), testCase.wanted.String())
			}
		})
	}
}

func TestStringSet_Union(t *testing.T) {
	set := MakeSet([]string{"ATG", "GAC", "GCC"})
	var tests = []struct {
		name               string
		set1, set2, wanted StringSet
	}{
		{
			name:   "EmptySet",
			set1:   set,
			set2:   StringSet{},
			wanted: set,
		},
		{
			name:   "Subsets",
			set1:   set,
			set2:   MakeSet([]string{"ATG", "GAC", "GCA", "GGG"}),
			wanted: MakeSet([]string{"ATG", "GAC", "GCA", "GGG", "GCC"}),
		},
		{
			name:   "DisjointSets",
			set1:   set,
			set2:   MakeSet([]string{"GCA", "GGG", "TTT"}),
			wanted: MakeSet([]string{"ATG", "GAC", "GCC", "GCA", "GGG", "TTT"}),
		},
		{
			name:   "SameSet",
			set1:   set,
			set2:   set,
			wanted: set,
		},
		{
			name:   "EmptySets",
			set1:   StringSet{},
			set2:   StringSet{},
			wanted: StringSet{},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ans := testCase.set1.Union(testCase.set2)
			if !ans.IsEqual(testCase.wanted) {
				t.Errorf("%s not equal to %s", ans.String(), testCase.wanted.String())
			}
		})
	}
}

func TestStringSet_String(t *testing.T) {
	set := MakeSet([]string{"AAA", "GGG", "TTT"})
	wanted := fmt.Sprintf("Set{AAA,GGG,TTT,(3)}")
	if ans := set.String(); ans != wanted {
		t.Errorf("Wanted: %v\n got: %v", wanted, ans)
	}
}

func BenchmarkStringSet_Intersection(b *testing.B) {
	seqs, _, _ := ParseFasta("test.fasta")
	kmers1, _ := Kmerize(seqs["Seq01"], 5)
	kmers2, _ := Kmerize(seqs["Seq02"], 5)
	for i := 0; i < b.N; i++ {
		kmers1.Intersection(kmers2)
	}
}
