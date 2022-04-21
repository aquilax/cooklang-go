package cooklang_test

import (
	"encoding/json"
	"fmt"

	"github.com/aquilax/cooklang-go"
)

func ExampleParseString_toString() {
	recipeIn := `>> servings: 6

Make 6 pizza balls using @tipo zero flour{820%g}, @water{533%ml}, @salt{24.6%g} and @fresh yeast{1.6%g}. Put in a #fridge for ~{2%days}.

Set #oven to max temperature and heat #pizza stone{} for about ~{40%minutes}.

Make some tomato sauce with @chopped tomato{3%cans} and @garlic{3%cloves} and @dried oregano{3%tbsp}. Put on a #pan and leave for ~{15%minutes} occasionally stirring.

Make pizzas putting some tomato sauce with #spoon on top of flattened dough. Add @fresh basil{18%leaves}, @parma ham{3%packs} and @mozzarella{3%packs}.

Put in an #oven for ~{4%minutes}.`
	r, _ := cooklang.ParseString(recipeIn)
	fmt.Print(r)
	// Output:
	// >> servings: 6
	//
	// Make 6 pizza balls using tipo zero flour, water, salt and fresh yeast. Put in a fridge for 2 days.
	//
	// Set oven to max temperature and heat pizza stone for about 40 minutes.
	//
	// Make some tomato sauce with chopped tomato and garlic and dried oregano. Put on a pan and leave for 15 minutes occasionally stirring.
	//
	// Make pizzas putting some tomato sauce with spoon on top of flattened dough. Add fresh basil, parma ham and mozzarella.
	//
	// Put in an oven for 4 minutes.
}

