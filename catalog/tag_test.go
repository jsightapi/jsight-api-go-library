package catalog

import (
	"testing"
)

func Test_tagTitle(t *testing.T) {
	tests := []struct {
		url  string
		want string
	}{
		{
			"",
			"/",
		},
		{
			"/",
			"/",
		},
		{
			"//",
			"/",
		},
		{
			"///",
			"/",
		},
		{
			"aaa",
			"/aaa",
		},
		{
			"/aaa",
			"/aaa",
		},
		{
			"//aaa",
			"/aaa",
		},
		{
			"/aaa/bbb",
			"/aaa",
		},
		{
			"aaa/bbb",
			"/aaa",
		},
		{
			`./../..`,
			"/..",
		},
		{
			`./ddd`,
			"/ddd",
		},
		{
			`../../`,
			"/..",
		},
		{
			`/Hello, 世界/some`,
			"/Hello, 世界",
		},
		{
			`\\\`,
			`/\\\`,
		},
		{
			`/{id}/aaa/bbb`,
			`/{id}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			if got := tagTitle(tt.url); got != tt.want {
				t.Errorf("tagTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}
