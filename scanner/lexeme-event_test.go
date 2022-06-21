package scanner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexemeEventType_IsBeginning(t *testing.T) {
	cc := map[LexemeEventType]bool{
		KeywordBegin:    true,
		KeywordEnd:      false,
		ParameterBegin:  true,
		ParameterEnd:    false,
		AnnotationBegin: true,
		AnnotationEnd:   false,
		SchemaBegin:     true,
		SchemaEnd:       false,
		TextBegin:       true,
		TextEnd:         false,
		ContextOpen:     false,
		ContextClose:    false,
		EnumBegin:       true,
		EnumEnd:         false,
	}

	for et, expected := range cc {
		t.Run(et.String(), func(t *testing.T) {
			actual := et.IsBeginning()
			assert.Equal(t, expected, actual)
		})
	}
}

func TestLexemeEventType_IsEnding(t *testing.T) {
	cc := map[LexemeEventType]bool{
		KeywordBegin:    false,
		KeywordEnd:      true,
		ParameterBegin:  false,
		ParameterEnd:    true,
		AnnotationBegin: false,
		AnnotationEnd:   true,
		SchemaBegin:     false,
		SchemaEnd:       true,
		TextBegin:       false,
		TextEnd:         true,
		ContextOpen:     false,
		ContextClose:    false,
		EnumBegin:       false,
		EnumEnd:         true,
	}

	for et, expected := range cc {
		t.Run(et.String(), func(t *testing.T) {
			actual := et.IsEnding()
			assert.Equal(t, expected, actual)
		})
	}
}

func TestLexemeEventType_IsSingle(t *testing.T) {
	cc := map[LexemeEventType]bool{
		KeywordBegin:    false,
		KeywordEnd:      false,
		ParameterBegin:  false,
		ParameterEnd:    false,
		AnnotationBegin: false,
		AnnotationEnd:   false,
		SchemaBegin:     false,
		SchemaEnd:       false,
		TextBegin:       false,
		TextEnd:         false,
		ContextOpen:     true,
		ContextClose:    true,
		EnumBegin:       false,
		EnumEnd:         false,
	}

	for et, expected := range cc {
		t.Run(et.String(), func(t *testing.T) {
			actual := et.IsSingle()
			assert.Equal(t, expected, actual)
		})
	}
}

func TestLexemeEventType_String(t *testing.T) {
	cc := map[LexemeEventType]string{
		KeywordBegin:    "keyword-begin",
		KeywordEnd:      "keyword-end",
		ParameterBegin:  "property-begin",
		ParameterEnd:    "property-end",
		AnnotationBegin: "annotation-begin",
		AnnotationEnd:   "annotation-end",
		SchemaBegin:     "schema-begin",
		SchemaEnd:       "schema-end",
		TextBegin:       "text-begin",
		TextEnd:         "text-end",
		ContextOpen:     "context-open",
		ContextClose:    "context-close",
		EnumBegin:       "enum-begin",
		EnumEnd:         "enum-end",
		255:             "Unknown-lexeme-event-type",
	}

	for et, expected := range cc {
		t.Run(expected, func(t *testing.T) {
			actual := et.String()
			assert.Equal(t, expected, actual)
		})
	}
}

func TestLexemeEventType_ToLexemeType(t *testing.T) {
	cc := map[LexemeEventType]LexemeType{
		KeywordBegin:    Keyword,
		KeywordEnd:      Keyword,
		ParameterBegin:  Parameter,
		ParameterEnd:    Parameter,
		AnnotationBegin: Annotation,
		AnnotationEnd:   Annotation,
		SchemaBegin:     Schema,
		SchemaEnd:       Schema,
		TextBegin:       Text,
		TextEnd:         Text,
		ContextOpen:     ContextExplicitOpening,
		ContextClose:    ContextExplicitClosing,
		EnumBegin:       Enum,
		EnumEnd:         Enum,
	}

	for et, expected := range cc {
		t.Run(et.String(), func(t *testing.T) {
			actual := et.ToLexemeType()
			assert.Equal(t, expected, actual)
		})
	}
}
