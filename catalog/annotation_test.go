package catalog

import (
	"testing"
)

func TestAnnotation(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"  leading whitespaces",
			"leading whitespaces",
		},
		{
			"trailing whitespaces  ",
			"trailing whitespaces",
		},
		{
			"line\nfeed",
			"line feed",
		},
		{
			"carriage\rreturn",
			"carriage return",
		},
		{
			"horizontal\ttab",
			"horizontal tab",
		},
		{
			" \n\r\tall  \n\n \r\r \t\t  together  \n\r\t",
			"all together",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Annotation(tt.name); got != tt.want {
				t.Errorf("Annotation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkAnnotation(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		Annotation(" \n\n\tall  \n\n \n\n \t\t  together  \n\n\t")
	}
}
