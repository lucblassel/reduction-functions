package main

import (
	"testing"
)

func TestIsEqual(t *testing.T) {
	var tests = []struct {
		name       string
		set1, set2 StringSet
		wanted     bool
	}{
		{
			name: "EqualFull", set1: StringSet{"ATG": true, "GAC": true, "GCC": true},
			set2: StringSet{"ATG": true, "GAC": true, "GCC": true}, wanted: true,
		},
		{
			name: "EqualEmpty", set1: StringSet{}, set2: StringSet{}, wanted: true,
		},
		{
			name: "DifferentElems", set1: StringSet{"ATG": true, "GAC": true, "GCC": true},
			set2: StringSet{"ATG": true, "GAC": true, "GCA": true}, wanted: false,
		},
		{
			name: "DifferentLengths", set1: StringSet{"ATG": true, "GAC": true, "GCC": true},
			set2: StringSet{"ATG": true, "GAC": true, "GCA": true, "GCC": true}, wanted: false,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ans := testCase.set1.IsEqual(testCase.set2)
			if ans != testCase.wanted {
				t.Errorf(
					"got equal=%v, wanted %v for sets:%s and %s",
					ans, testCase.wanted, testCase.set1.ToString(), testCase.set2.ToString(),
				)
			}
		})
	}
}

func TestMakeSets(t *testing.T) {
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
				t.Errorf("%s not equal to %s", ans.ToString(), testCase.wanted.ToString())
			}
		})
	}
}

func TestSetIntersection(t *testing.T) {
	var tests = []struct {
		name               string
		set1, set2, wanted StringSet
	}{
		{
			name:   "EmptySet",
			set1:   MakeSet([]string{"ATG", "GAC", "GCC"}),
			set2:   StringSet{},
			wanted: StringSet{},
		},
		{
			name:   "Subsets",
			set1:   MakeSet([]string{"ATG", "GAC", "GCC"}),
			set2:   MakeSet([]string{"ATG", "GAC", "GCA", "GGG"}),
			wanted: MakeSet([]string{"ATG", "GAC"}),
		},
		{
			name:   "DisjointSets",
			set1:   MakeSet([]string{"ATG", "GAC", "GCC"}),
			set2:   MakeSet([]string{"GCA", "GGG", "TTT"}),
			wanted: StringSet{},
		},
		{
			name:   "SameSet",
			set1:   MakeSet([]string{"ATG", "GAC", "GCC"}),
			set2:   MakeSet([]string{"ATG", "GAC", "GCC"}),
			wanted: MakeSet([]string{"ATG", "GAC", "GCC"}),
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
				t.Errorf("%s not equal to %s", ans.ToString(), testCase.wanted.ToString())
			}
		})
	}
}

func TestSetUnion(t *testing.T) {
	var tests = []struct {
		name               string
		set1, set2, wanted StringSet
	}{
		{
			name:   "EmptySet",
			set1:   MakeSet([]string{"ATG", "GAC", "GCC"}),
			set2:   StringSet{},
			wanted: MakeSet([]string{"ATG", "GAC", "GCC"}),
		},
		{
			name:   "Subsets",
			set1:   MakeSet([]string{"ATG", "GAC", "GCC"}),
			set2:   MakeSet([]string{"ATG", "GAC", "GCA", "GGG"}),
			wanted: MakeSet([]string{"ATG", "GAC", "GCA", "GGG", "GCC"}),
		},
		{
			name:   "DisjointSets",
			set1:   MakeSet([]string{"ATG", "GAC", "GCC"}),
			set2:   MakeSet([]string{"GCA", "GGG", "TTT"}),
			wanted: MakeSet([]string{"ATG", "GAC", "GCC", "GCA", "GGG", "TTT"}),
		},
		{
			name:   "SameSet",
			set1:   MakeSet([]string{"ATG", "GAC", "GCC"}),
			set2:   MakeSet([]string{"ATG", "GAC", "GCC"}),
			wanted: MakeSet([]string{"ATG", "GAC", "GCC"}),
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
				t.Errorf("%s not equal to %s", ans.ToString(), testCase.wanted.ToString())
			}
		})
	}
}
