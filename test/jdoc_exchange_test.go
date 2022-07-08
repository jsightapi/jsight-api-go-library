package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jsightapi/jsight-api-go-library/kit"
)

func TestJDocExchange(t *testing.T) {
	jsonFilesPaths := jsonFilePaths(GetTestDataDir())

	for _, jsonPath := range jsonFilesPaths {
		t.Run(cutRepositoryPath(jsonPath), func(t *testing.T) {
			json, err := os.ReadFile(jsonPath)
			require.NoError(t, err)

			japiPath, err := japiFilePath(jsonPath)
			require.NoError(t, err)

			j, err := kit.NewJapi(japiPath)
			require.NoError(t, err)

			je := j.ValidateJAPI()
			if je != nil {
				logJAPIError(t, je)
				t.FailNow()
			}

			actual, err := j.ToJsonIndent()
			require.NoError(t, err)

			expected := string(json)

			ok := assert.JSONEq(t, expected, string(actual))

			if !ok {
				t.Log("Actual JSON:")
				t.Log(string(actual))
			}
		})
	}
}

func jsonFilePaths(dir string) []string {
	filenames := make([]string, 0, 30)

	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			base := filepath.Base(path)

			if info.IsDir() && (base == ".unused" || base == "scanner") {
				return filepath.SkipDir
			}

			if !info.IsDir() && filepath.Ext(path) == ".json" {
				filenames = append(filenames, path)
			}
			return nil
		})

	if err != nil {
		panic(err)
	}

	return filenames
}

func japiFilePath(japiPath string) (string, error) {
	ext := filepath.Ext(japiPath)
	jsonPath := japiPath[:len(japiPath)-len(ext)] + ".jst"
	return jsonPath, nil
}
