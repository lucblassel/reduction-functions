package reductions

import (
	"errors"
	"fmt"
	"gonum.org/v1/gonum/stat/combin"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"
	"github.com/hillbig/rsdic"
)

// CountSurjections counts the number of surjections from a set of
// nStart elements to a set of nEnd elements
func CountSurjections(nStart, nEnd int) int {
	res := 0
	for k := 0; k <= nEnd; k++ {
		sign := int(math.Pow(-1, float64(nEnd-k)))
		combs := combin.Binomial(nEnd, k)
		power := int(math.Pow(float64(k), float64(nStart)))
		res += sign * combs * power
	}
	return res
}

// RandomMapping returns a random mapping from a input set to an output set
func RandomMapping(nStart, nEnd int) []int {
	rand.Seed(time.Now().UnixNano())
	choices := make([]int, nStart)
	for i := 0; i < nStart; i++ {
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
	for i := 0; i < nEnd; i++ {
		pos := rand.Intn(nStart)
		for chosen[pos] {
			pos = rand.Intn(nStart)
		}
		choices[pos] = i
		chosen[pos] = true
	}
	for i := 0; i < nStart; i++ {
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
			output = append(output, elem+string(char))
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
func GetRandomReduction(inputAlphabet, outputAlphabet string, inputSize, outputSize int) (map[string]string, error) {
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

// MakeReductionFunction create a reduction function from a mapping
func MakeReductionFunction(surjection map[string]string) func(string) string {
	return func(read string) string {
		var order int
		for k := range surjection {
			order = len(k)
			break
		}

		var builder strings.Builder
		builder.WriteString(read[0 : order-1])
		for i := 0; i <= len(read)-order; i++ {
			s := surjection[read[i:i+order]]
			if s == "." {
				continue
			}
			builder.WriteString(s)
		}
		return builder.String()
	}
}

// MakeReductionFunctionKeepOffsets create a reduction function from a mapping
func MakeReductionFunctionKeepOffsets(surjection map[string]string) func(string) (string, string) {
	return func(read string) (string, string) {
		var order int
		for k := range surjection {
			order = len(k)
			break
		}

		var builder strings.Builder
		encoded := ""
		op, count := "M", order-1

		builder.WriteString(read[0 : order-1])

		for i := 0; i <= len(read)-order; i++ {
			s := surjection[read[i:i+order]]
			if s == "." {
				if op == "M" {
					encoded += fmt.Sprintf("M%d", count)
					count = 0
					op = "D"
				}
				count++
				continue
			}
			if op == "D" {
				encoded += fmt.Sprintf("D%d", count)
				count = 0
				op = "M"
			}
			count++
			builder.WriteString(s)
		}
		encoded += fmt.Sprintf("%s%d", op, count)

		return builder.String(), encoded
	}
}


// MakeReductionFunctionBitVector create a reduction function from a mapping
func MakeReductionFunctionBitVector(surjection map[string]string) func(string) (string, *rsdic.RSDic) {
	return func(read string) (string, *rsdic.RSDic) {
		var order int
		for k := range surjection {
			order = len(k)
			break
		}
		offsets := rsdic.New()

		for i:=0; i<order-1; i++ {
			offsets.PushBack(true)
		}

		var builder strings.Builder
		builder.WriteString(read[0 : order-1])
		for i := 0; i <= len(read)-order; i++ {
			s := surjection[read[i:i+order]]
			if s == "." {
				offsets.PushBack(false)
				continue
			} else {
				offsets.PushBack(true)
			}
			builder.WriteString(s)
		}
		return builder.String(), offsets
	}
}