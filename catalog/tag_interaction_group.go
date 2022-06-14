package catalog

type TagInteractionGroup interface {
	append(i InteractionId)
	MarshalJSON() ([]byte, error)
}

func newTagInteractionGroup(p Protocol) TagInteractionGroup {
	switch p {
	case jsonRpc:
		return newTagJsonRpcInteractionGroup()
	default: // case http:
		return newTagHttpInteractionGroup()
	}
}