func ExampleParseString() {
	recipe := `>> servings: 6

Make 6 pizza balls using @tipo zero flour{820%g}, @water{533%ml}, @salt{24.6%g} and @fresh yeast{1.6%g}. Put in a #fridge for ~{2%days}.

Set #oven to max temperature and heat #pizza stone{} for about ~{40%minutes}.

Make some tomato sauce with @chopped tomato{3%cans} and @garlic{3%cloves} and @dried oregano{3%tbsp}. Put on a #pan and leave for ~{15%minutes} occasionally stirring.

Make pizzas putting some tomato sauce with #spoon on top of flattened dough. Add @fresh basil{18%leaves}, @parma ham{3%packs} and @mozzarella{3%packs}.

Put in an #oven for ~{4%minutes}.`
	r, _ := cooklang.ParseString(recipe)
	j, _ := json.MarshalIndent(r, "", "  ")
	fmt.Println(string(j))
	// Output:
	// {
	//   "Steps": [
	//     {
	//       "Directions": "Make 6 pizza balls using tipo zero flour, water, salt and fresh yeast. Put in a fridge for 2 days.",
	//       "Timers": [
	//         {
	//           "Name": "",
	//           "Duration": 2,
	//           "Unit": "days"
	//         }
	//       ],
	//       "Ingredients": [
	//         {
	//           "Name": "tipo zero flour",
	//           "Amount": {
	//             "IsNumeric": true,
	//             "Quantity": 820,
	//             "QuantityRaw": "820",
	//             "Unit": "g"
	//           }
	//         },
	//         {
	//           "Name": "water",
	//           "Amount": {
	//             "IsNumeric": true,
	//             "Quantity": 533,
	//             "QuantityRaw": "533",
	//             "Unit": "ml"
	//           }
	//         },
	//         {
	//           "Name": "salt",
	//           "Amount": {
	//             "IsNumeric": true,
	//             "Quantity": 24.6,
	//             "QuantityRaw": "24.6",
	//             "Unit": "g"
	//           }
	//         },
	//         {
	//           "Name": "fresh yeast",
	//           "Amount": {
	//             "IsNumeric": true,
	//             "Quantity": 1.6,
	//             "QuantityRaw": "1.6",
	//             "Unit": "g"
	//           }
	//         }
	//       ],
	//       "Cookware": [
	//         {
	//           "IsNumeric": false,
	//           "Name": "fridge",
	//           "Quantity": 1,
	//           "QuantityRaw": ""
	//         }
	//       ],
	//       "Comments": null
	//     },
	//     {
	//       "Directions": "Set oven to max temperature and heat pizza stone for about 40 minutes.",
	//       "Timers": [
	//         {
	//           "Name": "",
	//           "Duration": 40,
	//           "Unit": "minutes"
	//         }
	//       ],
	//       "Ingredients": [],
	//       "Cookware": [
	//         {
	//           "IsNumeric": false,
	//           "Name": "oven",
	//           "Quantity": 1,
	//           "QuantityRaw": ""
	//         },
	//         {
	//           "IsNumeric": false,
	//           "Name": "pizza stone",
	//           "Quantity": 1,
	//           "QuantityRaw": ""
	//         }
	//       ],
	//       "Comments": null
	//     },
	//     {
	//       "Directions": "Make some tomato sauce with chopped tomato and garlic and dried oregano. Put on a pan and leave for 15 minutes occasionally stirring.",
	//       "Timers": [
	//         {
	//           "Name": "",
	//           "Duration": 15,
	//           "Unit": "minutes"
	//         }
	//       ],
	//       "Ingredients": [
	//         {
	//           "Name": "chopped tomato",
	//           "Amount": {
	//             "IsNumeric": true,
	//             "Quantity": 3,
	//             "QuantityRaw": "3",
	//             "Unit": "cans"
	//           }
	//         },
	//         {
	//           "Name": "garlic",
	//           "Amount": {
	//             "IsNumeric": true,
	//             "Quantity": 3,
	//             "QuantityRaw": "3",
	//             "Unit": "cloves"
	//           }
	//         },
	//         {
	//           "Name": "dried oregano",
	//           "Amount": {
	//             "IsNumeric": true,
	//             "Quantity": 3,
	//             "QuantityRaw": "3",
	//             "Unit": "tbsp"
	//           }
	//         }
	//       ],
	//       "Cookware": [
	//         {
	//           "IsNumeric": false,
	//           "Name": "pan",
	//           "Quantity": 1,
	//           "QuantityRaw": ""
	//         }
	//       ],
	//       "Comments": null
	//     },
	//     {
	//       "Directions": "Make pizzas putting some tomato sauce with spoon on top of flattened dough. Add fresh basil, parma ham and mozzarella.",
	//       "Timers": [],
	//       "Ingredients": [
	//         {
	//           "Name": "fresh basil",
	//           "Amount": {
	//             "IsNumeric": true,
	//             "Quantity": 18,
	//             "QuantityRaw": "18",
	//             "Unit": "leaves"
	//           }
	//         },
	//         {
	//           "Name": "parma ham",
	//           "Amount": {
	//             "IsNumeric": true,
	//             "Quantity": 3,
	//             "QuantityRaw": "3",
	//             "Unit": "packs"
	//           }
	//         },
	//         {
	//           "Name": "mozzarella",
	//           "Amount": {
	//             "IsNumeric": true,
	//             "Quantity": 3,
	//             "QuantityRaw": "3",
	//             "Unit": "packs"
	//           }
	//         }
	//       ],
	//       "Cookware": [
	//         {
	//           "IsNumeric": false,
	//           "Name": "spoon",
	//           "Quantity": 1,
	//           "QuantityRaw": ""
	//         }
	//       ],
	//       "Comments": null
	//     },
	//     {
	//       "Directions": "Put in an oven for 4 minutes.",
	//       "Timers": [
	//         {
	//           "Name": "",
	//           "Duration": 4,
	//           "Unit": "minutes"
	//         }
	//       ],
	//       "Ingredients": [],
	//       "Cookware": [
	//         {
	//           "IsNumeric": false,
	//           "Name": "oven",
	//           "Quantity": 1,
	//           "QuantityRaw": ""
	//         }
	//       ],
	//       "Comments": null
	//     }
	//   ],
	//   "Metadata": {
	//     "servings": "6"
	//   }
	// }
}
