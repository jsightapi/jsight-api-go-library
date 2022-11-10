package jerr

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jsightapi/jsight-schema-go-library/bytes"
)

func TestGetQuote(t *testing.T) {
	s :=
		`Some text with
line break
and one more`

	content := bytes.NewBytes(s)

	s1 := GetQuote(content, 3)
	assert.Equal(t, s1, "Some text with")

	s2 := GetQuote(content, 16)
	assert.Equal(t, s2, "line break")

	s3 := GetQuote(content, 30)
	assert.Equal(t, s3, "and one more")
}
