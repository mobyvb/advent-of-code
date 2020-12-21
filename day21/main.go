package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Food struct {
	ingredients []string
	allergens   []string
}

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

	allFood := []*Food{}
	for _, line := range fileTextLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		newFood := &Food{}
		parts := strings.Split(line, " (contains ")
		newFood.ingredients = strings.Split(parts[0], " ")
		allergenString := strings.ReplaceAll(parts[1], ")", "")
		newFood.allergens = strings.Split(allergenString, ", ")
		allFood = append(allFood, newFood)
	}

	part1(allFood)
}

func part1(allFood []*Food) {
	// allergenMap maps each allergen to a list of possible ingredients containing that allergen
	allergenMap := make(map[string][]string)

	for _, f := range allFood {
		for _, a := range f.allergens {
			allergenIngredients := allergenMap[a]
			if allergenIngredients == nil {
				allergenMap[a] = f.ingredients
				continue
			}
			newAllergenIngredients := []string{}
			for _, i1 := range f.ingredients {
				for _, i2 := range allergenIngredients {
					// only keep ingredients that were in previous foods containing this allergen
					if i1 == i2 {
						newAllergenIngredients = append(newAllergenIngredients, i1)
					}
				}
			}
			allergenMap[a] = newAllergenIngredients
		}
	}

	allIngredients := make(map[string]int)
	for _, f := range allFood {
		for _, i := range f.ingredients {
			allIngredients[i]++
		}
	}

	ingredientsWithAllergens := make(map[string]bool)
	for _, iList := range allergenMap {
		for _, i := range iList {
			ingredientsWithAllergens[i] = true
		}
	}

	answer := 0
	ingredientsWithNoAllergens := []string{}
	for i, _ := range allIngredients {
		if !ingredientsWithAllergens[i] {
			ingredientsWithNoAllergens = append(ingredientsWithNoAllergens, i)
			answer += allIngredients[i]
		}
	}
	fmt.Println(answer)
}

func part2(allFood []*Food) {
}
