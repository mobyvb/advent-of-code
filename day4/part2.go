package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type fieldValidation func(string) bool

func yearValidation(min, max int) fieldValidation {
	return func(s string) bool {
		x, err := strconv.Atoi(s)
		if err != nil {
			return false
		}
		if x < min || x > max {
			return false
		}
		return true

	}
}

func main() {
	requiredKeys := map[string]fieldValidation{
		"byr": yearValidation(1920, 2002),
		"iyr": yearValidation(2010, 2020),
		"eyr": yearValidation(2020, 2030),
		"hgt": func(s string) bool {
			unitIndex := len(s) - 2
			valueString := s[:unitIndex]
			unit := s[unitIndex:]
			value, err := strconv.Atoi(valueString)
			if err != nil {
				return false
			}
			if unit != "in" && unit != "cm" {
				return false
			}
			if unit == "cm" && (value < 150 || value > 193) {
				return false
			}
			if unit == "in" && (value < 59 || value > 76) {
				return false
			}
			return true
		},
		"hcl": func(s string) bool {
			if len(s) != 7 {
				return false
			}
			match, _ := regexp.MatchString("#([0-9]|[a-f]){6}", s)
			return match
		},
		"ecl": func(s string) bool {
			match, _ := regexp.MatchString("amb|blu|brn|gry|grn|hzl|oth", s)
			return match
		},
		"pid": func(s string) bool {
			if len(s) != 9 {
				return false
			}
			match, _ := regexp.MatchString("[0-9]{9}", s)
			return match
		},
		// "cid",
	}

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

		valid := true
		for key, validation := range requiredKeys {
			if _, ok := passportData[key]; !ok {
				valid = false
			} else if !validation(passportData[key]) {
				valid = false
			}
		}
		if valid {
			validPassports++

		}
	}
	fmt.Println(validPassports)
}
