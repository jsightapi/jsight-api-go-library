package test

import (
	"path/filepath"

	"github.com/jsightapi/jsight-schema-go-library/test"
)

func GetTestDataDir() string {
	return filepath.Join(test.GetProjectRoot(), "testdata")
}
