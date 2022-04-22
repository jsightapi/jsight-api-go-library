package test

import (
	"j/schema/test"
	"path/filepath"
)

func GetTestDataDir() string {
	return filepath.Join(test.GetProjectRoot(), "testdata")
}
