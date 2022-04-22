package catalog

import (
	dd "bytes"
	"encoding/json"
)

type Resource struct {
	Key ResourceMethodId
	Val *ResourceMethod
}

type OrderedResources []Resource

func (oo OrderedResources) MarshalJSON() ([]byte, error) {
	var buf dd.Buffer
	buf.WriteRune('{')

	for i, kv := range oo {
		if i != 0 {
			buf.WriteRune(',')
		}

		// marshal key
		key, err := json.Marshal(kv.Key)
		if err != nil {
			return nil, err
		}
		buf.Write(key)
		buf.WriteRune(':')

		// marshal value
		val, err := json.Marshal(kv.Val)
		if err != nil {
			return nil, err
		}
		buf.Write(val)
	}

	buf.WriteRune('}')
	return buf.Bytes(), nil
}
