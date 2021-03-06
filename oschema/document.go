//
// Key schema struct, constructors that are part of the
// OrientDB schema or support representing the schema.
//
package oschema

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/quux00/ogonori/ogl"
)

type ODocument struct {
	RID        ORID
	Version    int32
	entryOrder []string           // field names in the order they were added to the ODocument
	Fields     map[string]*OField // key: property-name
	// TODO: may want a mapping of ids => OField
	Classname string // TODO: probably needs to change *OClass (once that is built)

	// private fields
	dirty bool // TODO: probably need to track which fields are dirty
}

//
// NewDocument should be called to create new ODocument objects,
// since some internal data structures need to be initialized
// before the ODocument is ready to use.
//
func NewDocument(className string) *ODocument {
	return &ODocument{
		Fields:    make(map[string]*OField),
		Classname: className,
		RID:       NewORID(),
		Version:   int32(-1),
	}
}

//
// TODO: have this replace NewDocument and change NewDocument to take RID and Version (???)
//
func NewEmptyDocument() *ODocument {
	return &ODocument{
		Fields:  make(map[string]*OField),
		RID:     NewORID(),
		Version: int32(-1),
	}
}

//
// Implements database/sql.Scanner interface
//
func (doc *ODocument) Scan(src interface{}) error {
	ogl.Debugln("** ODocument.Scan")

	locdoc := src.(*ODocument)
	*doc = *locdoc

	// switch src.(type) {
	// case *ODocument:
	// 	locdoc := src.(*ODocument)
	// 	*doc = *locdoc
	// default:
	// 	return errors.New("Say what???")
	// }
	return nil
}

//
// Implements database/sql/driver.Valuer interface
// TODO: haven't detected when this is called yet (probably when serializing ODocument for insertion into DB??)
//
func (doc *ODocument) Value() (driver.Value, error) {
	ogl.Debugln("** ODocument.Value")

	return []byte(`{"b": 2}`), nil // FIXME: bogus
}

//
// Implements database/sql/driver.ValueConverter interface
// TODO: haven't detected when this is called yet
//
func (doc *ODocument) ConvertValue(v interface{}) (driver.Value, error) {
	ogl.Printf("** ODocument.ConvertValue: %T: %v\n", v, v)

	return []byte(`{"a": 1}`), nil // FIXME: bogus
}

//
// Testing out JSON marshalling -> this method may change to something else
//
func (doc *ODocument) ToJSON() ([]byte, error) {
	return json.Marshal(doc)
}

//
// FieldNames returns the names of all the fields currently in this ODocument.
// These fields may not have already been committed to the database.
//
func (doc *ODocument) FieldNames() []string {
	names := make([]string, 0, len(doc.entryOrder))
	for _, name := range doc.entryOrder {
		names = append(names, name)
	}
	return names
}

//
// GetFieldById looks up the OField in this document with the specified field id
// (aka property-id). If no field is found with that id, nil is returned.
//
func (doc *ODocument) GetFieldById(id int32) *OField {
	for _, fld := range doc.Fields {
		if fld.Id == id {
			return fld
		}
	}
	return nil
}

//
// GetFieldByName looks up the OField in this document with the specified field.
// If no field is found with that name, nil is returned.
//
func (doc *ODocument) GetField(fname string) *OField {
	return doc.Fields[fname]
}

//
// AddField adds a fully created field directly rather than by some of its
// attributes, as the other "Field" methods do.
// The same *ODocument is returned to allow call chaining.
//
func (doc *ODocument) AddField(name string, field *OField) *ODocument {
	ogl.Debugf("ODocument.AddField name== %v\n", name)   // DEBUG
	ogl.Debugf("ODocument.AddField field== %v\n", field) // DEBUG
	doc.Fields[name] = field
	doc.entryOrder = append(doc.entryOrder, name)
	doc.dirty = true
	return doc
}

func (doc *ODocument) SetDirty(b bool) {
	doc.dirty = b
}

//
// Field is used to add a new field to a document. This will usually be done just
// before calling Save and sending it to the database.  The field type will be inferred
// via type switch analysis on `val`.  Use FieldWithType to specify the type directly.
// The same *ODocument is returned to allow call chaining.
//
func (doc *ODocument) Field(name string, val interface{}) *ODocument {
	// TODO: need to add more types: EMBEDDEDRECORD, EMBEDDEDLIST, LINK, etc. ...
	var ftype byte
	switch val.(type) {
	case string:
		ftype = STRING
	case bool:
		ftype = BOOLEAN
	case int:
		ftype = INTEGER // TODO: need to check the size of the int => might be LONG
	case int32:
		ftype = INTEGER
	case int64:
		ftype = LONG
	case int16:
		ftype = SHORT
	case int8:
		ftype = BYTE
	case float32:
		ftype = FLOAT
	case float64:
		ftype = DOUBLE
	case []byte:
		ftype = BINARY
	case OEmbeddedMap:
		ftype = EMBEDDEDMAP
		// TODO: more types need to be added
	default:
		ftype = ANY // TODO: no idea if this is correct
	}
	return doc.FieldWithType(name, val, ftype)
}

//
// FieldWithType is used to add a new field to a document. This will usually be done just
// before calling Save and sending it to the database. The `fieldType` must correspond
// one of the OrientDB type in the schema pkg constants.  It will follow the same list
// as: https://github.com/orientechnologies/orientdb/wiki/Types
// The same *ODocument is returned to allow call chaining.
//
func (doc *ODocument) FieldWithType(name string, val interface{}, fieldType byte) *ODocument {
	fld := &OField{
		Name:  name,
		Value: val,
		Typ:   fieldType,
	}
	return doc.AddField(name, fld)
}

//
// String implements Stringer interface
//
func (doc *ODocument) String() string {
	buf := new(bytes.Buffer)
	_, err := buf.WriteString(fmt.Sprintf("ODocument[Classname: %s; RID: %s; Version: %d; fields: \n",
		doc.Classname, doc.RID, doc.Version))
	if err != nil {
		panic(err)
	}

	for _, fld := range doc.Fields {
		_, err = buf.WriteString(fmt.Sprintf("  %s\n", fld.String()))
		if err != nil {
			panic(err)
		}
	}

	buf.Truncate(buf.Len() - 1)
	buf.WriteString("]\n")
	return buf.String()
}

//
// StringNoFields is a String() method that elides the fields.
// This is useful when the fields include links and there are
// circular links.
//
func (doc *ODocument) StringNoFields() string {
	return fmt.Sprintf("ODocument[Classname: %s; RID: %s; Version: %d; fields: [...]]",
		doc.Classname, doc.RID, doc.Version)
}
