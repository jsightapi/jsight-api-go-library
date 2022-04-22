package catalog

import (
	"j/japi/directive"
	"sync"
)

type Server struct {
	Annotation       string            `json:"annotation,omitempty"`
	BaseUrl          string            `json:"baseUrl"`
	BaseUrlVariables *baseUrlVariables `json:"baseUrlVariables,omitempty"`
}

type baseUrlVariables struct {
	Schema    *Schema             `json:"schema"`
	Directive directive.Directive `json:"-"`
}

// Servers represent available servers.
// gen:OrderedMap
type Servers struct {
	data  map[string]*Server
	order []string
	mx    sync.RWMutex
}
