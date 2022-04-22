package catalog

import (
	"encoding/json"
	"strings"
)

type Tag struct {
	Name            TagName
	Title           string
	Description     string
	ResourceMethods *TagResourceMethods
	Children        *Tags
}

var _ json.Marshaler = &Tags{}

func newEmptyTag(r ResourceMethodId) *Tag {
	t := Tag{
		ResourceMethods: &TagResourceMethods{},
		Children:        &Tags{},
	}
	t.Title = tagTitle(r.path.String())
	t.Name = tagName(t.Title)
	return &t
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

func (t *Tag) appendResourceMethodId(r ResourceMethodId) {
	list, ok := t.ResourceMethods.Get(r.path)
	if !ok {
		list = newResourceMethodIdList()
		t.ResourceMethods.Set(r.path, list)
	}
	list.append(r)
}

func (t *Tag) MarshalJSON() ([]byte, error) {
	var data struct {
		Name            TagName             `json:"name"`
		Title           string              `json:"title"`
		Description     string              `json:"description,omitempty"`
		ResourceMethods *TagResourceMethods `json:"resourceMethods"`
		Children        *Tags               `json:"children,omitempty"`
	}

	data.Name = t.Name
	data.Title = t.Title
	data.Description = t.Description
	data.ResourceMethods = t.ResourceMethods
	if t.Children.Len() > 0 {
		data.Children = t.Children
	}

	return json.Marshal(data)
}
