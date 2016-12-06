package dsutil

import "google.golang.org/appengine/datastore"

type entityLoader map[string]interface{}

func (l entityLoader) Load(props []datastore.Property) error {
	for _, p := range props {
		obj := map[string]interface{}{
			"Value":   p.Value,
			"NoIndex": p.NoIndex,
		}
		if p.Multiple {
			if _, ok := l[p.Name]; !ok {
				l[p.Name] = make([]interface{}, 0)
			}
			l[p.Name] = append(l[p.Name].([]interface{}), obj)
		} else {
			l[p.Name] = obj
		}
	}
	return nil
}

func (l entityLoader) Save() ([]datastore.Property, error) {
	// unused
	return nil, nil
}

type entityValueLoader map[string]interface{}

func (l entityValueLoader) Load(props []datastore.Property) error {
	for _, p := range props {
		if p.Multiple {
			if _, ok := l[p.Name]; !ok {
				l[p.Name] = make([]interface{}, 0)
			}
			l[p.Name] = append(l[p.Name].([]interface{}), p.Value)
		} else {
			l[p.Name] = p.Value
		}
	}
	return nil
}

func (l entityValueLoader) Save() ([]datastore.Property, error) {
	// unused
	return nil, nil
}
