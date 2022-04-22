package test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"j/japi/core"
	"j/schema/reader"
	"path/filepath"
	"testing"
)

func TestGetAllTypesSchemas(t *testing.T) {
	filename := filepath.Join(GetTestDataDir(), "jsight_0.3", "others", "full.jst")
	f := reader.Read(filename)
	japi := core.NewJApiCore(f)
	err := japi.ValidateJAPI()
	require.Nil(t, err)

	c := japi.Catalog()
	types := c.UserTypes

	assert.Equal(t, 5, types.Len())

	assert.True(t, types.Has("@error"))
	assert.True(t, types.Has("@user"))
	assert.True(t, types.Has("@profile"))
	assert.True(t, types.Has("@task"))
	assert.True(t, types.Has("@attachment"))
	// assert.True(t, types.Has("@userType")) // enum
}
