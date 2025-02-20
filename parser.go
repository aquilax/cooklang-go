// Package cooklang provides a parser for .cook defined recipes as defined in
// https://cooklang.org/docs/spec/
package cooklang

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"gopkg.in/yaml.v3"
)

const (
	commentsLinePrefix     = "--"
	metadataLinePrefix     = ">>"
	metadataValueSeparator = ":"
	prefixIngredient       = '@'
	prefixCookware         = '#'
	prefixTimer            = '~'
	prefixBlockComment     = '['
	prefixInlineComment    = '-'

	ItemTypeText       ItemType = "text"
	ItemTypeComment    ItemType = "comment"
	ItemTypeCookware   ItemType = "cookware"
	ItemTypeIngredient ItemType = "ingredient"
	ItemTypeTimer      ItemType = "timer"

	CommentTypeLine    CommentType = 1
	CommentTypeBlock   CommentType = 2
	CommentTypeEndLine CommentType = 3
)

type ItemType string

// CommentType defines what type is the comment
type CommentType int

// Cookware represents a cookware item
type Cookware struct {
	IsNumeric   bool    // true if the amount is numeric
	Name        string  // cookware name
	Quantity    float64 // quantity of the cookware
	QuantityRaw string  // quantity of the cookware as raw text
}

type CookwareV2 struct {
	Type     ItemType `json:"type"`
	Name     string   `json:"name"`
	Quantity float64  `json:"quantity"`
}

