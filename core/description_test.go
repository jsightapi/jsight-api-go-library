package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_longestWhitespacePrefix(t *testing.T) {
	cc := map[string]string{
		"":                          "",
		"aaa\nbbb":                  "",
		"  aaa\nbbb":                "",
		"aaa\n  bbb":                "",
		"\taaa\nbbb":                "",
		"aaa\n\tbbb":                "",
		"\taaa\n\tbbb":              "\t",
		"   aaa\n   bbb":            "   ",
		"\t\taaa\n\tbbb":            "\t",
		"\t \taaa\n\t bbb":          "\t ",
		"\t\taaa\n\t\tbbb\n\t  ccc": "\t",
		"\t\t  \t\n\t\t aaa\n\t\t bbb\n\t\t       ":                        "\t\t ",
		"\n\t\t  \t\n\t\t aaa\n\t\t bbb\n\t\t       ":                      "",
		"\t \t \t aaa\n\t \t bbb\n\t ccc":                                  "\t ",
		"\t \t \t aaa\n\t \t bbb\n\t ccc\n\t \n\t \n\t \t \t \t \t \t \t ": "\t ",
	}

	for given, expected := range cc {
		t.Run(given, func(t *testing.T) {
			actual := longestWhitespacePrefix(splitLines([]byte(given)))
			assert.Equal(t, expected, string(actual))
		})
	}
}

func Test_descriptionTrimBrackets(t *testing.T) {
	cc := map[string]string{
		"(aaa)":                        "aaa",
		" \t\t\t    (aaa)\t\n  \r\t  ": "aaa",
		"()":                           "",
		"   ( )   ":                    " ",
		"   (\t)   ":                   "\t",
		"\n\n(\r\r)\n\n":               "\r\r",
		"(":                            "(",
		")":                            ")",
		"(aaa":                         "(aaa",
		"aaa)":                         "aaa)",
	}

	for given, expected := range cc {
		t.Run(given, func(t *testing.T) {
			actual := descriptionTrimBrackets([]byte(given))
			assert.Equal(t, expected, string(actual))
		})
	}
}
