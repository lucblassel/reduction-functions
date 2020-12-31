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

// PhiRecord
type PhiRecord struct {
	Phi, C, F, Mu float64
}

func (record PhiRecord) String() string {
	return fmt.Sprintf(
		"{phi: %.4e,C: %.4e,F: %.4e,mu: %.4e}",
		record.Phi, record.C, record.F, record.Mu,
	)
}

func ComputeFarRatios(farSet []DistanceRecord) []float64 {
	ratios := make([]float64, len(farSet))
	for i, record := range farSet {
		ratios[i] = record.ReducedDistance / record.RawDistance
	}
	return ratios
}

func ComputeFarTerms(farRatios []float64, mu float64) []float64 {
	terms := make([]float64, len(farRatios))
	for i := range farRatios {
		terms[i] = math.Pow(farRatios[i]-mu, 2)
	}
	return terms
}

func ComputeCloseTerms(closeSet []DistanceRecord) []float64 {
	terms := make([]float64, len(closeSet))
	for i := range closeSet {
		terms[i] = closeSet[i].ReducedDistance
	}
	return terms
}

func SumArray(slice []float64) float64 {
	sum := 0.
	for i := range slice {
		sum += slice[i]
	}
	return sum
}

// ObjectivePhi is the objective function to evaluate the reduction function
func ObjectivePhi(closeSet, farSet []DistanceRecord) PhiRecord {

	C := SumArray(ComputeCloseTerms(closeSet)) / float64(len(closeSet))
	farRatios := ComputeFarRatios(farSet)
	mu := SumArray(farRatios) / float64(len(farSet))
	F := SumArray(ComputeFarTerms(farRatios, mu)) / float64(len(farRatios))

	return PhiRecord{
		Phi: C + F/mu,
		C:   C,
		F:   F,
		Mu:  mu,
	}
}
