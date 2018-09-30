package gaetest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"text/template"

	"github.com/yssk22/go/x/xlog"
	"github.com/yssk22/go/x/xtime"

	"context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
)

var _floatRe = regexp.MustCompile("\\.0+$")

const FixtureLoggerKey = "web.gae.gaetest.fixture"

// CleanupStorage cleans up all storage services (memcache and datastore)
func CleanupStorage(ctx context.Context, namespaces ...string) error {
	if err := ResetMemcache(ctx, namespaces...); err != nil {
		return err
	}
	if err := CleanupDatastore(ctx, namespaces...); err != nil {
		return err
	}
	return nil
}

// ResetMemcache resets all memcache entries
func ResetMemcache(ctx context.Context, namespaces ...string) error {
	if len(namespaces) == 0 {
		return memcache.Flush(ctx)
	}
	for _, ns := range namespaces {
		var _ctx = ctx
		if ns != "" {
			_ctx, _ = appengine.Namespace(_ctx, ns)
		}
		if err := memcache.Flush(_ctx); err != nil {
			return err
		}
	}
	return nil
}

// ResetFixtureFromFile is to reset all data in datastore and reload the fixtures from the file.
func ResetFixtureFromFile(ctx context.Context, path string, bindings interface{}, namespaces ...string) error {
	if err := CleanupDatastore(ctx, namespaces...); err != nil {
		return err
	}
	err := FixtureFromFile(ctx, path, bindings)
	return err
}

// CleanupDatastore is to remove all data in datastore
func CleanupDatastore(ctx context.Context, namespaces ...string) error {
	ctx, logger := xlog.WithContextAndKey(ctx, "", FixtureLoggerKey)
	if len(namespaces) == 0 {
		numDeleted, err := cleanupDatastore(ctx)
		if err != nil {
			return err
		}
		logger.Infof("Clean up %d ents", numDeleted)
		return nil
	}
	for _, ns := range namespaces {
		var _ctx = ctx
		if ns != "" {
			_ctx, _ = appengine.Namespace(_ctx, ns)
		}
		numDeleted, err := cleanupDatastore(_ctx)
		if err != nil {
			return err
		}
		logger.Infof("Clean up %d ents in %q", numDeleted, ns)
	}
	return nil
}

func cleanupDatastore(ctx context.Context) (int, error) {
	var dummy []interface{}
	numDeleted := 0
	count := 1
	for {
		var keys []*datastore.Key
		var err error
		if keys, err = datastore.NewQuery("").KeysOnly().GetAll(ctx, dummy); err != nil {
			return 0, err
		}
		if err := datastore.DeleteMulti(ctx, keys); err != nil {
			return 0, err
		}
		numDeleted += len(keys)
		count, _ = datastore.NewQuery("").KeysOnly().Count(ctx)
		if count == 0 {
			return numDeleted, nil
		}
	}
}

// FixtureFromMap is to load fixtures from []map[string]interface{}
func FixtureFromMap(ctx context.Context, arr []map[string]interface{}) error {
	for _, v := range arr {
		var vv map[string]interface{}
		buff, _ := json.Marshal(v)
		json.Unmarshal(buff, vv)
		if err := loadJsonToDatastore(ctx, nil, v); err != nil {
			return err
		}
	}
	return nil
}

// FixtureFromFile is to load fixtures from a file.
func FixtureFromFile(ctx context.Context, path string, bindings interface{}) error {
	return DatastoreFixture(ctx, path, bindings)
}

// DatastoreFixture loads the fixtures from path to datastore.
func DatastoreFixture(ctx context.Context, path string, bindings interface{}) error {
	data, err := loadFile(path, bindings)
	if err != nil {
		return fmt.Errorf("Could not load fixture file from %s: %v", err)
	}
	var arr []map[string]interface{}
	if err = json.Unmarshal(data, &arr); err != nil {
		return fmt.Errorf("Could not load the json file from %q - JSON Parse error: %v", path, err)
	}
	for _, v := range arr {
		if err := loadJsonToDatastore(ctx, nil, v); err != nil {
			return err
		}
	}
	return nil
}

type jsonSaver map[string]interface{}

func (js jsonSaver) Load(ps []datastore.Property) error {
	return nil
}

func (js jsonSaver) Save() ([]datastore.Property, error) {
	props := []datastore.Property{}
	for k, v := range map[string]interface{}(js) {
		if !strings.HasPrefix(k, "_") {
			for _, val := range convertJsonValueToProperties(k, v) {
				props = append(props, val)
			}
		}
	}
	return props, nil
}

