package jerr

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jsightapi/jsight-schema-go-library/bytes"
)

func TestDetectNewLine(t *testing.T) {
	s :=
		`Some text with
line break
`
	nl := DetectNewLineSymbol(bytes.Bytes(s))
	assert.Equal(t, byte('\n'), nl)

	s = "Some text with \r line break"
	nl = DetectNewLineSymbol(bytes.Bytes(s))
	assert.Equal(t, byte('\r'), nl)
}

func TestLineBeginning(t *testing.T) {
	s :=
		`Some text with
line break
and one more
`
	content := bytes.Bytes(s)
	nl := DetectNewLineSymbol(content)

	lb1 := LineBeginning(content, 13, nl)
	assert.Equal(t, bytes.Index(0), lb1)
	lb2 := LineBeginning(content, 16, nl)
	assert.Equal(t, bytes.Index(15), lb2)
	lb3 := LineBeginning(content, 28, nl)
	assert.Equal(t, bytes.Index(26), lb3)
}

func TestLineEnd(t *testing.T) {
	s :=
		`Some text with
line break
and one more`
	content := bytes.Bytes(s)
	nl := DetectNewLineSymbol(content)

	le1 := LineEnd(content, 13, nl)
	assert.Equal(t, bytes.Index(14), le1)
	le2 := LineEnd(content, 16, nl)
	assert.Equal(t, bytes.Index(25), le2)
	le3 := LineEnd(content, 28, nl)
	assert.Equal(t, bytes.Index(38), le3)
}

func TestPositionInLine(t *testing.T) {
	s :=
		`Some text with
line break
and one more`
	content := bytes.Bytes(s)
	nl := DetectNewLineSymbol(content)

	p1 := PositionInLine(content, 3, nl)
	assert.Equal(t, bytes.Index(3), p1)
	p2 := PositionInLine(content, 16, nl)
	assert.Equal(t, bytes.Index(1), p2)
	p3 := PositionInLine(content, 30, nl)
	assert.Equal(t, bytes.Index(4), p3)
}

func TestSourceSubString(t *testing.T) {
	s :=
		`Some text with
line break
and one more`

	content := bytes.Bytes(s)
	nl := DetectNewLineSymbol(content)

	s1 := GetQuote(content, 3, nl)
	assert.Equal(t, s1, "Some text with")

	s2 := GetQuote(content, 16, nl)
	assert.Equal(t, s2, "line break")

	s3 := GetQuote(content, 30, nl)
	assert.Equal(t, s3, "and one more")
}

func TestLineNumber(t *testing.T) {
	s :=
		`Some text with
line break
and one more`

	content := bytes.Bytes(s)
	nl := DetectNewLineSymbol(content)

	l1 := LineNumber(content, 3, nl)
	assert.Equal(t, l1, bytes.Index(1))

	l2 := LineNumber(content, 16, nl)
	assert.Equal(t, l2, bytes.Index(2))

	l3 := LineNumber(content, 30, nl)
	assert.Equal(t, l3, bytes.Index(3))
}
