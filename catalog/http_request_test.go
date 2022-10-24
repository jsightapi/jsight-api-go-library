package catalog

import (
	"fmt"
	jschemaLib "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHTTPRequestBody_Validate(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		tests := []struct {
			schema string
			types  map[string]string
			json   string
		}{
			{
				schema: `{"id": 1, "name": "Tom"}`,
				json:   `{"id": 1, "name": "Tom"}`,
			},
			{
				schema: `{"id": 1, "name": "Tom"}`,
				json:   `{"id": 222, "name": "Some string"}`,
			},
			{
				schema: `
{
	"id": 1, // {min: 1}
	"name": "Tom"
}`,
				json: `{"id": 1, "name": "Tom"}`,
			},
			{
				schema: `{"friend": @cat}`,
				types: map[string]string{
					"@cat": `{"id": 1, "name": "Tom"}`,
				},
				json: `{"friend": {"id": 333, "name": "Max"}}`,
			},
		}
		for _, tt := range tests {
			title := fmt.Sprintf("schema: %s; val: %s", tt.schema, tt.json)
			t.Run(title, func(t *testing.T) {
				b := makeHTTPRequestBody(t, []byte(tt.schema), userTypes(tt.types))
				assert.NoError(
					t,
					b.Validate([]byte(tt.json)),
					fmt.Sprintf("Validate(%s)", tt.json),
				)
			})
		}
	})

	t.Run("negative", func(t *testing.T) {
		tests := []struct {
			schema string
			types  map[string]string
			json   string
		}{
			{
				schema: `{"id": 1, "name": "Tom"}`,
				json:   `{"id": "wrong", "name": "Tom"}`,
			},
			{
				schema: `{"id": 1, "name": "Tom"}`,
				json:   `{"id": 2, "name": false}`,
			},
			{
				schema: `
{
	"id": 1, // {min: 1}
	"name": "Tom"
}`,
				json: `{"id": 0, "name": "Tom"}`,
			},
			{
				schema: `{"friend": @cat}`,
				types: map[string]string{
					"@cat": `{"id": 1, "name": "Tom"}`,
				},
				json: `{"friend": {"id": "wrong", "name": "Max"}}`,
			},
		}
		for _, tt := range tests {
			name := fmt.Sprintf("schema: %s; val: %s", tt.schema, tt.json)
			t.Run(name, func(t *testing.T) {
				b := makeHTTPRequestBody(t, []byte(tt.schema), userTypes(tt.types))
				assert.Error(
					t,
					b.Validate([]byte(tt.json)),
					fmt.Sprintf("Validate(%s)", tt.json),
				)
			})
		}
	})
}

func makeHTTPRequestBody(t *testing.T, schema []byte, userTypes *UserSchemas) HTTPRequestBody {
	es, err := NewExchangeJSightSchema("", schema, userTypes, map[string]jschemaLib.Rule{}, &UserTypes{})
	if err != nil {
		t.Fatal(err)
	}
	return HTTPRequestBody{
		Format: SerializeFormatJSON,
		Schema: es,
	}
}

func userTypes(tt map[string]string) *UserSchemas {
	uts := &UserSchemas{}
	if len(tt) != 0 {
		for name, schema := range tt {
			uts.Set(name, jschema.New(name, schema))
		}
	}
	return uts
}
