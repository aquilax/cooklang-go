package canonical_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"
	"testing"

	"github.com/aquilax/cooklang-go"
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

	b, _ := ioutil.ReadAll(jsonFile)

	var result *SpecTests
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func contains(s []string, searchterm string) bool {
	i := sort.SearchStrings(s, searchterm)
	return i < len(s) && s[i] == searchterm
}

func compareResult(got *cooklang.Recipe, want Result) error {
	// To do check results
	return nil
}

func TestCanonical(t *testing.T) {
	specs, err := loadSpecs(specFileName)
	if err != nil {
		panic(err)
	}
	skipCases := []string{}
	sort.Strings(skipCases)
	for name, spec := range (*specs).Tests {
		name := name
		spec := spec
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if contains(skipCases, name) {
				t.Skip(name)
			}
			r, err := cooklang.ParseString(spec.Source)
			if err != nil {
				t.Errorf("%s ParseString returned %v", name, err)
			}
			if err = compareResult(r, spec.Result); err != nil {
				t.Errorf("parseString() got = %v, want %v", r, spec.Result)
			}
		})
	}
}
