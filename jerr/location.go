package jerr

import (
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/fs"
)

type Location struct {
	file   *fs.File
	quote  string
	index  bytes.Index
	line   bytes.Index
	column bytes.Index
}

// NewLocation a bit optimized version of getting all info
func NewLocation(f *fs.File, i bytes.Index) Location {
	loc := Location{
		file:  f,
		index: i,
		quote: quote(f.Content(), i),
	}
	loc.line, loc.column = f.Content().LineAndColumn(i)
	return loc
}

func (l Location) Quote() string {
	return l.quote
}

func (l Location) Line() bytes.Index {
	return l.line
}

func (l Location) Column() bytes.Index {
	return l.column
}

func (l Location) Index() bytes.Index {
	return l.index
}
