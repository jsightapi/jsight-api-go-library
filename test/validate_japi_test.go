package test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"j/japi/jerr"
	"j/japi/kit"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestValidateJapi(t *testing.T) {
	filenames := jstFilenames(GetTestDataDir())

	t.Run("positive", func(t *testing.T) {
		positive := positiveJstFilenames(filenames)
		for _, f := range positive {
			t.Run(cutRepositoryPath(f), func(t *testing.T) {
				j := requireNewJapi(t, f)
				je := assertValidateJapi(t, j)

				// show debug info
				if je != nil {
					logJAPIError(t, je)
					t.FailNow()
				}
			})
		}
	})

	t.Run("negative", func(t *testing.T) {
		negative := negativeJstFilenames(filenames)
		for _, f := range negative {
			t.Run(cutRepositoryPath(f), func(t *testing.T) {
				j := requireNewJapi(t, f)
				je := requireValidateJapiError(t, j)
				logJAPIError(t, je)

				want, err := wantIndex(f)
				require.NoError(t, err)

				assert.Equal(t, want, int(je.Index()))

				expectedError, err := getExpectedError(f)
				require.NoError(t, err)

				if expectedError != "" {
					assert.Equal(t, expectedError, je.Error())
				}
			})
		}
	})
}

func cutRepositoryPath(p string) string {
	p, err := filepath.Abs(p)
	if err != nil {
		panic(err)
	}

	parts := strings.Split(p, string(filepath.Separator))
	var idx int
	for _, p := range parts {
		idx++
		if p == "testdata" {
			break
		}
	}

	return filepath.Join(parts[idx:]...)
}

func requireNewJapi(t *testing.T, filename string) kit.JApi {
	j, err := kit.NewJapi(filename)
	require.Nil(t, err, "NewJapi should not return an error")
	return j
}

func assertValidateJapi(t *testing.T, j kit.JApi) *jerr.JAPIError {
	je := j.ValidateJAPI()
	if je != nil {
		t.Log("ValidateJAPI should NOT return an error")
	}
	return je
}

func requireValidateJapiError(t *testing.T, j kit.JApi) *jerr.JAPIError {
	je := j.ValidateJAPI()
	require.NotNil(t, je, "ValidateJAPI should return an error")
	return je
}

func wantIndex(filename string) (int, error) {
	base := filepath.Base(filename)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	s := strings.Split(name, "_")
	if len(s) < 2 {
		return 0, errors.New("the error index is not specified in the file name")
	}
	i, err := strconv.Atoi(s[1])
	if err != nil {
		return 0, errors.New("the error index is not specified in the file name")
	}
	return i, nil
}

func getExpectedError(filename string) (string, error) {
	filename = strings.TrimSuffix(filename, ".jst") + ".error"
	c, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return "", err
	}
	return strings.TrimSpace(string(c)), nil
}

func logJAPIError(t *testing.T, e *jerr.JAPIError) {
	t.Log("Got:")
	t.Log("- Line: " + strconv.Itoa(int(e.Line())))
	t.Log("- Index: " + strconv.Itoa(int(e.Index())))
	t.Log("- Message: " + e.Error())
	t.Log("- Quote: " + e.Quote())
}

func jstFilenames(dir string) []string {
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

			if !info.IsDir() && filepath.Ext(path) == ".jst" {
				filenames = append(filenames, path)
			}
			return nil
		})

	if err != nil {
		panic(err)
	}

	return filenames
}

func positiveJstFilenames(filenames []string) []string {
	list := make([]string, 0, 10)
	for _, f := range filenames {
		if !strings.HasPrefix(filepath.Base(f), "err_") {
			list = append(list, f)
		}
	}
	return list
}

func negativeJstFilenames(filenames []string) []string {
	list := make([]string, 0, 10)
	for _, f := range filenames {
		if strings.HasPrefix(filepath.Base(f), "err_") {
			list = append(list, f)
		}
	}
	return list
}
