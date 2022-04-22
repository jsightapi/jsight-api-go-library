package core

import "j/japi/directive"

func removeDirectiveFromSlice(slice []*directive.Directive, i int) []*directive.Directive {
	return append(slice[:i], slice[i+1:]...)
}
