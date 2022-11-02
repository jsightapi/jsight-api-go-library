package catalog

import (
	"fmt"
	"testing"

	jschemaLib "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema"
	"github.com/stretchr/testify/assert"
)

func TestPathVariables_Validate(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		tests := []struct {
			schema    string
			userTypes map[string]string
			key       string
			value     string
		}{
			// int
			{
				schema: `{"id": 1, "name": "Tom"}`,
				key:    `id`,
				value:  `123`,
			},
			{
				schema: `{"id": 1, "name": "Tom"}`,
				key:    `id`,
				value:  `0`,
			},
			{
				schema: `{"id": 1, "name": "Tom"}`,
				key:    `id`,
				value:  `-123`,
			},

			// string
			{
				schema: `{"id": 1, "name": "Tom"}`,
				key:    `name`,
				value:  `"string"`,
			},

			// rules
			{
				schema: `{
	"id": 1, // {min: 1, max: 100}
	"name": "Tom"
}`,
				key:   `id`,
				value: `1`,
			},
			{
				schema: `{
	"id": 1, // {min: 1, max: 100}
	"name": "Tom"
}`,
				key:   `id`,
				value: `100`,
			},

			// user type
			{
				schema: `{
	"id": @uint,
	"name": "Tom"
}`,
				userTypes: map[string]string{
					"@uint": "1 // {min: 0}",
				},
				key:   `id`,
				value: `100`,
			},
		}
		for _, tt := range tests {
			name := fmt.Sprintf("schema: %s; key: %s; val: %s", tt.schema, tt.key, tt.value)
			t.Run(name, func(t *testing.T) {
				p := makePathVariables(t, tt.schema, tt.userTypes)
				assert.Nil(
					t,
					p.Validate(tt.key, tt.value),
					fmt.Sprintf("Validate(%s, %s)", tt.key, tt.value),
				)
			})
		}
	})

	t.Run("negative", func(t *testing.T) {
		tests := []struct {
			schema    string
			userTypes map[string]string
			key       string
			value     string
		}{
			// int
			{
				schema: `{"id": 1, "name": "Tom"}`,
				key:    `id`,
				value:  `"string"`,
			},
			{
				schema: `{"id": 1, "name": "Tom"}`,
				key:    `id`,
				value:  `null`,
			},
			{
				schema: `{"id": 1, "name": "Tom"}`,
				key:    `id`,
				value:  `true`,
			},
			{
				schema: `{"id": 1, "name": "Tom"}`,
				key:    `id`,
				value:  `12.34`,
			},

			// string
			{
				schema: `{"id": 1, "name": "Tom"}`,
				key:    `name`,
				value:  `123`,
			},
			{
				schema: `{"id": 1, "name": "Tom"}`,
				key:    `name`,
				value:  `12.34`,
			},
			{
				schema: `{"id": 1, "name": "Tom"}`,
				key:    `name`,
				value:  `null`,
			},
			{
				schema: `{"id": 1, "name": "Tom"}`,
				key:    `name`,
				value:  `false`,
			},
			{
				schema: `{"id": 1, "name": "Tom"}`,
				key:    `name`,
				value:  `{}`,
			},
			{
				schema: `{"id": 1, "name": "Tom"}`,
				key:    `name`,
				value:  `[]`,
			},

			// rules
			{
				schema: `{
	"id": 1, // {min: 1, max: 100}
	"name": "Tom"
}`,
				key:   `id`,
				value: `0`,
			},
			{
				schema: `{
	"id": 1, // {min: 1, max: 100}
	"name": "Tom"
}`,
				key:   `id`,
				value: `101`,
			},
			{
				schema: `{
	"id": 1, // {min: 1, max: 100}
	"name": "Tom"
}`,
				key:   `id`,
				value: `-1`,
			},

			// unknown property
			{
				schema: `{"id": 1, "name": "Tom"}`,
				key:    `unknown`,
				value:  `"string"`,
			},
		}
		for _, tt := range tests {
			name := fmt.Sprintf("schema: %s; key: %s; val: %s", tt.schema, tt.key, tt.value)
			t.Run(name, func(t *testing.T) {
				p := makePathVariables(t, tt.schema, tt.userTypes)
				assert.NotNil(
					t,
					p.Validate(tt.key, tt.value),
					fmt.Sprintf("Validate(%s, %s)", tt.key, tt.value),
				)
			})
		}
	})
}

func makePathVariables(t *testing.T, schema string, userTypes map[string]string) PathVariables {
	uss := &UserSchemas{}
	uts := &UserTypes{}
	enumRules := map[string]jschemaLib.Rule{}

	if userTypes != nil {
		for k, v := range userTypes {
			uss.Set(k, jschema.New(k, v))
		}
		for k, v := range userTypes {
			ut := &UserType{}
			var err error
			ut.Schema, err = NewExchangeJSightSchema(k, []byte(v), uss, enumRules, uts)
			if err != nil {
				t.Fatal(err)
			}
			uts.Set(k, ut)
		}
	}

	es, err := NewExchangeJSightSchema("", []byte(schema), uss, enumRules, uts)
	if err != nil {
		t.Fatal(err)
	}

	return PathVariables{
		Schema: es,
	}
}
