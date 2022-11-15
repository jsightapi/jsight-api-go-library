package core

import (
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/ischema"

	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCore_getPropertiesNames(t *testing.T) {
	cc := map[string]struct {
		given    map[string]ischema.Node
		expected string
	}{
		"nil": {},

		"empty": {
			given: map[string]ischema.Node{},
		},

		"with data": {
			given: map[string]ischema.Node{
				"foo":  &ischema.LiteralNode{},
				"bar":  &ischema.LiteralNode{},
				"fizz": &ischema.LiteralNode{},
				"buzz": &ischema.LiteralNode{},
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
