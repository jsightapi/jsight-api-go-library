package catalog

import (
	"regexp"
	"strings"
)

func Annotation(s string) string {
	space := regexp.MustCompile(`\s+`)
	return space.ReplaceAllString(strings.TrimSpace(s), " ")
}
