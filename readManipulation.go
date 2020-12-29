package main

import (
	"errors"
	"fmt"
	"strings"
)

var basePairs = map[byte]rune{'A': 'T', 'G': 'C', 'C': 'G', 'T': 'A'}

// ReverseComplement gives the reverse complement of a given sequence
func ReverseComplement(seq string) (string, error) {
	if len(seq) == 0 {
		return "", errors.New("Cannot reverse complement empty string")
	}
	var sb strings.Builder
	for i := len(seq) - 1; i >= 0; i-- {
		rc, ok := basePairs[seq[i]]
		if !ok {
			return "", fmt.Errorf("Unkown nucleotide: %v", rc)
		}
		sb.WriteRune(rc)
	}
	return sb.String(), nil
}

// Canonize returns the canonical kmer
func Canonize(kmer string) (string, error) {
	rc, err := ReverseComplement(kmer)

	if err != nil {
		return "", err
	}

	if rc < kmer {
		return rc, nil
	}
	return kmer, nil
}

// Kmerize returns the set of canonical k-mers in a given sequence
func Kmerize(seq string, k int) (StringSet, error) {

	if len(seq) < k {
		return nil, errors.New("k is larger than the lenght of given read")
	}

	if k <= 1 {
		return nil, errors.New("k must be an integer > 1")
	}

	kmers := make([]string, 0, len(seq)-k+1)

	for i := 0; i < len(seq)-k+1; i++ {

		canonical, err := Canonize(seq[i : i+k])

		if err != nil {
			return nil, err
		}

		kmers = append(kmers, canonical)
	}

	return MakeSet(kmers), nil
}
