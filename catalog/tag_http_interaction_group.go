package catalog

import (
	"encoding/json"
)

type TagHttpInteractionGroup struct {
	protocol     Protocol
	interactions []InteractionId
}

func newTagHttpInteractionGroup() *TagHttpInteractionGroup {
	return &TagHttpInteractionGroup{
		protocol:     http,
		interactions: make([]InteractionId, 0, 5),
	}
}

func (l *TagHttpInteractionGroup) append(i InteractionId) {
	l.interactions = append(l.interactions, i)
}

func (l TagHttpInteractionGroup) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.interactions)
}
