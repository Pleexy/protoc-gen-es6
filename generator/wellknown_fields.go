package generator

import (
	"fmt"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	pgs "github.com/lyft/protoc-gen-star"
	"strconv"
	"strings"
)

type timestampField struct {
	Field
}

func (t *timestampField) GenerateSerializeBlock(p Printer, val string) {
	p.Printf(`if (%s) {
  writer.writeMessage(%d, %s, %s)
}
`, t.CheckNotEmptyExp(val), t.Number(), val, t.SerializeFunction())
}
func (t *timestampField) GenerateDeserializeBlock(p Printer, val string) {
	p.Printf(`reader.readMessage(%s, %s);
`, val, t.DeserializeFunction() )
}

func (t *timestampField) SerializeFunction() string {
	return `(val, writer) => {
    const seconds = Math.floor(val.getTime() / 1000)
    const nanos = val.getMilliseconds() * 1000000
    if (seconds !== 0) {
      writer.writeInt64(1, seconds);  
    }
    if (nanos !== 0) {
      writer.writeInt32(2, nanos);
    }
  }`
}

func (t *timestampField) DeserializeFunction() string {
	return `(val, reader) => {
  let seconds = 0;
  let nanos = 0;
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    const fieldNumber = reader.getFieldNumber();
    switch (fieldNumber) {
    case 1:	
      seconds = /** @type {number} */ (reader.readInt64());
      break;
    case 2:
      nanos = /** @type {number} */ (reader.readInt32());
      break;
    default:
	  reader.skipField();
	  break;
    }
  }
 val.setDate((seconds * 1000) + (nanos / 1000000))
}`
}



type wellKnownField struct {
	Field
	WellKnownType string
	FromFunction string
	ToFunction string
}
/*
func NewTimeStampField(pgsField pgs.Field, msg *MessageGenerator, o *Options) (FieldGenerator, error) {
	tf := &timestampField{Field:  newField(pgsField, msg, o)}
	tf.protoType = "google.protobuf.Timestamp"
	tf.es6Type = "Date"
	tf.typeValidationFunc = func(name string) string { return fmt.Sprintf("%s == null || %s  instanceof Date", name, name)}
	tf.defaultValue = ""
	tf.checkEmptyFunc = func(name string) string { return name }
	return tf, nil
}
 */

func NewWellKnownField(pgsField pgs.Field, msg *MessageGenerator, o *Options) (FieldGenerator, error) {
	if  *pgsField.Descriptor().Type != descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		return nil, nil
	}
	if ! strings.HasPrefix(*pgsField.Descriptor().TypeName, ".google.protobuf.") {
		// unsupported well known type
		return nil, nil
	}
	/*if *pgsField.Descriptor().TypeName == ".google.protobuf.Timestamp" {
		return NewTimeStampField(pgsField, msg, o)
	} */

	m := &wellKnownField{Field:newField(pgsField,msg,o)}

	prefix, typeName, err := GetTypeNameAndPrefixForField(pgsField, msg.File)
	if err != nil {
		return nil, err
	}

	switch typeName {
	case "Timestamp":
		m.es6Type = "Date"
		m.typeValidationFunc = func(name string) string { return fmt.Sprintf("%s == null || %s  instanceof Date", name, name)}
		m.defaultValue = ""
		m.checkEmptyFunc = func(name string) string { return name }
		m.FromFunction = "fromDate"
		m.ToFunction = "toDate"
	case "Struct":
		m.es6Type = "{ [string]: any }"
		m.typeValidationFunc = func(val string) string { return "" }
		m.checkEmptyFunc = func(val string) string { return fmt.Sprintf("%s && %s.constructor === Object && Object.entries(%s).length > 0", val, val, val) }
		m.FromFunction = "fromJavaScript"
		m.ToFunction = "toJavaScript"
	case "Value":
		m.es6Type = "any"
		m.typeValidationFunc = func(val string) string { return "" }
		m.checkEmptyFunc = func(val string) string { return val }
		m.FromFunction = "fromJavaScript"
		m.ToFunction = "toJavaScript"

	default:
		return nil, nil // process as usual message field
	}
	m.WellKnownType = prefix + typeName
	return m, nil
}

func (m *wellKnownField) GenerateSerializeBlock(p Printer, val string) {
	p.Printf(`if (%s) {
  writer.writeMessage(%d, %s, %s)
}
`, m.CheckNotEmptyExp(val), m.Number(), val, m.SerializeFunction())
}
func (m *wellKnownField) GenerateDeserializeBlock(p Printer, val string) {
	p.Print(m.DeserializeBlock(val)+"\n")
}

func (m *wellKnownField) SerializeFunction() string {
	return fmt.Sprintf("(val, writer) => %s.serializeBinaryToWriter(%s.%s(val), writer)",
		m.WellKnownType, m.WellKnownType,m.FromFunction)
}

func (m *wellKnownField) DeserializeBlock(valName string) string {
	varName := "tmpWKT"+ strconv.Itoa(int(m.Number()))
	return fmt.Sprintf(`const %s = new %s();
reader.readMessage(%s, %s.deserializeBinaryFromReader);
%s = %s.%s();`,varName, m.WellKnownType, varName, m.WellKnownType, valName, varName, m.ToFunction)
}