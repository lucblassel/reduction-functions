package reductions

import (
	"fmt"
	"math"
	"strings"
)

// HomopolymerCompression returns the homopolymer compressed version of a read
func HomopolymerCompression(read string) string {
	var lastChar rune
	var builder strings.Builder

	for _, char := range read {
		if char == lastChar {
			continue
		}
		builder.WriteRune(char)
		lastChar = char
	}

	return builder.String()
}

// Identity returns the input read unchanged
func Identity(read string) string {
	return read
}

// PhiRecord records all the computed terms of the Objective function
type PhiRecord struct {
	Phi, C, F, Mu float64
}

// String impelents the Stringer interface for PhiRecord structs
func (record PhiRecord) String() string {
	return fmt.Sprintf(
		"{phi: %.4e,C: %.4e,F: %.4e,mu: %.4e}",
		record.Phi, record.C, record.F, record.Mu,
	)
}

// ComputeDistanceRatio returns the ratio of the ReducedDistance / RawDistance for
// a slice of DistanceRecords
func ComputeDistanceRatio(farSet []DistanceRecord) []float64 {
	ratios := make([]float64, len(farSet))
	for i, record := range farSet {
		ratios[i] = record.ReducedDistance / record.RawDistance
	}
	return ratios
}

// ComputeSquareDifference return the square difference between all elements and a float mu
func ComputeSquareDifference(farRatios []float64, mu float64) []float64 {
	terms := make([]float64, len(farRatios))
	for i := range farRatios {
		terms[i] = math.Pow(farRatios[i]-mu, 2)
	}
	return terms
}

// GetReducedDistance returns the all the reduced distances of a slice of DistanceRecords
func GetReducedDistance(closeSet []DistanceRecord) []float64 {
	terms := make([]float64, len(closeSet))
	for i := range closeSet {
		terms[i] = closeSet[i].ReducedDistance
	}
	return terms
}

// SumArray returns the sum of a slice of float64
func SumArray(slice []float64) float64 {
	sum := 0.
	for i := range slice {
		sum += slice[i]
	}
	return sum
}

// ObjectivePhi is the objective function to evaluate the reduction function
func ObjectivePhi(closeSet, farSet []DistanceRecord) PhiRecord {

	C := SumArray(GetReducedDistance(closeSet)) / float64(len(closeSet))
	farRatios := ComputeDistanceRatio(farSet)
	mu := SumArray(farRatios) / float64(len(farSet))
	F := SumArray(ComputeSquareDifference(farRatios, mu)) / float64(len(farRatios))

	return PhiRecord{
		Phi: C + F/mu,
		C:   C,
		F:   F,
		Mu:  mu,
	}
}
