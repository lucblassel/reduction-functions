package reductions

import (
	"math/rand"
	"testing"
)

func isSurjection(nEnd int, mapping []int) bool{
	outputs := make(map[int]bool)
	for _, elem := range mapping {
		outputs[elem] = true
	}
	for i:=0; i<nEnd; i++ {
		if !outputs[i] {
			return false
		}
	}
	return true
}

func TestRandomMapping(t *testing.T) {
	for i:=0; i<100; i++ {
		start, end := rand.Intn(99) + 1, rand.Intn(99) + 1
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
	for i:=0; i<100; i++ {
		start, end := rand.Intn(99) + 1, rand.Intn(99) + 1
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
