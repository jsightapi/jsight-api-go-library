// Autogenerated code!
// DO NOT EDIT!
//
// Generated by OrderedMap generator from the internal/cmd/generator command.

package catalog

import (
	"bytes"
	"encoding/json"
)

// Set sets a value with specified key.
func (m *Tags) Set(k TagName, v *Tag) {
	m.mx.Lock()
	defer m.mx.Unlock()

	if m.data == nil {
		m.data = map[TagName]*Tag{}
	}
	if !m.has(k) {
		m.order = append(m.order, k)
	}
	m.data[k] = v
}

// SetToTop do the same as Set, but new key will be placed on top of the order
// map.
func (m *Tags) SetToTop(k TagName, v *Tag) {
	m.mx.Lock()
	defer m.mx.Unlock()

	if m.data == nil {
		m.data = map[TagName]*Tag{}
	}
	if !m.has(k) {
		m.order = append([]TagName{k}, m.order...)
	}
	m.data[k] = v
}

// Update updates a value with specified key.
func (m *Tags) Update(k TagName, fn func(v *Tag) *Tag) {
	m.mx.Lock()
	defer m.mx.Unlock()

	if !m.has(k) {
		// Prevent from possible nil pointer dereference if map value type is a
		// pointer.
		return
	}

	m.data[k] = fn(m.data[k])
}

// GetValue gets a value by key.
func (m *Tags) GetValue(k TagName) *Tag {
	m.mx.RLock()
	defer m.mx.RUnlock()

	return m.data[k]
}

// Get gets a value by key.
func (m *Tags) Get(k TagName) (*Tag, bool) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	v, ok := m.data[k]
	return v, ok
}

// Has checks that specified key is set.
func (m *Tags) Has(k TagName) bool {
	m.mx.RLock()
	defer m.mx.RUnlock()

	return m.has(k)
}

func (m *Tags) has(k TagName) bool {
	_, ok := m.data[k]
	return ok
}

// Len returns count of values.
func (m *Tags) Len() int {
	m.mx.RLock()
	defer m.mx.RUnlock()

	return len(m.data)
}

// Iterate iterates over map key/values.
// Will block in case of slow consumer.
// Should be used only for read only operations. Attempt to change something
// inside loop will lead to deadlock.
// Use Tags.Map when you have to update value.
func (m *Tags) Iterate() <-chan TagsItem {
	ch := make(chan TagsItem)
	go func() {
		m.mx.RLock()
		defer m.mx.RUnlock()

		for _, k := range m.order {
			ch <- TagsItem{
				Key:   k,
				Value: m.data[k],
			}
		}
		close(ch)
	}()
	return ch
}

// IterateReverse do the same as Iterate but in reverse order.
func (m *Tags) IterateReverse() <-chan TagsItem {
	ch := make(chan TagsItem)
	go func() {
		m.mx.RLock()
		defer m.mx.RUnlock()

		for i := len(m.order) - 1; i >= 0; i-- {
			k := m.order[i]
			ch <- TagsItem{
				Key:   k,
				Value: m.data[k],
			}
		}
		close(ch)
	}()
	return ch
}

// Map iterates and changes values in the map.
func (m *Tags) Map(fn mapTagsFunc) error {
	m.mx.Lock()
	defer m.mx.Unlock()

	for _, k := range m.order {
		v, err := fn(k, m.data[k])
		if err != nil {
			return err
		}
		m.data[k] = v
	}
	return nil
}

type mapTagsFunc = func(k TagName, v *Tag) (*Tag, error)

// TagsItem represent single data from the Tags.
type TagsItem struct {
	Key   TagName
	Value *Tag
}

var _ json.Marshaler = &Tags{}

func (m *Tags) MarshalJSON() ([]byte, error) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	var buf bytes.Buffer
	buf.WriteRune('{')

	for i, k := range m.order {
		if i != 0 {
			buf.WriteRune(',')
		}

		// marshal key
		key, err := json.Marshal(k)
		if err != nil {
			return nil, err
		}
		buf.Write(key)
		buf.WriteRune(':')

		// marshal value
		val, err := json.Marshal(m.data[k])
		if err != nil {
			return nil, err
		}
		buf.Write(val)
	}

	buf.WriteRune('}')
	return buf.Bytes(), nil
}
