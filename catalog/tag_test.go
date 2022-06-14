package catalog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newEmptyTag(t *testing.T) {
	tag := newEmptyTag(HttpInteractionId{
		protocol: http,
		path:     "/foo/bar",
	})

	assert.Equal(t, make(map[Protocol]TagInteractionGroup), tag.InteractionGroups)
	assert.Equal(t, &Tags{}, tag.Children)
	assert.Equal(t, "/foo", tag.Title)
	assert.Equal(t, TagName("@foo"), tag.Name)
}

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
