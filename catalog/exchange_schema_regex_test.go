package catalog

import (
	schema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegexMarshaller_Marshal(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		s, err := PrepareRegexSchema("foo", bytes.NewBytes("/bar-\\d/"))
		require.NoError(t, err)

		n, err := s.GetAST()
		require.NoError(t, err)
		assert.Equal(t, schema.ASTNode{
			TokenType:  schema.TokenTypeString,
			SchemaType: schema.TokenTypeString,
			Value:      "/bar-\\d/",
		}, n)

		example, err := s.Example()
		require.NoError(t, err)
		assert.Equal(t, []byte("bar-4"), example)
	})

	t.Run("negative", func(t *testing.T) {
		s, err := PrepareRegexSchema("", bytes.NewBytes("invalid"))
		require.NoError(t, err)

		_, err = s.GetAST()
		assert.EqualError(t, err, `ERROR (code 1500): Regex should starts with '/' character, but found 'i'
	in line 1 on file 
	> invalid
	--^`)
	})
}
