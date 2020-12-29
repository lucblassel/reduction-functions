package main

import (
	"bufio"
	"os"
)

//ParseFasta reads a fasta file and returns a map of sequences and ids as strings
func ParseFasta(path string) (map[string]string, error) {
	sequences := make(map[string]string)
	file, err := os.Open(path)
	var key, sequence string
	if err != nil {
		return sequences, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

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
		return map[string]string{}, err
	}

	return sequences, nil
}
