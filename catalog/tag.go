package catalog

import (
	"encoding/json"
	"strings"
)

type Tag struct {
	Name              TagName
	Title             string
	Description       string
	InteractionGroups map[Protocol]TagInteractionGroup
	Children          *Tags
}

var _ json.Marshaler = &Tags{}

func newEmptyTag(r InteractionId) *Tag {
	title := tagTitle(r.Path().String())
	return &Tag{
		InteractionGroups: make(map[Protocol]TagInteractionGroup),
		Children:          &Tags{},
		Title:             title,
		Name:              tagName(title),
	}
}

func tagTitle(path string) string {
	p := strings.Split(path, "/")
	for len(p) != 0 {
		if p[0] != "" && p[0] != "." {
			break
		}
		p = p[1:]
	}
	if len(p) == 0 {
		return "/"
	}
	return "/" + p[0]
}

func (t *Tag) appendInteractionId(k InteractionId) {
	list, ok := t.InteractionGroups[k.Protocol()]
	if !ok {
		list = newTagInteractionGroup(k.Protocol())
		t.InteractionGroups[k.Protocol()] = list
	}
	list.append(k)
}

func (t *Tag) MarshalJSON() ([]byte, error) {
	var data struct {
		Name              TagName               `json:"name"`
		Title             string                `json:"title"`
		Description       string                `json:"description,omitempty"`
		InteractionGroups []TagInteractionGroup `json:"interactionGroups"`
		Children          *Tags                 `json:"children,omitempty"`
	}

	data.Name = t.Name
	data.Title = t.Title
	data.Description = t.Description
	data.InteractionGroups = make([]TagInteractionGroup, 0, 2)
	for _, g := range t.InteractionGroups {
		data.InteractionGroups = append(data.InteractionGroups, g)
	}
	if t.Children.Len() > 0 {
		data.Children = t.Children
	}

	return json.Marshal(data)
}
