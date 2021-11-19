package cooklang

import (
	"reflect"
	"testing"
)

func TestParseString(t *testing.T) {
	tests := []struct {
		name    string
		recipe  string
		want    *Recipe
		wantErr bool
	}{
		{
			"Empty string returns an error",
			"",
			nil,
			true,
		},
		{
			"Parses single line comments",
			"--This is a comment",
			&Recipe{
				Steps: []Step{
					{
						Comments: "This is a comment",
					},
				},
				Metadata: make(Metadata),
			},
			false,
		},
		{
			"Parses metadata",
			">> key: value",
			&Recipe{
				Steps: []Step{},
				Metadata: Metadata{
					"key": "value",
				},
			},
			false,
		},
		{
			"Parses recipe line",
			"Place @bacon strips{1%kg} on a baking sheet and glaze with @syrup{1.2%tbsp}.",
			&Recipe{
				Steps: []Step{
					{
						Directions: "Place bacon strips on a baking sheet and glaze with syrup.",
						Ingredients: []Ingredient{
							{
								Name: "bacon strips",
								Amount: IngredientAmount{
									1.0,
									"kg",
								},
							},
							{
								Name: "syrup",
								Amount: IngredientAmount{
									1.2,
									"tbsp",
								},
							},
						},
						Timers:    []Timer{},
						Equipment: []Equipment{},
					},
				},
				Metadata: make(Metadata),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseString(tt.recipe)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findIngredient(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		want     string
		endIndex int
	}{
		{
			"works with single word ingredients",
			"@word1 word2",
			"word1",
			6,
		},
		{
			"works with multiple words ingredients",
			"@word1 word2{}",
			"word1 word2{}",
			14,
		},
		{
			"works with multiple words ingredients with quantities",
			"@word1 word2{1%kg}",
			"word1 word2{1%kg}",
			18,
		},
		{
			"works when there are more then one ingredient",
			"@word1 test @word2{1%kg}",
			"word1",
			6,
		},
		{
			"works when the ingredient is at the end of the line",
			"@word1",
			"word1",
			6,
		},
		{
			"works when multi word ingredient is at the end of the line",
			"@word1{1%kg}",
			"word1{1%kg}",
			12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findNodeEndIndex('@', tt.line)
			raw := tt.line[1:got]
			if raw != tt.want {
				t.Errorf("findNodeEndIndex() got = %v, want %v", raw, tt.want)
			}
			if got != tt.endIndex {
				t.Errorf("findNodeEndIndex() got1 = %v, want %v", got, tt.endIndex)
			}
		})
	}
}
