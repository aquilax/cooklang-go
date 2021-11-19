package cooklang

import "fmt"

type Parser struct{}

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

type Recipe struct {
	Steps    []Step
	Metadata map[string]string
}

func NewParser() *Parser {
	return &Parser{}
}

func ParseString(s string) (*Recipe, error) {
	if s == "" {
		return nil, fmt.Errorf("Recipe string must not be empty")
	}
	// TODO parse recipe
	return nil, nil
}
