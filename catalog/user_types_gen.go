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
func (m *UserTypes) Set(k string, v *UserType) {
	m.mx.Lock()
	defer m.mx.Unlock()

	if m.data == nil {
		m.data = map[string]*UserType{}
	}
	if !m.has(k) {
		m.order = append(m.order, k)
	}
	m.data[k] = v
}

// SetToTop do the same as Set, but new key will be placed on top of the order
// map.
func (m *UserTypes) SetToTop(k string, v *UserType) {
	m.mx.Lock()
	defer m.mx.Unlock()

	if m.data == nil {
		m.data = map[string]*UserType{}
	}
	if !m.has(k) {
		m.order = append([]string{k}, m.order...)
	}
	m.data[k] = v
}

// Update updates a value with specified key.
func (m *UserTypes) Update(k string, fn func(v *UserType) *UserType) {
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
func (m *UserTypes) GetValue(k string) *UserType {
	m.mx.RLock()
	defer m.mx.RUnlock()

	return m.data[k]
}

// Get gets a value by key.
func (m *UserTypes) Get(k string) (*UserType, bool) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	v, ok := m.data[k]
	return v, ok
}

// Has checks that specified key is set.
func (m *UserTypes) Has(k string) bool {
	m.mx.RLock()
	defer m.mx.RUnlock()

	return m.has(k)
}

func (m *UserTypes) has(k string) bool {
	_, ok := m.data[k]
	return ok
}

// Len returns count of values.
func (m *UserTypes) Len() int {
	m.mx.RLock()
	defer m.mx.RUnlock()

	return len(m.data)
}

// Find finds first matched item from the map.
func (m *UserTypes) Find(fn findUserTypesFunc) (UserTypesItem, bool) {
	for _, k := range m.order {
		if fn(k, m.data[k]) {
			return UserTypesItem{
				Key:   k,
				Value: m.data[k],
			}, true
		}
	}
	return UserTypesItem{}, false
}

type findUserTypesFunc = func(k string, v *UserType) bool

// Each iterates and perform given function on each item in the map.
func (m *UserTypes) Each(fn eachUserTypesFunc) error {
	for _, k := range m.order {
		if err := fn(k, m.data[k]); err != nil {
			return err
		}
	}
	return nil
}

// EachReverse act almost the same as Each but in reverse order.
func (m *UserTypes) EachReverse(fn eachUserTypesFunc) error {
	for i := len(m.order) - 1; i >= 0; i-- {
		k := m.order[i]
		if err := fn(k, m.data[k]); err != nil {
			return err
		}
	}
	return nil
}

type eachUserTypesFunc = func(k string, v *UserType) error

func (m *UserTypes) EachSafe(fn eachSafeUserTypesFunc) {
	for _, k := range m.order {
		fn(k, m.data[k])
	}
}

type eachSafeUserTypesFunc = func(k string, v *UserType)

// Map iterates and changes values in the map.
func (m *UserTypes) Map(fn mapUserTypesFunc) error {
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

type mapUserTypesFunc = func(k string, v *UserType) (*UserType, error)

// UserTypesItem represent single data from the UserTypes.
type UserTypesItem struct {
	Key   string
	Value *UserType
}

var _ json.Marshaler = &UserTypes{}

func (m *UserTypes) MarshalJSON() ([]byte, error) {
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
