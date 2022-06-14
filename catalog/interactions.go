package catalog

import (
	"sync"
)

// Interactions represent available resource methods.
// gen:OrderedMap
type Interactions struct {
	data  map[InteractionId]*HttpInteraction
	order []InteractionId
	mx    sync.RWMutex
}
