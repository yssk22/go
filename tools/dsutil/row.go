package dsutil

import (
	"fmt"

	"github.com/speedland/go/x/xtime"
	"google.golang.org/appengine/datastore"
)

type Row struct {
	Key  *datastore.Key
	Data entity
}

type entity map[string]interface{}

func (ent entity) Load(props []datastore.Property) error {
	for _, p := range props {
		obj := map[string]interface{}{
			"Value":   p.Value,
			"NoIndex": p.NoIndex,
		}
		if p.Multiple {
			if _, ok := ent[p.Name]; !ok {
				ent[p.Name] = make([]interface{}, 0)
			}
			ent[p.Name] = append(ent[p.Name].([]interface{}), obj)
		} else {
			ent[p.Name] = obj
		}
	}
	return nil
}

func (ent entity) Save() ([]datastore.Property, error) {
	var props []datastore.Property
	for k, v := range ent {
		switch v.(type) {
		case map[string]interface{}:
			prop, err := ent.mapToProp(v.(map[string]interface{}))
			if err != nil {
				return nil, err
			}
			prop.Name = k
			prop.Multiple = false
			props = append(props, *prop)
			break
		case []interface{}:
			for _, vv := range v.([]interface{}) {
				vvv, ok := vv.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("unexpected value: %v", vv)
				}
				prop, err := ent.mapToProp(vvv)
				prop.Name = k
				prop.Multiple = true
				if err != nil {
					return nil, err
				}
				props = append(props, *prop)
			}
			break
		default:
			return nil, fmt.Errorf("unexpected value: %v", v)
		}
	}
	return props, nil
}

func (ent entity) mapToProp(m map[string]interface{}) (*datastore.Property, error) {
	prop := &datastore.Property{}
	noIndex, ok := m["NoIndex"].(bool)
	if !ok {
		return nil, fmt.Errorf("missing `NoIndex` key")
	}
	prop.NoIndex = noIndex
	value, ok := m["Value"]
	if !ok {
		return nil, fmt.Errorf("missing `Value` key")
	}
	switch value.(type) {
	case string:
		if v, err := xtime.Parse(value.(string)); err == nil {
			prop.Value = v
		} else {
			prop.Value = value
		}
		break
	default:
		prop.Value = value
		break
	}
	return prop, nil
}
