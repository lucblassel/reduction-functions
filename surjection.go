package reductions

import (
	"errors"
	"gonum.org/v1/gonum/stat/combin"
	"math"
	"math/rand"
	"sort"
	"time"
)

// CountSurjections counts the number of surjections from a set of
// nStart elements to a set of nEnd elements
func CountSurjections(nStart, nEnd int) int {
	res := 0
	for k:=0; k<=nEnd; k++ {
		sign := int(math.Pow(-1, float64(nEnd - k)))
		combs := combin.Binomial(nEnd, k)
		power := int(math.Pow(float64(k), float64(nStart)))
		res += sign * combs * power
	}
	return res
}

// RandomMapping returns a random mapping from a input set to an output set
func RandomMapping(nStart, nEnd int) []int{
		rand.Seed(time.Now().UnixNano())
		choices := make([]int, nStart)
		for i:=0; i < nStart; i++ {
			choices[i] = rand.Intn(nEnd)
		}
		return choices
}

// Surjection returns a random surjection from an input set into an output set
func Surjection(nStart, nEnd int) ([]int, error) {
	if nEnd > nStart {
		return nil, errors.New("end set must be smaller or of the same size as the starting set")
	}
	rand.Seed(time.Now().UnixNano())
	choices := make([]int, nStart)
	chosen := make([]bool, nStart)
	for i:=0; i<nEnd; i++ {
		pos:=rand.Intn(nStart)
		for chosen[pos] {
			pos=rand.Intn(nStart)
		}
		choices[pos] = i
		chosen[pos] = true
	}
	for i:=0; i<nStart; i++ {
		if chosen[i] {
			continue
		}
		choices[i] = rand.Intn(nEnd)
		chosen[i] = true
	}
	return choices, nil
}

// GetTuples returns all possible samplings of size n with replacements from an input alphabet
func GetTuples(n, depth int, universe string, samples []string) []string {
	if depth == n {
		return samples
	}
	output := make([]string, 0)
	for _, elem := range samples {
		for _, char := range universe {
			output = append(output, elem + string(char))
		}
	}
	if len(samples) == 0 {
		for _, char := range universe {
			output = append(output, string(char))
		}
	}
	return GetTuples(n, depth+1, universe, output)
}

// GetRandomReduction generates a random surjection between a set of input and output sequences
func GetRandomReduction(inputAlphabet, outputAlphabet  string, inputSize, outputSize int) (map[string]string, error){
	inputPerms := GetTuples(inputSize, 0, inputAlphabet, []string{})
	outputPerms := GetTuples(outputSize, 0, outputAlphabet, []string{})
	sort.Strings(inputPerms)
	sort.Strings(outputPerms)
	mapping, err := Surjection(len(inputPerms), len(outputPerms))
	if err != nil {
		return nil, err
	}
	function := make(map[string]string)
	for inIdx, outIdx := range mapping {
		function[inputPerms[inIdx]] = outputPerms[outIdx]
	}
	return function, nil
}