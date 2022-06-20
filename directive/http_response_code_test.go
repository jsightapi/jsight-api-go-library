package directive

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsHTTPResponseCode(t *testing.T) {
	cc := map[string]bool{
		"100": true,
		"200": true,
		"300": true,
		"400": true,
		"500": true,
		"526": true,

		"99":  false,
		"527": false,
		"999": false,

		"":      false,
		"AAA":   false,
		"0":     false,
		"00100": false,
		"3.14":  false,
	}

	for given, expected := range cc {
		t.Run(given, func(t *testing.T) {
			actual := IsHTTPResponseCode(given)
			assert.Equal(t, expected, actual)
		})
	}
}
