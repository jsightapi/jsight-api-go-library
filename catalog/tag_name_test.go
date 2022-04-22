package catalog

import "testing"

func Test_tagName(t *testing.T) {
	tests := []struct {
		title string
		want  TagName
	}{
		{
			"/",
			"@_",
		},
		{
			"/aaa",
			"@aaa",
		},
		{
			"/aaa_bbb",
			"@aaa__bbb",
		},
		{
			"/{id}",
			"@_7Bid_7D",
		},
		{
			"/{pet_id}",
			"@_7Bpet__id_7D",
		},
		{
			"/%",
			"@_25",
		},
		{
			"/_25",
			"@__25",
		},
		{
			"/_25",
			"@__25",
		},
		{
			"/_%",
			"@___25",
		},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			if got := tagName(tt.title); got != tt.want {
				t.Errorf("newTagName() = %v, want %v", got, tt.want)
			}
		})
	}
}
