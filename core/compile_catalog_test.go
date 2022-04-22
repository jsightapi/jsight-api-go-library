package core

import (
	"j/japi/catalog"
	"j/japi/directive"
	"j/schema/bytes"
	"j/schema/fs"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJApiCore_BuildResourceMethodsPathVariables(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		cc := map[string]struct {
			given    *JApiCore
			expected *JApiCore
		}{
			"empty": {
				&JApiCore{
					catalog: &catalog.Catalog{
						ResourceMethods: &catalog.ResourceMethods{},
					},
				},
				&JApiCore{
					catalog: &catalog.Catalog{
						ResourceMethods: &catalog.ResourceMethods{},
					},
				},
			},
		}

		for n, c := range cc {
			t.Run(n, func(t *testing.T) {
				err := c.given.BuildResourceMethodsPathVariables()
				require.Nil(t, err)

				assert.Equal(t, c.expected, c.given)
			})
		}
	})

	t.Run("negative", func(t *testing.T) {
		cc := map[string]*JApiCore{
			`Has unused parameters "bar" in schema`: {
				rawPathVariables: []rawPathVariable{
					{
						schema: catalog.Schema{
							ContentJSight: &catalog.SchemaContentJSight{
								Properties: catalog.NewProperties(
									map[string]*catalog.SchemaContentJSight{
										"foo": nil,
										"bar": nil,
									},
									[]string{"foo", "bar"},
								),
							},
						},
						parameters: []PathParameter{
							{
								path:      "/path/{foo}/{bar}",
								parameter: "foo",
							},
						},
						pathDirective: *directive.New(
							directive.Get,
							directive.NewCoords(
								fs.NewFile("foo", bytes.Bytes("/path/{foo}/{bar}")),
								0,
								5,
							),
						),
					},
				},

				catalog: &catalog.Catalog{
					ResourceMethods: &catalog.ResourceMethods{},
				},
			},
		}

		for expected, core := range cc {
			t.Run(expected, func(t *testing.T) {
				err := core.BuildResourceMethodsPathVariables()
				assert.EqualError(t, err, expected)
			})
		}

		t.Run("nil", func(t *testing.T) {
			assert.Panics(t, func() {
				var core *JApiCore

				_ = core.BuildResourceMethodsPathVariables()
			})
		})
	})
}

func TestCore_propertiesToMap(t *testing.T) {
	cc := map[string]struct {
		given    *catalog.Properties
		expected map[string]*catalog.SchemaContentJSight
	}{
		"nil": {},

		"empty": {
			given: &catalog.Properties{},
		},

		"with data": {
			given: catalog.NewProperties(
				map[string]*catalog.SchemaContentJSight{
					"foo": nil,
					"bar": {
						Note: "fake",
					},
				},
				[]string{"foo", "bar"},
			),
			expected: map[string]*catalog.SchemaContentJSight{
				"foo": nil,
				"bar": {
					Note: "fake",
				},
			},
		},
	}

	for n, c := range cc {
		t.Run(n, func(t *testing.T) {
			actual := (&JApiCore{}).propertiesToMap(c.given)
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestCore_getPropertiesNames(t *testing.T) {
	cc := map[string]struct {
		given    map[string]*catalog.SchemaContentJSight
		expected string
	}{
		"nil": {},

		"empty": {
			given: map[string]*catalog.SchemaContentJSight{},
		},

		"with data": {
			given: map[string]*catalog.SchemaContentJSight{
				"foo":  nil,
				"bar":  nil,
				"fizz": nil,
				"buzz": nil,
			},
			expected: "foo, bar, fizz, buzz",
		},
	}

	split := func(s string) []string {
		return strings.Split(s, ", ")
	}

	for n, c := range cc {
		t.Run(n, func(t *testing.T) {
			actual := (&JApiCore{}).getPropertiesNames(c.given)
			assert.ElementsMatch(t, split(c.expected), split(actual))
		})
	}
}
