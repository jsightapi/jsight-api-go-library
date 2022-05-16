package core

import (
	"testing"

	"github.com/jsightapi/jsight-schema-go-library/fs"
	"github.com/stretchr/testify/assert"

	"github.com/jsightapi/jsight-api-go-library/directive"
)

func TestWithBannedDirectives(t *testing.T) {
	err := NewJApiCore(fs.NewFile("", []byte("JSIGHT 0.3")), WithBannedDirectives(directive.Jsight)).ValidateJAPI()
	assert.EqualError(t, err, "directive not allowed (JSIGHT)")
}
