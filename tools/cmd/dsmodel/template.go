package main

const codeTemplate = `// Code generated by "dsmodel -type={{.Type}}"; DO NOT EDIT

package {{.Package}}

import(
    {{range $key, $as := .Dependencies -}}
    {{if $as -}}
    {{$as}} "{{$key}}"
    {{else -}}
    "{{$key}}"
    {{end -}}
    {{end }}
)

type {{.Type}}Kind struct {
    useDefaultIfNil bool
}

const {{.Type}}KindLoggerKey = "dsmodel.{{snakecase .Type}}"

func (k *{{.Type}}Kind) New() *{{.Type}} {
    a := &{{.Type}}{}
    {{range .Fields -}}{{if .Default -}}
    a.{{.FieldName}} = {{.Default}}
    {{end}}{{end -}}
    return a
}

func (k *{{.Type}}Kind) UseDefaultIfNil(b bool) *{{.Type}}Kind {
    k.useDefaultIfNil = b
    return k
}

// Get gets the kind entity from datastore
func (k *{{.Type}}Kind) Get(ctx context.Context, key interface{}) (*datastore.Key, *{{.Type}}, error) {
    keys, ents, err := k.GetMulti(ctx, key)
    if err != nil {
        return nil, nil, err
    }
    return keys[0], ents[0], nil
}

// GetMulti do Get with multiple keys
func (k *{{.Type}}Kind) GetMulti(ctx context.Context, keys ...interface{}) ([]*datastore.Key, []*{{.Type}}, error) {
    logger := xlog.WithContext(ctx).WithKey({{.Type}}KindLoggerKey)
    size := len(keys)
    if size == 0 {
        return nil, nil, nil
    }
    dsKeys := make([]*datastore.Key, size, size)
    for i := range keys {
        dsKeys[i] = helper.NewKey(ctx, "{{.Type}}", keys[i])
    }
    ents := make([]*{{.Type}}, size, size)
    err := helper.GetMulti(ctx, dsKeys, ents)
    if helper.IsDatastoreError(err) {
        return nil, nil, err
    }

    if k.useDefaultIfNil {
        for i:=0; i < size; i++ {
            if ents[i] == nil {
                ents[i] = k.New()
                ents[i].{{.IDField}} = dsKeys[i].StringID() // TODO: Support non-string key as ID
            }
        }
    }

    logger.Debug(func(p *xlog.Printer){
        p.Println("{{.Type}}#GetMulti (UseDefault: %b)", k.useDefaultIfNil)
        for i:=0; i < size; i++ {
            s := fmt.Sprintf("%v", ents[i])
            if len(s) > 20 {
                p.Printf("\t%s - %s...\n", dsKeys[i], s[:20])
            }else{
                p.Printf("\t%s - %s\n", dsKeys[i], s)
            }
        }
    })

    return dsKeys, ents, nil
}
`
