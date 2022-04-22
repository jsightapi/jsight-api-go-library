package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_camelCaseToUnderscore(t *testing.T) {
	cc := map[string]string{
		"":       "",
		"foo":    "foo",
		"Foo":    "foo",
		"fooBar": "foo_bar",
		"FooBar": "foo_bar",
	}

	for given, expected := range cc {
		t.Run(given, func(t *testing.T) {
			actual := camelCaseToUnderscore(given)
			assert.Equal(t, expected, actual)
		})
	}
}

func Benchmark_camelCaseToUnderscore(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		camelCaseToUnderscore("VeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryVeryLongString")
	}
}
