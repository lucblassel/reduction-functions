package reductions

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// getLineScanner returns a line by line scanner and the corresponding file (for closing)
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
func ParseFasta(path string) (map[string]string, []string, error) {
	sequences := make(map[string]string)
	order := make([]string, 0, 0)
	var key, sequence string

	scanner, file, err := getLineScanner(path)
	if err != nil {
		return nil, nil, err
	}

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("parsing line", line[:10], "...")
		if rune(line[0]) == '>' {
			if key != "" {
				sequences[key] = sequence
				sequence = ""
			}
			key = strings.TrimSpace(line[1:])
			order = append(order, key)
		} else if strings.TrimSpace(line) == "" {
			continue
		} else {
			sequence += line
		}
	}
	sequences[key] = sequence

	err = file.Close()
	if err != nil {
		return nil, nil, err
	}

	return sequences, order, nil
}

// write a single sequence with ID to fasta file
func writeSequence(file *os.File, name , sequence string) error {
	_, err := file.WriteString(fmt.Sprintf(">%s\n", name))
	if err != nil {
		return err
	}
	_, err = file.WriteString(fmt.Sprintf("%s\n", sequence))
	if err != nil {
		return err
	}
	return nil
}


// WriteFasta saves a collection of sequences to a FASTA formatted file
func WriteFasta(sequences map[string]string, path string, order []string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	if len(order) != 0 && len(order) != len(sequences) {
		return errors.New("Key order and sequences must have the same length")
	}

	defer file.Close()

	if len(order) == 0 {
		for name, sequence := range sequences {
			err := writeSequence(file, name, sequence)
			if err != nil {
				return err
			}
		}
	} else {
		for _, name := range order {
			err := writeSequence(file, name, sequences[name])
			if err != nil {
				return err
			}
		}
	}

	err = file.Sync()
	if err != nil {
		return err
	}
	return nil
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

// CheckSurjectionFile unmarshalls a reduction function .json to a map
func CheckSurjectionFile(path string, output *map[string]string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, output)
	if err != nil {
		return err
	}
	return nil
}
