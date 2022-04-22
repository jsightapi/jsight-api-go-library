package test

import (
	"j/japi/kit"
	"path/filepath"
	"testing"
)

func BenchmarkJAPI(b *testing.B) {
	filename := filepath.Join(GetTestDataDir(), "jsight_0.3", "others", "full.jst")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		j, err := kit.NewJapi(filename)
		if err != nil {
			b.Error(err)
		}

		je := j.ValidateJAPI()
		if je != nil {
			b.Error(je)
		}

		_, err = j.ToJson()
		if err != nil {
			b.Error(je)
		}
	}
}
