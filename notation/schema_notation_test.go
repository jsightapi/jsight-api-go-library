package notation

import (
	"testing"
)

func TestNewSchemaNotation(t *testing.T) {
	tests := []struct {
		name    string
		want    SchemaNotation
		wantErr bool
	}{
		{
			"",
			SchemaNotationJSight,
			false,
		},

		{
			"jsight",
			SchemaNotationJSight,
			false,
		},
		{
			"regex",
			SchemaNotationRegex,
			false,
		},
		{
			"any",
			SchemaNotationAny,
			false,
		},
		{
			"empty",
			SchemaNotationEmpty,
			false,
		},
		{
			"unknown",
			"",
			true,
		},
		{
			"unknown 2",
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSchemaNotation(tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSchemaNotation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewSchemaNotation() got = %v, want %v", got, tt.want)
			}
		})
	}
}
