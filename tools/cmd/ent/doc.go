// ent tool generates the datastore access logics and
// validations for your datastore entities.
//
//     //go:generate -type=MyEnt
//     type MyEnt struct {
//        // pubic fields here
//     }
//
// then `go generate` command will generate my_ent_datastore.go with
// `MyEntKind` type and `MyEntQuery`` type that implements
// the datastore access logics.
//
// Some logics can be configured by using field tags in a model.
//
// 1. ent tag
//
//     type MyEnt struct {
//         ID        `ent:"id""`
//         UpdatedAt `ent:"timestamp"`
//     }
//
// ent tag has some attributes:
//
//     - 'id': used for the key in the datastore
//     - 'timestamp': automatically updated when kind puts the entity.
//     - 'form': indicates the field an be collected via {Type}Kind.FromForm function.
//
// 2. default tag.
//
//     type MyEnt struct {
//         ID        `ent:"id""`
//         Title     `default:"*Untitled*"`
//         UpdatedAt `ent:"timestamp"`
//     }
//
// default tag is a tag to configure default values of the instance
// created by {ModelName}.New() function. In the example above,
// `MyEnt.New()` would return the instance with the '*Untitled*'
// value of it's Title field.
//
package main
