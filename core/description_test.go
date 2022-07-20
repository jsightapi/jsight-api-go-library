package core

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_longestWhitespacePrefix(t *testing.T) {
	tests := []struct {
		arg  [][]byte
		want []byte
	}{
		{
			[][]byte{},
			[]byte(""),
		},
		{
			[][]byte{
				[]byte(""),
			},
			[]byte(""),
		},
		{
			[][]byte{
				[]byte(""),
				[]byte(""),
			},
			[]byte(""),
		},
		{
			[][]byte{
				[]byte("aaa"),
				[]byte("bbb"),
			},
			[]byte(""),
		},
		{
			[][]byte{
				[]byte("  aaa"),
				[]byte("bbb"),
			},
			[]byte(""),
		},
		{
			[][]byte{
				[]byte("aaa"),
				[]byte("  bbb"),
			},
			[]byte(""),
		},
		{
			[][]byte{
				[]byte("\taaa"),
				[]byte("bbb"),
			},
			[]byte(""),
		},
		{
			[][]byte{
				[]byte("aaa"),
				[]byte("\tbbb"),
			},
			[]byte(""),
		},
		{
			[][]byte{
				[]byte("\taaa"),
				[]byte("\tbbb"),
			},
			[]byte("\t"),
		},
		{
			[][]byte{
				[]byte("\taaa"),
				[]byte("\t\t\tbbb"),
			},
			[]byte("\t"),
		},
		{
			[][]byte{
				[]byte("  aaa"),
				[]byte("  bbb"),
			},
			[]byte("  "),
		},
		{
			[][]byte{
				[]byte("  aaa"),
				[]byte(" bbb"),
			},
			[]byte(" "),
		},
		{
			[][]byte{
				[]byte("\t aaa"),
				[]byte("\tbbb"),
			},
			[]byte("\t"),
		},
		{
			[][]byte{
				[]byte("\t  aaa"),
				[]byte("\t bbb"),
			},
			[]byte("\t "),
		},
		{
			[][]byte{
				[]byte("\t\taaa"),
				[]byte("\t\tbbb"),
				[]byte("\t bbb"),
			},
			[]byte("\t"),
		},
		{
			[][]byte{
				[]byte("\t\t \t"),
				[]byte("\t\t aaa"),
				[]byte("\t\t bbb"),
				[]byte("\t\t     "),
			},
			[]byte("\t\t "),
		},
		{
			[][]byte{
				[]byte(""),
				[]byte("  "),
				[]byte("\t\t"),
			},
			[]byte(""),
		},
		{
			[][]byte{
				[]byte("\t \t aaa"),
				[]byte("\t \t bbb"),
				[]byte("\t \tccc"),
			},
			[]byte("\t \t"),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := longestWhitespacePrefix(tt.arg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_descriptionRemoveParentheses(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		tests := []struct {
			arg  string
			want string
		}{
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
				"bbb)",
				"bbb)",
			},
			{
				"(\naaa",
				"(\naaa",
			},
			{
				"bbb\n)",
				"bbb\n)",
			},
			{
				"(\n)",
				"",
			},
			{
				"(\n\n)",
				"",
			},
			{
				"(\n\n\n)",
				"",
			},
			{
				"\n\n(\n\n)\n\n",
				"",
			},
			{
				"(\naaa\n)",
				"aaa",
			},
			{
				"(\naaa\nbbb\n)",
				"aaa\nbbb",
			},
			{
				"( \naaa\n )",
				"aaa",
			},
			{
				"(\n aaa \n)",
				" aaa ",
			},
			{
				"(\n aaa\n bbb\n)",
				" aaa\n bbb",
			},
			{
				" \t (\naaa\n) \t ",
				"aaa",
			},
			{
				" \n (\n\naaa\nbbb\n\n) \n ",
				"aaa\nbbb",
			},
		}
		for _, tt := range tests {
			t.Run(tt.arg, func(t *testing.T) {
				b := []byte(tt.arg)
				got, err := descriptionRemoveParentheses(b)
				assert.NoError(t, err, "%q", tt)
				assert.Equalf(t, tt.want, string(got), "Test data: %q", tt.arg)
			})
		}
	})

	t.Run("negative", func(t *testing.T) {
		tests := []string{
			"()",
			" () ",
			"( )",
			"(  )",
			"(\t)",
			"(\t\t)",
			"(aaa)",
			"( aaa )",
			"(\naaa)",
			"(aaa\n)",
			"(aaa\nbbb)",
			" (aaa) ",
			"\t(aaa)\t",
		}
		for _, tt := range tests {
			t.Run(tt, func(t *testing.T) {
				b := []byte(tt)
				_, err := descriptionRemoveParentheses(b)
				assert.Error(t, err, "Test data: %q", tt)
			})
		}
	})
}

func Test_description(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		tests := []struct {
			arg  []byte
			want []byte
		}{
			{
				[]byte("\n(\naaa\nbbb\n)\n"),
				[]byte("aaa\nbbb"),
			},
			{
				[]byte("\r(\raaa\rbbb\r)\r"),
				[]byte("aaa\nbbb"),
			},
			{
				[]byte("\r\n(\r\naaa\r\nbbb\r\n)\r\n"),
				[]byte("aaa\nbbb"),
			},
			{
				[]byte("\n\r(\n\raaa\n\rbbb\n\r)\n\r"), // incorrect sequence of newline characters
				[]byte("aaa\n\nbbb"),
			},
			{
				[]byte("\n(\n  aaa\n  bbb\n)\n"),
				[]byte("aaa\nbbb"),
			},
			{
				[]byte("\n(\n  aaa  \n  bbb  \n)\n"),
				[]byte("aaa  \nbbb"),
			},
			{
				[]byte("\n(\naaa\n  bbb  \nccc\n)\n"),
				[]byte("aaa\n  bbb  \nccc"),
			},
			{
				[]byte("\n(\n\taaa\n bbb\n\tccc\n)\n"),
				[]byte("\taaa\n bbb\n\tccc"),
			},
			{
				[]byte("\n(\n  aaa\n  bbb\n ccc\n)\n"),
				[]byte(" aaa\n bbb\nccc"),
			},
		}
		for _, tt := range tests {
			t.Run(string(tt.arg), func(t *testing.T) {
				got, err := description(tt.arg)
				assert.NoError(t, err, got, "arg: %q", tt.arg)
				assert.Equalf(t, tt.want, got, "arg: %q\ngot: %q\nwant: %q", tt.arg, got, tt.want)
			})
		}
	})
}
