package scanner

import (
	"github.com/jsightapi/jsight-schema-go-library/bytes"
)

type LexemeEvent struct {
	type_    LexemeEventType
	position bytes.Index
}

type LexemeEventType uint8

const (
	KeywordBegin LexemeEventType = iota
	KeywordEnd
	ParameterBegin
	ParameterEnd
	AnnotationBegin
	AnnotationEnd
	SchemaBegin
	SchemaEnd
	JsonArrayBegin
	JsonArrayEnd
	TextBegin
	TextEnd
	ContextOpen
	ContextClose
)

func (e LexemeEventType) IsBeginning() bool {
	switch e {
	case KeywordBegin,
		ParameterBegin,
		AnnotationBegin,
		SchemaBegin,
		JsonArrayBegin,
		TextBegin:
		return true
	default:
		return false
	}
}

func (e LexemeEventType) IsEnding() bool {
	switch e {
	case KeywordEnd,
		ParameterEnd,
		AnnotationEnd,
		SchemaEnd,
		JsonArrayEnd,
		TextEnd:
		return true
	default:
		return false
	}
}

func (e LexemeEventType) IsSingle() bool {
	switch e {
	case ContextOpen, ContextClose:
		return true
	default:
		return false
	}
}

func (e LexemeEventType) String() string {
	if s, ok := lexemeEventTypeStringMap[e]; ok {
		return s
	}
	return "Unknown-lexeme-event-type"
}

var lexemeEventTypeStringMap = map[LexemeEventType]string{
	KeywordBegin:    "keyword-begin",
	KeywordEnd:      "keyword-end",
	ParameterBegin:  "property-begin",
	ParameterEnd:    "property-end",
	AnnotationBegin: "annotation-begin",
	AnnotationEnd:   "annotation-end",
	SchemaBegin:     "schema-begin",
	SchemaEnd:       "schema-end",
	JsonArrayBegin:  "array-begin",
	JsonArrayEnd:    "array-end",
	TextBegin:       "text-begin",
	TextEnd:         "text-end",
	ContextOpen:     "context-open",
	ContextClose:    "context-close",
}

func (e LexemeEventType) ToLexemeType() LexemeType {
	switch e {
	case KeywordBegin, KeywordEnd:
		return Keyword
	case AnnotationBegin, AnnotationEnd:
		return Annotation
	case SchemaBegin, SchemaEnd:
		return Schema
	case TextBegin, TextEnd:
		return Text
	case JsonArrayBegin, JsonArrayEnd:
		return Array
	case ContextOpen:
		return ContextExplicitOpening
	case ContextClose:
		return ContextExplicitClosing
	case ParameterBegin, ParameterEnd:
		return Parameter
	default:
		panic("Unknown lexeme event type")
	}
}
