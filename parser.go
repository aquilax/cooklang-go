package cooklang

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	COMMENTS_LINE_PREFIX     = "--"
	METADATA_LINE_PREFIX     = ">>"
	METADATA_VALUE_SEPARATOR = ":"
	PREFIX_INGREDIENT        = '@'
	PREFIX_Cookware          = '#'
	PREFIX_TIMER             = '~'
)

type Cookware struct {
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
	Cookware    []Cookware
	Comments    []string
}

type Metadata = map[string]string

type Recipe struct {
	Steps    []Step
	Metadata Metadata
}

func ParseFile(fileName string) (*Recipe, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseStream(bufio.NewReader(f))
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
		if strings.TrimSpace(line) != "" {
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
		recipe.Steps = append(recipe.Steps, Step{Comments: []string{commentLine}})
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
		Cookware:    make([]Cookware, 0),
	}
	skipIndex := -1
	var directions strings.Builder
	var err error
	var skipNext int
	var ingredient *Ingredient
	var Cookware *Cookware
	var timer *Timer
	for index, ch := range line {
		if skipIndex > index {
			continue
		}
		if ch == PREFIX_INGREDIENT {
			// ingredient ahead
			ingredient, skipNext, err = getIngredient(line[index:])
			if err != nil {
				return nil, err
			}
			skipIndex = index + skipNext
			step.Ingredients = append(step.Ingredients, *ingredient)
			directions.WriteString((*ingredient).Name)
		} else if ch == PREFIX_Cookware {
			// Cookware ahead
			Cookware, skipNext, err = getCookware(line[index:])
			if err != nil {
				return nil, err
			}
			skipIndex = index + skipNext
			step.Cookware = append(step.Cookware, *Cookware)
			directions.WriteString((*Cookware).Name)
		} else if ch == PREFIX_TIMER {
			//timer ahead
			timer, skipNext, err = getTimer(line[index:])
			if err != nil {
				return nil, err
			}
			skipIndex = index + skipNext
			step.Timers = append(step.Timers, *timer)
			directions.WriteString(fmt.Sprintf("%v %s", (*timer).Duration, (*timer).Unit))
		} else {
			// raw string
			directions.WriteRune(ch)
		}
	}
	step.Directions = directions.String()
	return &step, nil
}

func getCookware(line string) (*Cookware, int, error) {
	endIndex := findNodeEndIndex(line)
	Cookware, err := getCookwareFromRawString(line[1:endIndex])
	return Cookware, endIndex, err
}

func getIngredient(line string) (*Ingredient, int, error) {
	endIndex := findNodeEndIndex(line)
	ingredient, err := getIngredientFromRawString(line[1:endIndex])
	return ingredient, endIndex, err
}

func getTimer(line string) (*Timer, int, error) {
	endIndex := findNodeEndIndex(line)
	timer, err := getTimerFromRawString(line[2 : endIndex-1])
	return timer, endIndex, err
}

func getFloat(s string) (float64, error) {
	index := strings.Index(s, "/")
	if index == -1 {
		return strconv.ParseFloat(s, 64)
	}
	var err error
	var numerator int
	var denominator int
	numerator, err = strconv.Atoi(s[:index])
	if err != nil {
		return 0, err
	}

	denominator, err = strconv.Atoi(s[index+1:])
	if err != nil {
		return 0, err
	}
	return float64(numerator) / float64(denominator), nil
}

func findNodeEndIndex(line string) int {
	endIndex := -1

	for index, ch := range line {
		if index == 0 {
			continue
		}
		if (ch == PREFIX_Cookware || ch == PREFIX_INGREDIENT || ch == PREFIX_TIMER) && endIndex == -1 {
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
		return &Ingredient{Name: s, Amount: IngredientAmount{Quantity: 1}}, nil
	}
	amount, err := getAmount(s[index+1 : len(s)-1])
	if err != nil {
		return nil, err
	}
	return &Ingredient{Name: s[:index], Amount: *amount}, nil
}

func getAmount(s string) (*IngredientAmount, error) {
	if s == "" {
		return &IngredientAmount{Quantity: 1}, nil
	}
	index := strings.Index(s, "%")
	if index == -1 {
		f, err := getFloat(s)
		if err != nil {
			return nil, err
		}
		return &IngredientAmount{Quantity: f}, nil
	}
	f, err := getFloat(s[:index])
	if err != nil {
		return nil, err
	}
	return &IngredientAmount{Quantity: f, Unit: s[index+1:]}, nil
}

func getCookwareFromRawString(s string) (*Cookware, error) {
	return &Cookware{strings.TrimRight(s, "{}")}, nil
}

func getTimerFromRawString(s string) (*Timer, error) {
	index := strings.Index(s, "%")
	if index == -1 {
		return nil, fmt.Errorf("invalid timer syntax: %s", s)
	}
	f, err := getFloat(s[:index])
	if err != nil {
		return nil, err
	}
	return &Timer{Duration: f, Unit: s[index+1:]}, nil
}
