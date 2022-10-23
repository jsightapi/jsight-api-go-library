package catalog

import (
	"fmt"
	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPathVariables_Validate(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		tests := []struct {
			schema string
			key    string
			value  string
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
		}
		for _, tt := range tests {
			name := fmt.Sprintf("schema: %s; key: %s, val: %s", tt.schema, tt.key, tt.value)
			t.Run(name, func(t *testing.T) {
				p := makePathVariables(t, []byte(tt.schema))
				assert.NoError(
					t,
					p.Validate([]byte(tt.key), []byte(tt.value)),
					fmt.Sprintf("Validate(%s, %s)", tt.key, tt.value),
				)
			})
		}
	})

	t.Run("negative", func(t *testing.T) {
		tests := []struct {
			schema string
			key    string
			value  string
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
			name := fmt.Sprintf("schema: %s; key: %s, val: %s", tt.schema, tt.key, tt.value)
			t.Run(name, func(t *testing.T) {
				p := makePathVariables(t, []byte(tt.schema))
				assert.Error(
					t,
					p.Validate([]byte(tt.key), []byte(tt.value)),
					fmt.Sprintf("Validate(%s, %s)", tt.key, tt.value),
				)
			})
		}
	})
}

func makePathVariables(t *testing.T, b []byte) PathVariables {
	es, err := NewExchangeJSightSchema("", b, &UserSchemas{}, map[string]jschema.Rule{}, &UserTypes{})
	if err != nil {
		t.Fatal(err)
	}
	return PathVariables{
		Schema: es,
	}
}
