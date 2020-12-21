package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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

	allergenMap := part1(allFood)
	part2(allergenMap)
}

func part1(allFood []*Food) (allergenMap map[string][]string) {
	// allergenMap maps each allergen to a list of possible ingredients containing that allergen
	allergenMap = make(map[string][]string)

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
	return allergenMap
}

func part2(allergenMap map[string][]string) {
	// another sudoku-style problem
	// while we have not determined the ingredient associated with each allergen,
	// iterate over all allergens, and find an allergen with only one ingredient
	// this ingredient is now a danger food. Remove references to it and its allergens from the allergen map
	type allergenPair struct {
		ingredient string
		allergen   string
	}
	dangerFoods := []*allergenPair{}
	for len(allergenMap) > 0 {
		newAllergenMap := make(map[string][]string)
		for a, iList := range allergenMap {
			if len(iList) == 1 {
				newDangerFood := &allergenPair{
					ingredient: iList[0],
					allergen:   a,
				}
				dangerFoods = append(dangerFoods, newDangerFood)
				for a2, iList := range allergenMap {
					if a2 == a {
						continue
					}
					for _, i := range iList {
						if i == newDangerFood.ingredient {
							continue
						}
						newAllergenMap[a2] = append(newAllergenMap[a2], i)
					}
				}

				break
			}
		}
		allergenMap = newAllergenMap
	}
	// sort dangerFoods alphabetically by allergen
	sort.SliceStable(dangerFoods, func(i, j int) bool {
		return strings.Compare(dangerFoods[i].allergen, dangerFoods[j].allergen) < 0
	})
	dangerIngredients := []string{}
	for _, f := range dangerFoods {
		dangerIngredients = append(dangerIngredients, f.ingredient)
	}
	canonicalDangerousIngredients := strings.Join(dangerIngredients, ",")
	fmt.Println(canonicalDangerousIngredients)
}
