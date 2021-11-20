# cooklang-go [![Go Reference](https://pkg.go.dev/badge/github.com/aquilax/cooklang-go.svg)](https://pkg.go.dev/github.com/aquilax/cooklang-go)

[Cooklang](https://cooklang.org/) parser in Go

## Usage

```go
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
	//           "Duration": 2,
	//           "Unit": "days"
	//         }
	//       ],
	//       "Ingredients": [
	//         {
	//           "Name": "tipo zero flour",
	//           "Amount": {
	//             "Quantity": 820,
	//             "Unit": "g"
	//           }
	//         },
	//         {
	//           "Name": "water",
	//           "Amount": {
	//             "Quantity": 533,
	//             "Unit": "ml"
	//           }
	//         },
	//         {
	//           "Name": "salt",
	//           "Amount": {
	//             "Quantity": 24.6,
	//             "Unit": "g"
	//           }
	//         },
	//         {
	//           "Name": "fresh yeast",
	//           "Amount": {
	//             "Quantity": 1.6,
	//             "Unit": "g"
	//           }
	//         }
	//       ],
	//       "Cookware": [
	//         {
	//           "Name": "fridge"
	//         }
	//       ],
	//       "Comments": null
	//     },
	//     {
	//       "Directions": "Set oven to max temperature and heat pizza stone for about 40 minutes.",
	//       "Timers": [
	//         {
	//           "Duration": 40,
	//           "Unit": "minutes"
	//         }
	//       ],
	//       "Ingredients": [],
	//       "Cookware": [
	//         {
	//           "Name": "oven"
	//         },
	//         {
	//           "Name": "pizza stone"
	//         }
	//       ],
	//       "Comments": null
	//     },
	//     {
	//       "Directions": "Make some tomato sauce with chopped tomato and garlic and dried oregano. Put on a pan and leave for 15 minutes occasionally stirring.",
	//       "Timers": [
	//         {
	//           "Duration": 15,
	//           "Unit": "minutes"
	//         }
	//       ],
	//       "Ingredients": [
	//         {
	//           "Name": "chopped tomato",
	//           "Amount": {
	//             "Quantity": 3,
	//             "Unit": "cans"
	//           }
	//         },
	//         {
	//           "Name": "garlic",
	//           "Amount": {
	//             "Quantity": 3,
	//             "Unit": "cloves"
	//           }
	//         },
	//         {
	//           "Name": "dried oregano",
	//           "Amount": {
	//             "Quantity": 3,
	//             "Unit": "tbsp"
	//           }
	//         }
	//       ],
	//       "Cookware": [
	//         {
	//           "Name": "pan"
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
	//             "Quantity": 18,
	//             "Unit": "leaves"
	//           }
	//         },
	//         {
	//           "Name": "parma ham",
	//           "Amount": {
	//             "Quantity": 3,
	//             "Unit": "packs"
	//           }
	//         },
	//         {
	//           "Name": "mozzarella",
	//           "Amount": {
	//             "Quantity": 3,
	//             "Unit": "packs"
	//           }
	//         }
	//       ],
	//       "Cookware": [
	//         {
	//           "Name": "spoon"
	//         }
	//       ],
	//       "Comments": null
	//     },
	//     {
	//       "Directions": "Put in an oven for 4 minutes.",
	//       "Timers": [
	//         {
	//           "Duration": 4,
	//           "Unit": "minutes"
	//         }
	//       ],
	//       "Ingredients": [],
	//       "Cookware": [
	//         {
	//           "Name": "oven"
	//         }
	//       ],
	//       "Comments": null
	//     }
	//   ],
	//   "Metadata": {
	//     "servings": "6"
	//   }
	// }
```
