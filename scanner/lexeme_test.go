package scanner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexemeType_String(t *testing.T) {
	cc := map[LexemeType]string{
		Keyword:                "keyword",
		Parameter:              "property",
		Annotation:             "annotation",
		Schema:                 "schema",
		Json:                   "json",
		Text:                   "text",
		ContextExplicitOpening: "context-opening",
		ContextExplicitClosing: "context-closing",
		255:                    "unknown-lexeme-type",
	}

	for lt, expected := range cc {
		t.Run(expected, func(t *testing.T) {
			actual := lt.String()
			assert.Equal(t, expected, actual)
		})
	}
}
