package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func loadFileList(path string) (files map[string]string, err error) {
	files = make(map[string]string)

	file, err := os.Open(path)

	if err != nil {
		return
	}

	defer func(file *os.File) {
		e := file.Close()
		if e != nil {
			err = e
		}
	}(file)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		value := scanner.Text()

		if value != "" {
			words := strings.Split(value, "\t")

			if len(words) != 2 {
				return files, errors.New(fmt.Sprintf("wrong file list line format, expected tab separated pair 'SAMPLE_NAME PATH', got '%s'", value))
			}

			fileName := words[0]
			filePath := words[1]

			files[fileName] = filePath
		}
	}

	if err := scanner.Err(); err != nil {
		return files, err
	}

	return
}
