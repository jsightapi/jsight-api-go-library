package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_longestWhitespacePrefix(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"",
			"",
		},
		{
			"aaa\nbbb",
			"",
		},
		{
			"  aaa\nbbb",
			"",
		},
		{
			"aaa\n  bbb",
			"",
		},
		{
			"\taaa\nbbb",
			"",
		},
		{
			"aaa\n\tbbb",
			"",
		},
		{
			"\taaa\n\tbbb",
			"\t",
		},
		{
			"   aaa\n   bbb",
			"   ",
		},
		{
			"\t\taaa\n\tbbb",
			"\t",
		},
		{
			"\t \taaa\n\t bbb",
			"\t ",
		},
		{
			"\t\taaa\n\t\tbbb\n\t  ccc",
			"\t",
		},
		{
			"\t\t  \t\n\t\t aaa\n\t\t bbb\n\t\t       ",
			"\t\t ",
		},
		{
			"\n\t\t  \t\n\t\t aaa\n\t\t bbb\n\t\t       ",
			"",
		},
		{
			"\t \t \t aaa\n\t \t bbb\n\t ccc",
			"\t ",
		},
		{
			"\t \t \t aaa\n\t \t bbb\n\t ccc\n\t \n\t \n\t \t \t \t \t \t \t ",
			"\t ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := splitLines([]byte(tt.name))
			assert.Equalf(t, []byte(tt.want), longestWhitespacePrefix(lines), "longestWhitespacePrefix with %v", tt.name)
		})
	}
}

func Test_descriptionTrimBrackets(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"(aaa)",
			"aaa",
		},
		{
			" \t\t\t    (aaa)\t\n  \r\t  ",
			"aaa",
		},
		{
			"()",
			"",
		},
		{
			"   ( )   ",
			" ",
		},
		{
			"   (\t)   ",
			"\t",
		},
		{
			"\n\n(\r\r)\n\n",
			"\r\r",
		},
		{
			"(",
			"(",
		},
		{
			")",
			")",
		},
		{
			"(aaa",
			"(aaa",
		},
		{
			"aaa)",
			"aaa)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, []byte(tt.want), descriptionTrimBrackets([]byte(tt.name)), "descriptionTrimBrackets(%v)", tt.name)
		})
	}
}

func Test_splitLines(t *testing.T) {
	tests := []struct {
		name string
		arg  []byte
		want [][]byte
	}{
		{
			"LF",
			[]byte("\naaa\nbbb\n\nccc\n"),
			[][]byte{[]byte(""), []byte("aaa"), []byte("bbb"), []byte(""), []byte("ccc"), []byte("")},
		},
		{
			"CRLF",
			[]byte("\r\naaa\r\nbbb\r\n\r\nccc\r\n"),
			[][]byte{[]byte(""), []byte("aaa"), []byte("bbb"), []byte(""), []byte("ccc"), []byte("")},
		},
		{
			"CR",
			[]byte("\raaa\rbbb\r\rccc\r"),
			[][]byte{[]byte(""), []byte("aaa"), []byte("bbb"), []byte(""), []byte("ccc"), []byte("")},
		},
		{
			"MIXED",
			[]byte("\raaa\nbbb\r\n\r\nccc\n"),
			[][]byte{[]byte(""), []byte("aaa"), []byte("bbb"), []byte(""), []byte("ccc"), []byte("")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := splitLines(tt.arg)
			assert.Equalf(t, tt.want, x, "splitLines(%v)", tt.arg)
		})
	}
}
