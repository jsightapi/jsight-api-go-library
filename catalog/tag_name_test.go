package catalog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_tagName(t *testing.T) {
	tests := map[string]TagName{
		"/":         "@_",
		"/aaa":      "@aaa",
		"/aaa_bbb":  "@aaa__bbb",
		"/{id}":     "@_7Bid_7D",
		"/{pet_id}": "@_7Bpet__id_7D",
		"/%":        "@_25",
		"/_25":      "@__25",
		"/_%":       "@___25",
	}

	for given, expected := range tests {
		t.Run(given, func(t *testing.T) {
			actual := tagName(given)
			assert.Equal(t, expected, actual)
		})
	}
}
