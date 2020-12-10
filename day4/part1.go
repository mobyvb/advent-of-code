package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	inputPath := os.Args[1]
	inputFile, err := os.Open(inputPath)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	fileScanner := bufio.NewScanner(inputFile)
	fileScanner.Split(bufio.ScanLines)
	var fileTextLines []string

	for fileScanner.Scan() {
		fileTextLines = append(fileTextLines, fileScanner.Text())
	}

	passportStrings := []string{}
	nextPassportString := ""
	for _, line := range fileTextLines {
		line = strings.TrimSpace(line)
		if line == "" {
			passportStrings = append(passportStrings, nextPassportString)
			nextPassportString = ""
			continue
		}
		nextPassportString += " "
		nextPassportString += line
	}

	validPassports := 0
	for _, passportString := range passportStrings {
		passportString := strings.Replace(passportString, "\n", " ", 0)
		passportData := make(map[string]string)
		fields := strings.Split(passportString, " ")
		for _, field := range fields {
			if field == "" {
				continue
			}
			fieldParts := strings.Split(field, ":")
			key := fieldParts[0]
			value := fieldParts[1]
			passportData[key] = value
		}
		requiredKeys := []string{
			"byr",
			"iyr",
			"eyr",
			"hgt",
			"hcl",
			"ecl",
			"pid",
			// "cid",
		}
		valid := true
		for _, key := range requiredKeys {
			if _, ok := passportData[key]; !ok {
				valid = false
			}
		}
		if valid {
			validPassports++

		}
	}
	fmt.Println(validPassports)
}
