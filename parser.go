package cooklang

import (
	"bufio"
	"fmt"
	"strings"
)

const (
	COMMENTS_LINE_PREFIX     = "--"
	METADATA_LINE_PREFIX     = ">>"
	METADATA_VALUE_SEPARATOR = ":"
)

type Equipment struct {
	Name string
}

type IngredientAmount struct {
	Quantity float64
	Unit     string
}

type Ingredient struct {
	Name   string
	Amount *IngredientAmount
}

type Timer struct {
	Duration float64
	Unit     string
}

type Step struct {
	Directions  string
	Timers      []Timer
	Ingredients []Ingredient
	Equipment   []Equipment
	Comments    string
}

type Metadata = map[string]string

type Recipe struct {
	Steps    []Step
	Metadata Metadata
}

func ParseString(s string) (*Recipe, error) {
	if s == "" {
		return nil, fmt.Errorf("Recipe string must not be empty")
	}

	// TODO parse recipe
	scanner := bufio.NewScanner(strings.NewReader(s))
	recipe := Recipe{
		make([]Step, 0),
		make(map[string]string),
	}
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		if line != "" {
			err := parseLine(line, &recipe)
			if err != nil {
				return nil, err
			}
		}
		fmt.Println(scanner.Text())
	}
	return &recipe, nil
}

func parseLine(line string, recipe *Recipe) error {
	if strings.HasPrefix(line, COMMENTS_LINE_PREFIX) {
		commentLine, err := parseSingleLineComment(line)
		if err != nil {
			return err
		}
		recipe.Steps = append(recipe.Steps, Step{Comments: commentLine})
	} else if strings.HasPrefix(line, METADATA_LINE_PREFIX) {
		key, value, err := parseMetadata(line)
		if err != nil {
			return err
		}
		recipe.Metadata[key] = value
	} else {
		step, err := parseRecipe(line)
		if err != nil {
			return err
		}
		recipe.Steps = append(recipe.Steps, *step)
	}
	return nil
}

func parseSingleLineComment(line string) (string, error) {
	return strings.TrimSpace(line[2:]), nil
}

func parseMetadata(line string) (string, string, error) {
	metadataLine := strings.TrimSpace(line[2:])
	index := strings.Index(metadataLine, METADATA_VALUE_SEPARATOR)
	if index < 1 {
		return "", "", fmt.Errorf("invalid metadata: %s", metadataLine)
	}
	return strings.TrimSpace(metadataLine[:index]), strings.TrimSpace(metadataLine[index+1:]), nil
}

func parseRecipe(line string) (*Step, error) {
	return &Step{}, nil
}
