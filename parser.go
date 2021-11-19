package cooklang

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	COMMENTS_LINE_PREFIX     = "--"
	METADATA_LINE_PREFIX     = ">>"
	METADATA_VALUE_SEPARATOR = ":"
	PREFIX_INGREDIENT        = '@'
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
	Amount IngredientAmount
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
		return nil, fmt.Errorf("recipe string must not be empty")
	}
	return ParseStream(strings.NewReader(s))
}

func ParseStream(s io.Reader) (*Recipe, error) {
	scanner := bufio.NewScanner(s)
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
	step := Step{
		Timers:      make([]Timer, 0),
		Ingredients: make([]Ingredient, 0),
		Equipment:   make([]Equipment, 0),
	}
	skipIndex := -1
	var directions strings.Builder
	var err error
	var skipNext int
	var ingredient *Ingredient
	for index, ch := range line {
		if skipIndex > index {
			continue
		}
		if ch == '@' {
			ingredient, skipNext, err = getIngredient(line[index:])
			if err != nil {
				return nil, err
			}
			skipIndex = index + skipNext
			step.Ingredients = append(step.Ingredients, *ingredient)
			directions.WriteString((*ingredient).Name)
			// ingredient ahead
		} else if ch == '#' {
			// equipment ahead
		} else if ch == '~' {
			//timer ahead
		} else {
			// raw string
			directions.WriteRune(ch)
		}
	}
	step.Directions = directions.String()
	return &step, nil
}

func getIngredient(line string) (*Ingredient, int, error) {
	endIndex := findNodeEndIndex(PREFIX_INGREDIENT, line)
	ingredient, error := getIngredientFromRawString(line[1:endIndex])
	return ingredient, endIndex, error
}

func findNodeEndIndex(prefix rune, line string) int {
	endIndex := -1

	for index, ch := range line {
		if index == 0 {
			continue
		}
		if ch == prefix && endIndex == -1 {
			break
		}
		if ch == '}' {
			endIndex = index + 1
			break
		}
	}
	if endIndex == -1 {
		endIndex = strings.Index(line, " ")
		if endIndex == -1 {
			endIndex = len(line)
		}
	}
	return endIndex
}

func getIngredientFromRawString(s string) (*Ingredient, error) {
	index := strings.Index(s, "{")
	if index == -1 {
		return &Ingredient{Name: s}, nil
	}
	amount, err := getAmount(s[index+1 : len(s)-1])
	if err != nil {
		return nil, err
	}
	return &Ingredient{Name: s[:index], Amount: *amount}, nil
}

func getAmount(s string) (*IngredientAmount, error) {
	index := strings.Index(s, "%")
	if index == -1 {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		return &IngredientAmount{Quantity: f}, nil
	}
	f, err := strconv.ParseFloat(s[:index], 64)
	if err != nil {
		return nil, err
	}
	return &IngredientAmount{Quantity: f, Unit: s[index+1:]}, nil
}
