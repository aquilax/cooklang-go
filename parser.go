// Package cooklang provides a parser for .cook defined recipes as defined in
// https://cooklang.org/docs/spec/
package cooklang

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	commentsLinePrefix     = "--"
	metadataLinePrefix     = ">>"
	metadataValueSeparator = ":"
	prefixIngredient       = '@'
	prefixCookware         = '#'
	prefixTimer            = '~'
	prefixBlockComment     = '['
)

// Cookware represents a cookware item
type Cookware struct {
	Name string // cookware name
}

// IngredientAmount represents the amount required of an ingredient
type IngredientAmount struct {
	IsNumeric   bool    // true if the amount is numeric
	Quantity    float64 // quantity of the ingredient
	QuantityRaw string  // quantity of the ingredient as raw text
	Unit        string  // optional ingredient unit
}

// Ingredient represents a recipe ingredient
type Ingredient struct {
	Name   string           // name of the ingredient
	Amount IngredientAmount // optional ingredient amount (default: 1)
}

// Timer represents a time duration
type Timer struct {
	Duration float64 // duration of the timer
	Unit     string  // time unit of the duration
}

// Step represents a recipe step
type Step struct {
	Directions  string       // step directions as plain text
	Timers      []Timer      // list of timers in the step
	Ingredients []Ingredient // list of ingredients used in the step
	Cookware    []Cookware   // list of cookware used in the step
	Comments    []string     // list of comments
}

// Metadata contains key value map of metadata
type Metadata = map[string]string

// Recipe contains a cooklang defined recipe
type Recipe struct {
	Steps    []Step   // list of steps for the recipe
	Metadata Metadata // metadata of the recipe
}

func (r Recipe) String() string {
	var sb strings.Builder
	for k, v := range r.Metadata {
		sb.WriteString(fmt.Sprintf("%s %s: %s\n", metadataLinePrefix, k, v))
	}
	if len(r.Metadata) > 0 {
		sb.WriteString("\n")
	}
	for _, s := range r.Steps {
		sb.WriteString(fmt.Sprintf("%s \n\n", s.Directions))
	}
	return sb.String()
}

// ParseFile parses a cooklang recipe file and returns the recipe or an error
func ParseFile(fileName string) (*Recipe, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseStream(bufio.NewReader(f))
}

// ParseString parses a cooklang recipe string and returns the recipe or an error
func ParseString(s string) (*Recipe, error) {
	if s == "" {
		return nil, fmt.Errorf("recipe string must not be empty")
	}
	return ParseStream(strings.NewReader(s))
}

// ParseStream parses a cooklang recipe text stream and returns the recipe or an error
func ParseStream(s io.Reader) (*Recipe, error) {
	scanner := bufio.NewScanner(s)
	recipe := Recipe{
		make([]Step, 0),
		make(map[string]string),
	}
	var line string
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line = scanner.Text()

		if strings.TrimSpace(line) != "" {
			err := parseLine(line, &recipe)
			if err != nil {
				return nil, fmt.Errorf("line %d: %w", lineNumber, err)
			}
		}
	}
	return &recipe, nil
}