func (c Cookware) asCookwareV2() CookwareV2 {
	return CookwareV2{
		Type:     ItemTypeCookware,
		Name:     c.Name,
		Quantity: c.Quantity,
	}
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

type IngredientV2 struct {
	Type     ItemType `json:"type"`
	Name     string   `json:"name"`
	Quantity float64  `json:"quantity"`
	Units    string   `json:"units,omitempty"`
}

func (i Ingredient) asIngredientV2() IngredientV2 {
	return IngredientV2{
		Type:     ItemTypeIngredient,
		Name:     i.Name,
		Quantity: i.Amount.Quantity,
		Units:    i.Amount.Unit,
	}
}

// Timer represents a time duration
type Timer struct {
	Name     string  // name of the timer
	Duration float64 // duration of the timer
	Unit     string  // time unit of the duration
}

type TimerV2 struct {
	Type     ItemType `json:"type"`
	Name     string   `json:"name,omitempty"`
	Quantity float64  `json:"quantity"`
	Unit     string   `json:"units"`
}

func (t Timer) asTimerV2() TimerV2 {
	return TimerV2{
		Type:     ItemTypeTimer,
		Name:     t.Name,
		Quantity: t.Duration,
		Unit:     t.Unit,
	}
}

// Comment represents comment text
type Comment struct {
	Type  CommentType
	Value string
}

type Text struct {
	Value string
}

type TextV2 struct {
	Type  ItemType `json:"type"`
	Value string   `json:"value"`
}

func (t Text) asTextV2() TextV2 {
	return TextV2{ItemTypeText, t.Value}
}

type jsonStep struct {
	Type     string `json:"type"`
	Value    string `json:"value,omitempty"`
	Name     string `json:"name,omitempty"`
	Quantity any    `json:"quantity,omitempty"`
	Units    string `json:"units,omitempty"`
}

func (t *Text) MarshalJson() ([]byte, error) {
	return json.Marshal(&jsonStep{
		Type:  "text",
		Value: t.Value,
	})
}

func newText(v string) Text {
	return Text{v}
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
type Metadata = map[string]any

// Recipe contains a cooklang defined recipe
type Recipe struct {
	Steps    []Step   // list of steps for the recipe
	Metadata Metadata // metadata of the recipe
}

type ParseV2Config struct {
	IgnoreTypes []ItemType
}

type StepV2 []any

// RecipeV2 contains a cooklang defined recipe
type RecipeV2 struct {
	Steps    []StepV2 `json:"steps"`    // list of steps for the recipe
	Metadata Metadata `json:"metadata"` // metadata of the recipe
}

type ParserV2 struct {
	config        *ParseV2Config
	inFrontMatter bool
	frontMatter   string
}

func (r Recipe) String() string {
	var sb strings.Builder
	for k, v := range r.Metadata {
		sb.WriteString(fmt.Sprintf("%s %s: %s\n", metadataLinePrefix, k, v))
	}
	if len(r.Metadata) > 0 {
		sb.WriteString("\n")
	}
	steps := len(r.Steps)
	for i, s := range r.Steps {
		sb.WriteString(fmt.Sprintln(s.Directions))
		if i != steps-1 {
			sb.WriteString("\n")
		}
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

func (p *ParserV2) ParseFile(fileName string) (*RecipeV2, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return p.ParseStream(bufio.NewReader(f))
}

// ParseString parses a cooklang recipe string and returns the recipe or an error
func ParseString(s string) (*Recipe, error) {
	if s == "" {
		return nil, fmt.Errorf("recipe string must not be empty")
	}
	return ParseStream(strings.NewReader(s))
}

func (p *ParserV2) ParseString(s string) (*RecipeV2, error) {
	if s == "" {
		return nil, fmt.Errorf("recipe string must not be empty")
	}
	return p.ParseStream(strings.NewReader(s))
}

func NewParserV2(config *ParseV2Config) *ParserV2 {
	return &ParserV2{
		config: config,
	}
}

// ParseStream parses a cooklang recipe text stream and returns the recipe or an error
func ParseStream(s io.Reader) (*Recipe, error) {
	scanner := bufio.NewScanner(s)
	recipe := Recipe{
		make([]Step, 0),
		make(map[string]any),
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

// ParseStream parses a cooklang recipe text stream and returns the recipe or an error
func (p *ParserV2) ParseStream(s io.Reader) (*RecipeV2, error) {
	scanner := bufio.NewScanner(s)
	recipe := RecipeV2{
		make([]StepV2, 0),
		make(map[string]any),
	}
	var line string
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line = scanner.Text()

		if strings.TrimSpace(line) != "" {
			err := p.parseLine(line, &recipe)
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
		recipe.Steps = append(recipe.Steps, Step{
			Comments: []string{commentLine},
		})
	} else if strings.HasPrefix(line, metadataLinePrefix) {
		key, value, err := parseMetadata(line)
		if err != nil {
			return err
		}
		recipe.Metadata[key] = value
	} else {
		step, err := parseRecipeLine(line)
		if err != nil {
			return err
		}
		recipe.Steps = append(recipe.Steps, *step)
	}
	return nil
}

func (p *ParserV2) parseLine(line string, recipe *RecipeV2) error {
	line = strings.TrimRight(line, " ")

	if line == "---" && !p.inFrontMatter {
		p.inFrontMatter = true
	} else if line == "---" && p.inFrontMatter {
		p.inFrontMatter = false
		y := strings.NewReader(p.frontMatter)
		err := yaml.NewDecoder(y).Decode(recipe.Metadata)
		if err != nil {
			return fmt.Errorf("decoding yaml front matter: %w", err)
		}
	} else if p.inFrontMatter {
		p.frontMatter = p.frontMatter + line + "\n"
	} else if strings.HasPrefix(line, commentsLinePrefix) {
		commentLine, err := parseSingleLineComment(line)
		if err != nil {
			return err
		}
		if !slices.Contains(p.config.IgnoreTypes, ItemTypeComment) {
			recipe.Steps = append(recipe.Steps, StepV2{Comment{CommentTypeLine, commentLine}})
		}
	} else if strings.HasPrefix(line, metadataLinePrefix) {
		key, value, err := parseMetadata(line)
		if err != nil {
			return err
		}
		recipe.Metadata[key] = value
	} else {
		step, err := p.parseRecipeLine(line)
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

func parseStepCB(line string, cb func(item any) (bool, error)) (string, error) {
	skipIndex := -1
	var directions strings.Builder
	var err error
	var skipNext int
	var ingredient *Ingredient
	var cookware *Cookware
	var timer *Timer
	var comment string
	var buffer strings.Builder
	for index, ch := range line {
		if skipIndex > index {
			continue
		}
		if ch == prefixIngredient {
			nextRune := peek(line[index+1:])
			if nextRune != ' ' {
				if buffer.Len() > 0 {
					if stop, err := cb(newText(buffer.String())); err != nil || stop {
						return directions.String(), err
					}
					buffer.Reset()
				}
				// ingredient ahead
				ingredient, skipNext, err = getIngredient(line[index:])
				if err != nil {
					return directions.String(), err
				}
				skipIndex = index + skipNext
				directions.WriteString((*ingredient).Name)
				if stop, err := cb(*ingredient); err != nil || stop {
					return directions.String(), err
				}
				continue

			}
		}
		if ch == prefixCookware {
			nextRune := peek(line[index+1:])
			if nextRune != ' ' {
				if buffer.Len() > 0 {
					if stop, err := cb(newText(buffer.String())); err != nil || stop {
						return directions.String(), err
					}
					buffer.Reset()
				}
				// Cookware ahead
				cookware, skipNext, err = getCookware(line[index:])
				if err != nil {
					return directions.String(), err
				}
				skipIndex = index + skipNext
				directions.WriteString((*cookware).Name)
				if stop, err := cb(*cookware); err != nil || stop {
					return directions.String(), err
				}
				continue
			}
		}
		if ch == prefixTimer {
			nextRune := peek(line[index+1:])
			if nextRune != ' ' {
				if buffer.Len() > 0 {
					if stop, err := cb(newText(buffer.String())); err != nil || stop {
						return directions.String(), err
					}
					buffer.Reset()
				}
				//timer ahead
				timer, skipNext, err = getTimer(line[index:])
				if err != nil {
					return directions.String(), err
				}
				skipIndex = index + skipNext
				directions.WriteString(fmt.Sprintf("%v %s", (*timer).Duration, (*timer).Unit))
				if stop, err := cb(*timer); err != nil || stop {
					return directions.String(), err
				}
				continue
			}
		}
		if ch == prefixBlockComment {
			nextRune := peek(line[index+1:])
			if nextRune == '-' {
				if buffer.Len() > 0 {
					if stop, err := cb(newText(buffer.String())); err != nil || stop {
						return directions.String(), err
					}
					buffer.Reset()
				}
				// block comment ahead
				comment, skipNext, err = getBlockComment(line[index:])
				if err != nil {
					return directions.String(), err
				}
				skipIndex = index + skipNext
				if stop, err := cb(Comment{CommentTypeBlock, comment}); err != nil || stop {
					return directions.String(), err
				}
				continue
			}
		}
		if ch == prefixInlineComment {
			nextRune := peek(line[index+1:])
			if nextRune == prefixInlineComment {
				if buffer.Len() > 0 {
					if stop, err := cb(newText(buffer.String())); err != nil || stop {
						return directions.String(), err
					}
					buffer.Reset()
				}
				// end-line comment ahead
				comment = strings.TrimSpace(line[index+len(commentsLinePrefix):])
				if err != nil {
					return directions.String(), err
				}
				if stop, err := cb(Comment{CommentTypeEndLine, comment}); err != nil || stop {
					return directions.String(), err
				}
				break
			}
		}
		// raw string
		buffer.WriteRune(ch)
		directions.WriteRune(ch)
	}
	if buffer.Len() > 0 {
		if stop, err := cb(newText(buffer.String())); err != nil || stop {
			return directions.String(), err
		}
		buffer.Reset()
	}
	return strings.TrimSpace(directions.String()), nil
}

func parseRecipeLine(line string) (*Step, error) {
	step := Step{
		Timers:      make([]Timer, 0),
		Ingredients: make([]Ingredient, 0),
		Cookware:    make([]Cookware, 0),
	}
	var err error
	step.Directions, err = parseStepCB(line, func(item any) (bool, error) {
		switch v := item.(type) {
		case Timer:
			step.Timers = append(step.Timers, v)
		case Ingredient:
			step.Ingredients = append(step.Ingredients, v)
		case Cookware:
			step.Cookware = append(step.Cookware, v)
		case Text:
			//
		case Comment:
			step.Comments = append(step.Comments, v.Value)
		default:
			return true, fmt.Errorf("unknown type %T", v)
		}
		return false, nil
	})
	if err != nil {
		return nil, err
	}
	return &step, nil
}

func (p *ParserV2) parseRecipeLine(line string) (*StepV2, error) {
	step := StepV2{}
	var err error
	_, err = parseStepCB(line, func(item any) (bool, error) {
		switch v := item.(type) {
		case Timer:
			if !slices.Contains(p.config.IgnoreTypes, ItemTypeTimer) {
				step = append(step, v.asTimerV2())
			}
		case Ingredient:
			if !slices.Contains(p.config.IgnoreTypes, ItemTypeIngredient) {
				step = append(step, v.asIngredientV2())
			}
		case Cookware:
			if !slices.Contains(p.config.IgnoreTypes, ItemTypeCookware) {
				step = append(step, v.asCookwareV2())
			}
		case Text:
			if !slices.Contains(p.config.IgnoreTypes, ItemTypeText) {
				step = append(step, v.asTextV2())
			}
		case Comment:
			if !slices.Contains(p.config.IgnoreTypes, ItemTypeComment) {
				step = append(step, v)
			}
		default:
			return true, fmt.Errorf("unknown type %T", v)
		}
		return false, nil
	})
	if err != nil {
		return nil, err
	}
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
	timer, err := getTimerFromRawString(line[1:endIndex])
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
	amount, err := getAmount(s[index+1:len(s)-1], 0)
	if err != nil {
		return nil, err
	}
	return &Ingredient{Name: s[:index], Amount: *amount}, nil
}

func getAmount(s string, defaultValue float64) (*IngredientAmount, error) {
	if s == "" {
		return &IngredientAmount{Quantity: defaultValue, QuantityRaw: "", IsNumeric: false}, nil
	}
	index := strings.Index(s, "%")
	if index == -1 {
		isNumeric, f, _ := getFloat(s)
		if !isNumeric {
			f = defaultValue
		}
		return &IngredientAmount{Quantity: f, QuantityRaw: strings.TrimSpace(s), IsNumeric: isNumeric}, nil
	}
	isNumeric, f, _ := getFloat(s[:index])
	if !isNumeric {
		f = defaultValue
	}
	return &IngredientAmount{Quantity: f, QuantityRaw: strings.TrimSpace(s[:index]), Unit: strings.TrimSpace(s[index+1:]), IsNumeric: isNumeric}, nil
}

func getCookwareFromRawString(s string) (*Cookware, error) {
	index := strings.Index(s, "{")
	if index == -1 {
		return &Cookware{Name: s, Quantity: 1}, nil
	}
	amount, err := getAmount(s[index+1:len(s)-1], 1)
	if err != nil {
		return nil, err
	}
	return &Cookware{Name: s[:index], Quantity: amount.Quantity, IsNumeric: amount.IsNumeric, QuantityRaw: amount.QuantityRaw}, nil
}

func getTimerFromRawString(s string) (*Timer, error) {
	name := ""
	index := strings.Index(s, "{")
	if index > -1 {
		name = strings.TrimSpace(s[:index])
		s = s[index+1:]
	}
	index = strings.Index(s, "%")
	if index == -1 {
		return &Timer{Name: s, Duration: 0, Unit: ""}, nil
	}
	isNumeric, f, err := getFloat(s[:index])
	if err != nil {
		return nil, err
	}
	if !isNumeric {
		return &Timer{Name: name, Duration: 0, Unit: s[index+1 : len(s)-1]}, nil
	}
	return &Timer{Name: name, Duration: f, Unit: s[index+1 : len(s)-1]}, nil
}
