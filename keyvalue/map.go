package keyvalue

import "fmt"

// Map is an alias for map[interface{}]interface{} that implements GetterSetter interface
type Map map[interface{}]interface{}

// Get implements Getter#Get
func (m Map) Get(key interface{}) (interface{}, error) {
	if v, ok := m[key]; ok {
		return v, nil
	}
	return nil, KeyError(fmt.Sprintf("%s", key))
}

// Set implements Setter#Set
func (m Map) Set(key interface{}, v interface{}) error {
	m[key] = v
	return nil
}

// NewMap returns a new Map (shorthand for `Map(make(map[string]interface{}))`)
func NewMap() Map {
	return Map(make(map[interface{}]interface{}))
}

// StringKeyMap is an alias for map[string]interface{} that implements GetterSetter interface
type StringKeyMap map[string]interface{}

// Get implements Getter#Get
func (m StringKeyMap) Get(key interface{}) (interface{}, error) {
	if v, ok := m[key.(string)]; ok {
		return v, nil
	}
	return nil, KeyError(fmt.Sprintf("%s", key))
}

// Set implements Setter#Set
func (m StringKeyMap) Set(key interface{}, v interface{}) error {
	m[key.(string)] = v
	return nil
}

// Del implements Setter#Del
func (m StringKeyMap) Del(key interface{}) error {
	delete(m, key.(string))
	return nil
}

// NewStringKeyMap returns a new Map (shorthand for `Map(make(map[string]interface{}))`)
func NewStringKeyMap() StringKeyMap {
	return StringKeyMap(make(map[string]interface{}))
}
