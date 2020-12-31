package reductions

import (
	"bufio"
	"fmt"
	"os"
)

func getLineScanner(path string) (*bufio.Scanner, *os.File, error) {
	var scanner *bufio.Scanner
	file, err := os.Open(path)
	if err != nil {
		return scanner, file, err
	}

	scanner = bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	return scanner, file, nil
}

//ParseFasta reads a fasta file and returns a map of sequences and ids as strings
func ParseFasta(path string) (map[string]string, error) {
	sequences := make(map[string]string)
	var key, sequence string

	scanner, file, err := getLineScanner(path)
	if err != nil {
		return sequences, err
	}

	for scanner.Scan() {
		line := scanner.Text()
		if rune(line[0]) == '>' {
			if key != "" {
				sequences[key] = sequence
				sequence = ""
			}
			key = line[1:]
		} else {
			sequence += line
		}
	}
	sequences[key] = sequence

	err = file.Close()
	if err != nil {
		return sequences, err
	}

	return sequences, nil
}

// ParseWFA reads the dataset of sequences produced by the generate_dataset tool
// of the WFA paper
func ParseWFA(path string) (map[string]string, error) {
	sequences := make(map[string]string)
	var key string

	scanner, file, err := getLineScanner(path)
	if err != nil {
		return sequences, err
	}

	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		if rune(line[0]) == '>' {
			i++
			key = fmt.Sprintf("seq_%d", i)
			sequences[key] = line[1:]
		}
		if rune(line[0]) == '<' {
			key = fmt.Sprintf("seq_%d_err", i)
			sequences[key] = line[1:]
		}
	}

	err = file.Close()
	if err != nil {
		return sequences, err
	}

	return sequences, nil
}

func printSequences(seqs map[string]string) {
	for key, sequence := range seqs {
		fmt.Printf("%s:\t%s\n", key, sequence)
	}
}