func parseLine(line string, recipe *Recipe) error {
	if strings.HasPrefix(line, commentsLinePrefix) {
		commentLine, err := parseSingleLineComment(line)
		if err != nil {
			return err
		}
		recipe.Steps = append(recipe.Steps, Step{Comments: []string{commentLine}})
	} else if strings.HasPrefix(line, metadataLinePrefix) {
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
	index := strings.Index(metadataLine, metadataValueSeparator)
	if index < 1 {
		return "", "", fmt.Errorf("invalid metadata: %s", metadataLine)
	}
	return strings.TrimSpace(metadataLine[:index]), strings.TrimSpace(metadataLine[index+1:]), nil
}

func peek(s string) rune {
	r, _ := utf8.DecodeRuneInString(s)
	return r
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
	var comment string
	for index, ch := range line {
		if skipIndex > index {
			continue
		}
		if ch == prefixIngredient {
			// ingredient ahead
			ingredient, skipNext, err = getIngredient(line[index:])
			if err != nil {
				return nil, err
			}
			skipIndex = index + skipNext
			step.Ingredients = append(step.Ingredients, *ingredient)
			directions.WriteString((*ingredient).Name)
			continue
		}
		if ch == prefixCookware {
			// Cookware ahead
			Cookware, skipNext, err = getCookware(line[index:])
			if err != nil {
				return nil, err
			}
			skipIndex = index + skipNext
			step.Cookware = append(step.Cookware, *Cookware)
			directions.WriteString((*Cookware).Name)
			continue
		}
		if ch == prefixTimer {
			//timer ahead
			timer, skipNext, err = getTimer(line[index:])
			if err != nil {
				return nil, err
			}
			skipIndex = index + skipNext
			step.Timers = append(step.Timers, *timer)
			directions.WriteString(fmt.Sprintf("%v %s", (*timer).Duration, (*timer).Unit))
			continue
		}
		if ch == prefixBlockComment {
			nextRune := peek(line[index+1:])
			if nextRune == '-' {
				// block comment ahead
				comment, skipNext, err = getBlockComment(line[index:])
				if err != nil {
					return nil, err
				}
				skipIndex = index + skipNext
				step.Comments = append(step.Comments, comment)
				continue
			}
		}
		// raw string
		directions.WriteRune(ch)
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

func getBlockComment(s string) (string, int, error) {
	index := strings.Index(s, "-]")
	if index == -1 {
		return "", 0, fmt.Errorf("invalid block comment")
	}
	return strings.TrimSpace(s[2:index]), index + 2, nil
}

func getFloat(s string) (bool, float64, error) {
	var fl float64
	var err error
	trimmedValue := strings.TrimSpace(s)
	if trimmedValue == "" {
		return false, 0, nil
	}
	index := strings.Index(trimmedValue, "/")
	if index == -1 {
		fl, err = strconv.ParseFloat(trimmedValue, 64)
		return err == nil, fl, err
	}
	var numerator int
	var denominator int
	numerator, err = strconv.Atoi(strings.TrimSpace(trimmedValue[:index]))
	if err != nil {
		return false, 0, err
	}

	denominator, err = strconv.Atoi(strings.TrimSpace(trimmedValue[index+1:]))
	if err != nil {
		return false, 0, err
	}
	return true, float64(numerator) / float64(denominator), nil
}

func findNodeEndIndex(line string) int {
	endIndex := -1

	for index, ch := range line {
		if index == 0 {
			continue
		}
		if (ch == prefixCookware || ch == prefixIngredient || ch == prefixTimer || ch == prefixBlockComment) && endIndex == -1 {
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
		return &IngredientAmount{Quantity: 0, QuantityRaw: "", IsNumeric: false}, nil
	}
	index := strings.Index(s, "%")
	if index == -1 {
		isNumeric, f, _ := getFloat(s)
		return &IngredientAmount{Quantity: f, QuantityRaw: strings.TrimSpace(s), IsNumeric: isNumeric}, nil
	}
	isNumeric, f, _ := getFloat(s[:index])
	return &IngredientAmount{Quantity: f, QuantityRaw: strings.TrimSpace(s[:index]), Unit: strings.TrimSpace(s[index+1:]), IsNumeric: isNumeric}, nil
}

func getCookwareFromRawString(s string) (*Cookware, error) {
	return &Cookware{strings.TrimRight(s, "{}")}, nil
}

func getTimerFromRawString(s string) (*Timer, error) {
	index := strings.Index(s, "%")
	if index == -1 {
		return nil, fmt.Errorf("invalid timer syntax: %s", s)
	}
	isNumeric, f, err := getFloat(s[:index])
	if err != nil {
		return nil, err
	}
	if !isNumeric {
		return &Timer{Duration: 0, Unit: s[index+1:]}, nil
	}
	return &Timer{Duration: f, Unit: s[index+1:]}, nil
}
