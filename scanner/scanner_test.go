package scanner

import (
	"fmt"
	iofs "io/fs"
	"path"
	"path/filepath"
	"testing"

	"github.com/jsightapi/jsight-schema-go-library/fs"
	"github.com/jsightapi/jsight-schema-go-library/reader"
	"github.com/stretchr/testify/assert"

	"github.com/jsightapi/jsight-api-go-library/test"
)

func TestScanner_successes(t *testing.T) {
	keywords := []string{ // only which do not require values
		"GET",
		"POST",
		"PUT",
		"PATCH",
		"DELETE",
		"Query",
		"Body",
		"200",
		"422",
		"500",
	}

	for _, k := range keywords {
		t.Run(k, func(t *testing.T) {
			assert.NotPanics(t, func() {
				s := newTestScanner(k)
				scan(s)
			})
		})
	}

	withValues := []string{
		"URL /",
		"URL / ",
		"GET / ",
		"URL /users/{id}/posts/{id}",
		"# comment \n URL /users // expl",
		"# comment \n URL /users # // expl",
	}
	for _, k := range withValues {
		t.Run(k, func(t *testing.T) {
			assert.NotPanics(t, func() {
				s := newTestScanner(k)
				scan(s)
			})
		})
	}
}

func TestScanner_Next(t *testing.T) {
	str := `GET /users`
	bytes := []byte(str)

	file := fs.NewFile("dummy.jst", bytes)
	s := NewJApiScanner(file)
	scan(s)
}

func TestScanner_files_success(t *testing.T) {
	tt := make([]string, 0, 20)

	err := filepath.Walk(path.Join(test.GetTestDataDir(), "jsight_0.3", "others", "scanner"), func(path string, info iofs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			tt = append(tt, path)
		}
		return nil
	})

	if err != nil {
		t.Error(err)
		return
	}

	for _, fp := range tt {
		t.Run(fp, func(t *testing.T) {
			assert.NotPanics(t, func() {
				t.Log(fp)
				file := reader.Read(fp)
				s := NewJApiScanner(file)
				scan(s)
			})
		})
	}
}

func TestScanner_full(t *testing.T) {
	filename := filepath.Join(test.GetTestDataDir(), "jsight_0.3", "others", "full.jst")
	file := reader.Read(filename)
	s := NewJApiScanner(file)
	assert.NotPanics(t, func() {
		scan(s)
	})
}

func scan(s *Scanner) {
	for {
		lexeme, je := s.Next()

		if je != nil {
			fmt.Println(je.Msg)
		}

		if lexeme == nil {
			break
		}
	}
}

func newTestScanner(s string) *Scanner {
	return newTestScannerB([]byte(s))
}

func newTestScannerB(bytes []byte) *Scanner {
	file := fs.NewFile("dummy.jst", bytes)
	return NewJApiScanner(file)
}
