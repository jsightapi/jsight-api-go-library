package catalog

import (
	"encoding/json"
)

type TagJsonRpcInteractionGroup struct {
	protocol     Protocol
	interactions []InteractionId
}

func newTagJsonRpcInteractionGroup() *TagJsonRpcInteractionGroup {
	return &TagJsonRpcInteractionGroup{
		protocol:     jsonRpc,
		interactions: make([]InteractionId, 0, 5),
	}
}

func (l *TagJsonRpcInteractionGroup) append(i InteractionId) {
	l.interactions = append(l.interactions, i)
}

func (l TagJsonRpcInteractionGroup) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.interactions)
}
