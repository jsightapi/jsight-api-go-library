package directive

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsHTTPResponseCode(t *testing.T) {
	assert.True(t, IsHTTPResponseCode("100"))
	assert.True(t, IsHTTPResponseCode("200"))
	assert.True(t, IsHTTPResponseCode("300"))
	assert.True(t, IsHTTPResponseCode("400"))
	assert.True(t, IsHTTPResponseCode("500"))
	assert.True(t, IsHTTPResponseCode("526"))

	assert.False(t, IsHTTPResponseCode("527"))
	assert.False(t, IsHTTPResponseCode("999"))

	assert.False(t, IsHTTPResponseCode(""))
	assert.False(t, IsHTTPResponseCode("AAA"))
	assert.False(t, IsHTTPResponseCode("0"))
	assert.False(t, IsHTTPResponseCode("00100"))
}
