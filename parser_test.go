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
