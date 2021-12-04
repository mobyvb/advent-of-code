package common

import (
	"bufio"
	"os"
	"strconv"
)

type File struct {
	path  string
	lines []string
}

func OpenFile(path string) (*File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	var fileTextLines []string

	for fileScanner.Scan() {
		fileTextLines = append(fileTextLines, fileScanner.Text())
	}
	return &File{
		path:  path,
		lines: fileTextLines,
	}, nil
}

// GetInts tries to convert each line in the file to an integer and return the list.
func (f *File) GetInts() ([]int, error) {
	intList := make([]int, len(f.lines))
	for i, l := range f.lines {
		value, err := strconv.Atoi(l)
		if err != nil {
			return intList, err
		}
		intList[i] = value
	}
	return intList, nil
}

// GetBinary tries to parse each line as binary into an integer.
func (f *File) GetBinary() ([]int, error) {
	intList := make([]int, len(f.lines))
	for i, l := range f.lines {
		value, err := strconv.ParseInt(l, 2, 32)
		if err != nil {
			return intList, err
		}
		intList[i] = int(value)
	}
	return intList, nil
}
