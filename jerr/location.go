package jerr

import (
	"j/schema/bytes"
	"j/schema/fs"
)

type Location struct {
	file  *fs.File
	quote string
	index bytes.Index
	line  bytes.Index
}

// NewLocation a bit optimized version of getting all info
func NewLocation(f *fs.File, i bytes.Index) Location {
	ef := Location{
		file:  f,
		index: i,
	}
	nl := DetectNewLineSymbol(f.Content())
	lb := LineBeginning(f.Content(), i, nl)
	ef.quote = quote(f.Content(), i, lb, nl)
	ef.line = LineNumber(f.Content(), i, nl)
	// ef.positionInLine = p - lb
	return ef
}

func (l Location) Filename() string {
	return l.file.Name()
}

func (l Location) Quote() string {
	return l.quote
}

func (l Location) Line() bytes.Index {
	return l.line
}

func (l Location) Index() bytes.Index {
	return l.index
}
