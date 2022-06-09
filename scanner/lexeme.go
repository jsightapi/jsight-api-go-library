package scanner

import (
	"fmt"

	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/fs"
)

type LexemeType uint8

const (
	Keyword                LexemeType = iota // Name of a Directive, i.e. URL, GET, Path, 200, etc
	Parameter                                // Parameter for directive (regexp/jsight fot TYPE)
	Annotation                               // user's annotation to directive in a free-text form
	Schema                                   // jSchema inside directive's body (Body, TYPE, 200, etc)
	Json                                     // json inside directive's body (CONFIG)
	Array                                    // Array of string values inside ENUM body
	Text                                     // Text inside directive Description body
	ContextExplicitOpening                   // Explicitly opens context, so that it can ba later explicitly closed
	ContextExplicitClosing                   // Explicitly opens context, so that it can ba later explicitly closed
)

func (t LexemeType) String() string {
	if s, ok := lexemeTypeStringMap[t]; ok {
		return s
	}
	return "unknown-lexeme-type"
}

var lexemeTypeStringMap = map[LexemeType]string{
	Keyword:                "keyword",
	Parameter:              "property",
	Annotation:             "annotation",
	Schema:                 "schema",
	Json:                   "json",
	Array:                  "array",
	Text:                   "text",
	ContextExplicitOpening: "context-opening",
	ContextExplicitClosing: "context-closing",
}

type Lexeme struct {
	file  *fs.File    // File containing the contents of the json and the file name
	type_ LexemeType  // Type of found lexeme
	begin bytes.Index // bytes.Index of the start character of the found lexeme in the byte slice
	end   bytes.Index // bytes.Index of the end character of the found lexeme in the byte slice
}

func NewLexeme(type_ LexemeType, begin bytes.Index, end bytes.Index, file *fs.File) *Lexeme {
	return &Lexeme{
		file:  file,
		type_: type_,
		begin: begin,
		end:   end,
	}
}

func (lex Lexeme) Value() bytes.Bytes {
	return lex.file.Content().Slice(lex.begin, lex.end)
}

func (lex Lexeme) File() *fs.File {
	return lex.file
}

func (lex Lexeme) Type() LexemeType {
	return lex.type_
}

func (lex Lexeme) Begin() bytes.Index {
	return lex.begin
}

func (lex Lexeme) End() bytes.Index {
	return lex.end
}

func (lex Lexeme) String() string {
	return fmt.Sprintf("%s [%d:%d]", lex.type_.String(), lex.begin, lex.end)
}
