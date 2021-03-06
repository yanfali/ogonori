package obinary

import (
	"bytes"

	"github.com/quux00/ogonori/oschema"
)

//
// ORecordSerializer is the interface for the binary Serializer/Deserializer.
// More than one implementation will be needed if/when OrientDB creates additional
// versions of the binary serialization format.
// TODO: may want to use this interface for the csv serializer also - if so need to move this interface up a level
//
type ORecordSerializer interface {
	//
	// Deserialize reads bytes from the bytes.Buffer and puts the data into the
	// ODocument object.  The ODocument must already be created; nil cannot be
	// passed in for the `doc` field.  The serialization version (the first byte
	// of the serialized record) should be stripped off (already read) from the
	// bytes.Buffer being passed in
	//
	Deserialize(dbc *DBClient, doc *oschema.ODocument, buf *bytes.Buffer) error // TODO: should this take an io.Reader instead of *bytes.Buffer ???

	//
	// Deserialize reads bytes from the bytes.Buffer and updates the ODocument object
	// passed in, but only for the fields specified in the `fields` slice.
	// The ODocument must already be created; nil cannot be passed in for the `doc` field.
	//
	DeserializePartial(doc *oschema.ODocument, buf *bytes.Buffer, fields []string) error

	//
	// Serialize reads the ODocument and serializes to bytes into the bytes.Buffer.
	//
	Serialize(doc *oschema.ODocument, buf *bytes.Buffer) error

	//
	// SerializeClass gets the class from the ODocument and serializes it to bytes
	// into the bytes.Buffer.
	//
	SerializeClass(doc *oschema.ODocument, buf *bytes.Buffer) error
}
