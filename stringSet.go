package main

import (
	"fmt"
	"strings"
)

// StringSet is a set of unique strings
type StringSet map[string]bool

// MakeSet makes a set out of a slice of strings
func MakeSet(seqs []string) StringSet {
	set := StringSet{}
	for _, seq := range seqs {
		set[seq] = true
	}
	return set
}

// Intersection returns the intersection of 2 sets of strings
func (set1 StringSet) Intersection(set2 StringSet) StringSet {
	intersection := StringSet{}
	for key := range set1 {
		if set2[key] {
			intersection[key] = true
		}
	}
	return intersection
}

// Union returns the union of 2 sets of strings
func (set1 StringSet) Union(set2 StringSet) StringSet {
	union := StringSet{}
	for key := range set1 {
		union[key] = true
	}
	for key := range set2 {
		union[key] = true
	}
	return union
}

// IsEqual returns true if the 2 sets are equal and false if not
func (set1 StringSet) IsEqual(set2 StringSet) bool {
	if len(set1) != len(set2) {
		return false
	}

	for k := range set1 {
		if !set2[k] {
			return false
		}
	}

	return true
}

// ToString returns a string readable string representation of a set
func (set1 StringSet) ToString() string {
	var sb strings.Builder
	sb.WriteString("Set{")
	for k := range set1 {
		sb.WriteString(k + ",")
	}
	sb.WriteString(fmt.Sprintf("(%d)}", len(set1)))
	return sb.String()
}
