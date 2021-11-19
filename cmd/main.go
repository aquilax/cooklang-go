package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/aquilax/cooklang-go"
)

const OFFSET_INDENT = 4

func main() {
	recipe, err := cooklang.ParseFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	printRecipe(*recipe, os.Stdout)
}

func collectIngredients(steps []cooklang.Step) []cooklang.Ingredient {
	var result []cooklang.Ingredient
	for i := range steps {
		result = append(result, steps[i].Ingredients...)
	}
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}

func coollectCookware(steps []cooklang.Step) []string {
	var result []string
	for i := range steps {
		for j := range steps[i].Cookware {
			result = append(result, steps[i].Cookware[j].Name)
		}
	}
	sort.Strings(result)
	return result
}

func formatFloat(num float64, precision int) string {
	fs := fmt.Sprintf("%%.%df", precision)
	s := fmt.Sprintf(fs, num)
	return strings.TrimRight(strings.TrimRight(s, "0"), ".")
}

func getIngredients(ing []cooklang.Ingredient) []string {
	var result []string
	for i := range ing {
		result = append(result, fmt.Sprintf("%s: %s %s", ing[i].Name, formatFloat(ing[i].Amount.Quantity, 2), ing[i].Amount.Unit))
	}
	sort.Strings(result)
	return result
}

func printRecipe(recipe cooklang.Recipe, out io.Writer) {
	offset := strings.Repeat(" ", OFFSET_INDENT)
	if len(recipe.Metadata) > 0 {
		fmt.Fprintln(out, "Metadata:")
		for k, v := range recipe.Metadata {
			fmt.Fprintf(out, "%s%s: %s\n", offset, k, v)
		}
		fmt.Fprintln(out, "")
	}
	allIngredients := collectIngredients(recipe.Steps)
	if len(allIngredients) > 0 {
		fmt.Fprintln(out, "Ingredients:")
		for i := range allIngredients {
			fmt.Fprintf(out, "%s%-30s%s %s\n", offset, allIngredients[i].Name, formatFloat(allIngredients[i].Amount.Quantity, 2), allIngredients[i].Amount.Unit)
		}
		fmt.Fprintln(out, "")
	}
	allCookware := coollectCookware(recipe.Steps)
	if len(allCookware) > 0 {
		fmt.Fprintln(out, "Cookware:")
		for i := range allCookware {
			fmt.Fprintf(out, "%s%s\n", offset, allCookware[i])
		}
		fmt.Fprintln(out, "")
	}
	if len(recipe.Steps) > 0 {
		fmt.Fprintln(out, "Steps:")
		for i := range recipe.Steps {
			fmt.Fprintf(out, "%s%2d. %s\n", offset, i+1, recipe.Steps[i].Directions)
			ingredients := "–"
			ing := getIngredients(recipe.Steps[i].Ingredients)
			if len(ing) > 0 {
				ingredients = strings.Join(ing, "; ")
			}

			fmt.Fprintf(out, "%s    [%s]\n", offset, ingredients)
		}
	}
}
