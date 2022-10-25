package core

import (
	"github.com/jsightapi/jsight-api-go-library/directive"
	jschemaLib "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/fs"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jsightapi/jsight-api-go-library/catalog"
)

func TestJApiCore_buildPathVariables(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		cc := map[string]struct {
			given    *JApiCore
			expected *JApiCore
		}{
			"empty": {
				&JApiCore{
					catalog: &catalog.Catalog{
						Interactions: &catalog.Interactions{},
					},
				},
				&JApiCore{
					catalog: &catalog.Catalog{
						Interactions: &catalog.Interactions{},
					},
				},
			},
		}

		for n, c := range cc {
			t.Run(n, func(t *testing.T) {
				err := c.given.buildPathVariables()
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
						schema: jschema.New("", `{"foo": 1, "bar": 2}`),
						parameters: []PathParameter{
							{
								path:      "/path/{foo}/{bar}",
								parameter: "foo",
							},
						},
						pathDirective: *directive.New(
							directive.Get,
							directive.NewCoords(
								fs.NewFile("foo", "/path/{foo}/{bar}"),
								0,
								5,
							),
						),
					},
				},

				catalog: &catalog.Catalog{
					Interactions: &catalog.Interactions{},
				},

				allProjectProperties: make(map[catalog.Path]catalog.Prop, 20),
			},
		}

		for expected, core := range cc {
			t.Run(expected, func(t *testing.T) {
				for i := range core.rawPathVariables {
					_ = core.rawPathVariables[i].schema.Build()
				}
				err := core.buildPathVariables()
				assert.EqualError(t, err, expected)
			})
		}

		t.Run("nil", func(t *testing.T) {
			assert.Panics(t, func() {
				var core *JApiCore

				_ = core.buildPathVariables()
			})
		})
	})
}

func TestCore_getPropertiesNames(t *testing.T) {
	cc := map[string]struct {
		given    map[string]jschemaLib.ASTNode
		expected string
	}{
		"nil": {},

		"empty": {
			given: map[string]jschemaLib.ASTNode{},
		},

		"with data": {
			given: map[string]jschemaLib.ASTNode{
				"foo":  {},
				"bar":  {},
				"fizz": {},
				"buzz": {},
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
