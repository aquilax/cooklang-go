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
								Name:   "bacon strips",
								Amount: IngredientAmount{true, 1.0, "1", "kg"},
							},
							{
								Name:   "syrup",
								Amount: IngredientAmount{true, 1.2, "1.2", "tbsp"},
							},
						},
						Timers:   []Timer{},
						Cookware: []Cookware{},
					},
				},
				Metadata: make(Metadata),
			},
			false,
		},
		{
			"Parses recipe line with no qty",
			"Top with @1000 island dressing{ }",
			&Recipe{
				Steps: []Step{
					{
						Directions: "Top with 1000 island dressing",
						Ingredients: []Ingredient{
							{
								Name:   "1000 island dressing",
								Amount: IngredientAmount{false, 0.0, "", ""},
							},
						},
						Timers:   []Timer{},
						Cookware: []Cookware{},
					},
				},
				Metadata: make(Metadata),
			},
			false,
		},
		{
			"Parses Cookware",
			"Place the beacon on the #stove and mix with a #standing mixer{} or #fork{2}. Then use #frying pan{three} or #frying pot{two small}",
			&Recipe{
				Steps: []Step{
					{
						Directions:  "Place the beacon on the stove and mix with a standing mixer or fork. Then use frying pan or frying pot",
						Ingredients: []Ingredient{},
						Timers:      []Timer{},
						Cookware: []Cookware{
							{Name: "stove", Quantity: 1, IsNumeric: false},
							{Name: "standing mixer", Quantity: 1, IsNumeric: false},
							{Name: "fork", Quantity: 2, QuantityRaw: "2", IsNumeric: true},
							{Name: "frying pan", Quantity: 1, QuantityRaw: "three", IsNumeric: false},
							{Name: "frying pot", Quantity: 1, QuantityRaw: "two small", IsNumeric: false},
						},
					},
				},
				Metadata: make(Metadata),
			},
			false,
		},
		{
			"Parses Timers",
			"Place the beacon in the oven for ~{20%minutes}.",
			&Recipe{
				Steps: []Step{
					{
						Directions:  "Place the beacon in the oven for 20 minutes.",
						Ingredients: []Ingredient{},
						Timers:      []Timer{{"", 20.00, "minutes"}},
						Cookware:    []Cookware{},
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
							{Name: "tipo zero flour", Amount: IngredientAmount{true, 820., "820", "g"}},
							{Name: "water", Amount: IngredientAmount{true, 533, "533", "ml"}},
							{Name: "salt", Amount: IngredientAmount{true, 24.6, "24.6", "g"}},
							{Name: "fresh yeast", Amount: IngredientAmount{true, 1.6, "1.6", "g"}},
						},
						Cookware: []Cookware{{Name: "fridge", Quantity: 1, IsNumeric: false, QuantityRaw: ""}},
					},
					{
						Directions:  "Set oven to max temperature and heat pizza stone for about 40 minutes.",
						Timers:      []Timer{{Duration: 40, Unit: "minutes"}},
						Ingredients: []Ingredient{},
						Cookware: []Cookware{
							{Name: "oven", Quantity: 1, IsNumeric: false, QuantityRaw: ""},
							{Name: "pizza stone", Quantity: 1, IsNumeric: false, QuantityRaw: ""},
						},
					},
					{
						Directions: "Make some tomato sauce with chopped tomato and garlic and dried oregano. Put on a pan and leave for 15 minutes occasionally stirring.",
						Timers:     []Timer{{Duration: 15, Unit: "minutes"}},
						Ingredients: []Ingredient{
							{Name: "chopped tomato", Amount: IngredientAmount{true, 3, "3", "cans"}},
							{Name: "garlic", Amount: IngredientAmount{true, 3, "3", "cloves"}},
							{Name: "dried oregano", Amount: IngredientAmount{true, 3, "3", "tbsp"}},
						},
						Cookware: []Cookware{{Name: "pan", Quantity: 1, IsNumeric: false, QuantityRaw: ""}},
					},
					{
						Directions: "Make pizzas putting some tomato sauce with spoon on top of flattened dough. Add fresh basil, parma ham and mozzarella.",
						Timers:     []Timer{},
						Ingredients: []Ingredient{
							{Name: "fresh basil", Amount: IngredientAmount{true, 18, "18", "leaves"}},
							{Name: "parma ham", Amount: IngredientAmount{true, 3, "3", "packs"}},
							{Name: "mozzarella", Amount: IngredientAmount{true, 3, "3", "packs"}},
						},
						Cookware: []Cookware{{Name: "spoon", Quantity: 1, IsNumeric: false, QuantityRaw: ""}},
					},
					{
						Directions:  "Put in an oven for 4 minutes.",
						Timers:      []Timer{{Duration: 4, Unit: "minutes"}},
						Ingredients: []Ingredient{},
						Cookware:    []Cookware{{Name: "oven", Quantity: 1, IsNumeric: false, QuantityRaw: ""}},
					},
				},
				Metadata: Metadata{"servings": "6"},
			},
			false,
		},
		{
			"Parses block comments",
			"Text [- with block comment -] rules",
			&Recipe{
				Steps: []Step{
					{
						Directions:  "Text  rules",
						Comments:    []string{"with block comment"},
						Timers:      []Timer{},
						Ingredients: []Ingredient{},
						Cookware:    []Cookware{},
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
				t.Errorf("ParseString() = %+v, want %+v", got, tt.want)
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

func Test_getTimer(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name    string
		args    args
		want    *Timer
		wantErr bool
	}{
		{
			"Gets named timer",
			args{
				"~potato{42%minutes}",
			},
			&Timer{
				"potato",
				42,
				"minutes",
			},
			false,
		},
		{
			"Gets unn-named timer",
			args{
				"~{42%minutes}",
			},
			&Timer{
				"",
				42,
				"minutes",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := getTimer(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTimer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTimer() got = %v, want %v", got, tt.want)
			}
		})
	}
}
