package core

import (
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/schema"

	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCore_getPropertiesNames(t *testing.T) {
	cc := map[string]struct {
		given    map[string]schema.Node
		expected string
	}{
		"nil": {},

		"empty": {
			given: map[string]schema.Node{},
		},

		"with data": {
			given: map[string]schema.Node{
				"foo":  &schema.LiteralNode{},
				"bar":  &schema.LiteralNode{},
				"fizz": &schema.LiteralNode{},
				"buzz": &schema.LiteralNode{},
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
