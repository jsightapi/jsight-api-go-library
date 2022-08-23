package catalog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newPathTag(t *testing.T) {
	tag := newPathTag(HTTPInteractionID{
		protocol: HTTP,
		path:     "/foo/bar",
	})

	assert.Equal(t, make(map[Protocol]TagInteractionGroup), tag.InteractionGroups)
	assert.Equal(t, &Tags{}, tag.Children)
	assert.Equal(t, "/foo", tag.Title)
	assert.Equal(t, TagName("@foo"), tag.Name)
}

func Test_pathTagTitle(t *testing.T) {
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
			if got := pathTagTitle(tt.url); got != tt.want {
				t.Errorf("pathTagTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}
