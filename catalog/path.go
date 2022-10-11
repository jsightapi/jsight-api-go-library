package catalog

import "strings"

type Path string

func (p Path) Similar(urlPath string) bool {
	pp := strings.Split(p.String(), "/")
	uu := strings.Split(urlPath, "/")
	if len(pp) != len(uu) {
		return false
	}
	for i, v := range pp {
		if v != uu[i] {
			if len(v) == 0 || v[0] != '{' {
				return false
			}
		}
	}
	return true
}

func (p Path) String() string {
	return string(p)
}
