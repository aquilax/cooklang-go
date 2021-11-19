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
						Comments: []string{"This is a comment"},
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
		{
			"Parses equipment",
			"Place the beacon on the #stove and mix with a #standing mixer{}.",
			&Recipe{
				Steps: []Step{
					{
						Directions:  "Place the beacon on the stove and mix with a standing mixer.",
						Ingredients: []Ingredient{},
						Timers:      []Timer{},
						Equipment: []Equipment{
							{Name: "stove"},
							{Name: "standing mixer"},
						},
					},
				},
				Metadata: make(Metadata),
			},
			false,
		},
		{
			"Parses timers",
			"Place the beacon in the oven for ~{20%minutes}.",
			&Recipe{
				Steps: []Step{
					{
						Directions:  "Place the beacon in the oven for 20 minutes.",
						Ingredients: []Ingredient{},
						Timers: []Timer{
							{
								20.00,
								"minutes",
							},
						},
						Equipment: []Equipment{},
					},
				},
				Metadata: make(Metadata),
			},
			false,
		},
		{
			"Full recipe",
			`>> servings: 6

Make 6 pizza balls using @tipo zero flour{820%g}, @water{533%ml}, @salt{24.6%g} and @fresh yeast{1.6%g}. Put in a #fridge for ~{2%days}.

Set #oven to max temperature and heat #pizza stone{} for about ~{40%minutes}.

Make some tomato sauce with @chopped tomato{3%cans} and @garlic{3%cloves} and @dried oregano{3%tbsp}. Put on a #pan and leave for ~{15%minutes} occasionally stirring.

Make pizzas putting some tomato sauce with #spoon on top of flattened dough. Add @fresh basil{18%leaves}, @parma ham{3%packs} and @mozzarella{3%packs}.

Put in an #oven for ~{4%minutes}.`,
			&Recipe{
				Steps: []Step{
					{
						Directions: "Make 6 pizza balls using tipo zero flour, water, salt and fresh yeast. Put in a fridge for 2 days.",
						Timers:     []Timer{{Duration: 2, Unit: "days"}},
						Ingredients: []Ingredient{
							{Name: "tipo zero flour", Amount: IngredientAmount{820., "g"}},
							{Name: "water", Amount: IngredientAmount{533, "ml"}},
							{Name: "salt", Amount: IngredientAmount{24.6, "g"}},
							{Name: "fresh yeast", Amount: IngredientAmount{1.6, "g"}},
						},
						Equipment: []Equipment{{Name: "fridge"}},
					},
					{
						Directions:  "Set oven to max temperature and heat pizza stone for about 40 minutes.",
						Timers:      []Timer{{Duration: 40, Unit: "minutes"}},
						Ingredients: []Ingredient{},
						Equipment:   []Equipment{{Name: "oven"}, {Name: "pizza stone"}},
					},
					{
						Directions: "Make some tomato sauce with chopped tomato and garlic and dried oregano. Put on a pan and leave for 15 minutes occasionally stirring.",
						Timers:     []Timer{{Duration: 15, Unit: "minutes"}},
						Ingredients: []Ingredient{
							{Name: "chopped tomato", Amount: IngredientAmount{3, "cans"}},
							{Name: "garlic", Amount: IngredientAmount{3, "cloves"}},
							{Name: "dried oregano", Amount: IngredientAmount{3, "tbsp"}},
						},
						Equipment: []Equipment{{Name: "pan"}},
					},
					{
						Directions: "Make pizzas putting some tomato sauce with spoon on top of flattened dough. Add fresh basil, parma ham and mozzarella.",
						Timers:     []Timer{},
						Ingredients: []Ingredient{
							{Name: "fresh basil", Amount: IngredientAmount{18, "leaves"}},
							{Name: "parma ham", Amount: IngredientAmount{3, "packs"}},
							{Name: "mozzarella", Amount: IngredientAmount{3, "packs"}},
						},
						Equipment: []Equipment{{Name: "spoon"}},
					},
					{
						Directions:  "Put in an oven for 4 minutes.",
						Timers:      []Timer{{Duration: 4, Unit: "minutes"}},
						Ingredients: []Ingredient{},
						Equipment:   []Equipment{{Name: "oven"}},
					},
				},
				Metadata: Metadata{"servings": "6"},
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
			got := findNodeEndIndex(tt.line)
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