func loadFile(path string, bindings interface{}) ([]byte, error) {
	t, err := template.New(filepath.Base(path)).ParseFiles(path)
	if err != nil {
		return nil, err
	}
	var buff bytes.Buffer
	err = t.Execute(&buff, bindings)
	return buff.Bytes(), err
}

func convertJsonValueToProperties(k string, v interface{}) []datastore.Property {
	var propertyList []datastore.Property
	var value = reflect.ValueOf(v)

	switch value.Kind() {
	case reflect.String:
		p := datastore.Property{Name: k}
		s := v.(string)
		if strings.HasPrefix(s, "[]") {
			p.Value = []byte(strings.TrimPrefix(s, "[]"))
			p.NoIndex = true
			propertyList = append(propertyList, p)
		} else {
			if dt, err := xtime.Parse(fmt.Sprintf("%s", v)); err == nil {
				p.Value = dt
				propertyList = append(propertyList, p)
			} else if d, err := xtime.Parse(fmt.Sprintf("%sT00:00:00Z", v)); err == nil {
				p.Value = d
				propertyList = append(propertyList, p)
			} else {
				p.Value = s
				propertyList = append(propertyList, p)
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// reach here from FixtureFromMap since it can contain non floating number.
		var vv int64
		switch v.(type) {
		case int:
			vv = int64(v.(int))
		case int8:
			vv = int64(v.(int8))
		case int16:
			vv = int64(v.(int16))
		case int32:
			vv = int64(v.(int32))
		case int64:
			vv = v.(int64)
		}
		propertyList = append(propertyList, datastore.Property{
			Name:  k,
			Value: vv,
		})
	case reflect.Float32, reflect.Float64:
		str := []byte(fmt.Sprintf("%f", v))
		if _floatRe.Match(str) {
			// should be int.
			propertyList = append(propertyList, datastore.Property{
				Name:  k,
				Value: int64(v.(float64)),
			})
		} else {
			propertyList = append(propertyList, datastore.Property{
				Name:  k,
				Value: v,
			})
		}
	case reflect.Bool:
		propertyList = append(propertyList, datastore.Property{
			Name:  k,
			Value: v,
		})
	case reflect.Map:
		for k1, v1 := range v.(map[string]interface{}) {
			if !strings.HasPrefix(k1, "_") {
				for _, val := range convertJsonValueToProperties(k1, v1) {
					propertyList = append(propertyList, datastore.Property{
						Name:     fmt.Sprintf("%s.%s", k, val.Name),
						Value:    val.Value,
						Multiple: val.Multiple,
					})
				}
			}
		}
		propertyList = append(propertyList)
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			propertyList = append(propertyList, datastore.Property{
				Name:     k,
				Value:    value.Index(i).Interface(),
				Multiple: true,
			})
		}
	default:
		break
	}
	return propertyList
}

func loadJsonToDatastore(ctx context.Context, pkey *datastore.Key, data map[string]interface{}) error {
	ctx, logger := xlog.WithContextAndKey(ctx, "", FixtureLoggerKey)
	var kind string
	var ns string
	var keyval interface{}
	var key *datastore.Key
	var ok bool
	var err error
	if _, ok = data["_kind"]; !ok {
		return fmt.Errorf("Missing key `_kind`")
	}
	kind = data["_kind"].(string)
	if keyval, ok = data["_key"]; !ok {
		return fmt.Errorf("Missing key `_key`")
	}
	if _, ok = data["_ns"]; ok {
		ns = data["_ns"].(string)
		ctx, err = appengine.Namespace(ctx, ns)
		if err != nil {
			return fmt.Errorf("Could not change the namespace of %q, check _ns value: ", ns, err)
		}
	}

	switch keyval.(type) {
	case int64:
		key = datastore.NewKey(ctx, kind, "", keyval.(int64), pkey)
	case string:
		key = datastore.NewKey(ctx, kind, keyval.(string), 0, pkey)
	default:
		return fmt.Errorf("Invalid `_key` type.")
	}
	if _, err := datastore.Put(ctx, key, jsonSaver(data)); err != nil {
		return err
	}
	logger.Infof("Fixture: %s loaded", key)
	if children, ok := data["_children"]; ok {
		for _, v := range children.([]interface{}) {
			if err := loadJsonToDatastore(ctx, key, v.(map[string]interface{})); err != nil {
				return err
			}
		}
	}
	return nil
}
