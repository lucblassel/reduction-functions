package reductions

import (
	"errors"
	"math/rand"
	"time"
)

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