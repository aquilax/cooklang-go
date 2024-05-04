package canonical_test

import (
	"encoding/json"
	"io"
	"os"
	"slices"
	"testing"

	"github.com/aquilax/cooklang-go"
	"github.com/stretchr/testify/assert"
)

type Result struct {
	Steps [][]struct {
		Type     string      `json:"type"`
		Value    string      `json:"value,omitempty"`
		Name     string      `json:"name,omitempty"`
		Quantity interface{} `json:"quantity,omitempty"`
		Units    string      `json:"units,omitempty"`
	} `json:"steps"`
	Metadata interface{} `json:"metadata"`
}

type TestCase struct {
	Source string `json:"source"`
	Result Result `json:"result"`
}

type SpecTests struct {
	Version int                 `json:"version"`
	Tests   map[string]TestCase `json:"tests"`
}

const specFileName = "canonical.json"

func loadSpecs(fileName string) (*SpecTests, error) {
	var err error
	var jsonFile *os.File
	jsonFile, err = os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	b, _ := io.ReadAll(jsonFile)

	var result *SpecTests
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func TestCanonical(t *testing.T) {
	specs, err := loadSpecs(specFileName)
	if err != nil {
		panic(err)
	}
	skipCases := []string{}
	skipResultChecks := []string{
		"testQuantityAsText",
		"testSingleWordCookwareWithUnicodePunctuation",
		"testSingleWordCookwareWithPunctuation",
		"testIngredientNoUnits",
		"testEquipmentQuantityMultipleWords",
		"testIngredientWithEmoji",
		"testSingleWordIngredientWithUnicodePunctuation",
		"testMutipleIngredientsWithoutStopper",
		"testTimerWithUnicodeWhitespace",
		"testIngredientWithoutStopper",
		"testSingleWordIngredientWithPunctuation",
		"testSingleWordTimer",
		"testSingleWordTimerWithUnicodePunctuation",
		"testInvalidSingleWordIngredient",
		"testInvalidMultiWordIngredient",
		"testMultiWordIngredientNoAmount",
		"testEquipmentQuantityOneWord",
		"testQuantityDigitalString",
		"testCookwareWithUnicodeWhitespace",
		"testFractionsLike",
		"testIngredientWithUnicodeWhitespace",
		"testInvalidMultiWordTimer",
		"testSingleWordTimerWithPunctuation",
		"testIngredientMultipleWordsWithLeadingNumber",
		"testIngredientNoUnitsNotOnlyString",
		"testInvalidMultiWordCookware",
	}
	for name, spec := range (*specs).Tests {
		name := name
		spec := spec
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			t.Parallel()
			if slices.Contains(skipCases, name) {
				t.Skip(name)
			}
			parserV2 := cooklang.NewParserV2(&cooklang.ParseV2Config{IgnoreTypes: []cooklang.ItemType{cooklang.ItemTypeComment}})

			r, err := parserV2.ParseString(spec.Source)
			assert.NoError(err)

			if !slices.Contains(skipResultChecks, name) {
				gotJson, err := json.Marshal(r)
				assert.NoError(err)
				expectJson, err := json.Marshal(spec.Result)
				assert.NoError(err)

				assert.JSONEq(string(expectJson), string(gotJson))
			}
		})
	}
}
