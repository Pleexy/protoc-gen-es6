package generator

import (
	"fmt"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	pgs "github.com/lyft/protoc-gen-star"
	"strings"
)

type EnumField struct {
	Field
}

func NewEnumFieldGenerator(pgsField pgs.Field, msg *MessageGenerator, o *Options) (FieldGenerator, error) {
	if *pgsField.Descriptor().Type != descriptor.FieldDescriptorProto_TYPE_ENUM {
		return nil, nil
	}
	e := &EnumField{
		Field:newField(pgsField, msg, o),
	}
	prefix, typeName, err := GetTypeNameAndPrefixForField(pgsField, msg.File)
	if err != nil {
		return nil, err
	}
	e.es6Type = fmt.Sprintf("$Values<typeof %s%s>", prefix, typeName)
	e.protoType = strings.TrimPrefix(pgsField.Descriptor().GetTypeName(), ".")
	e.checkEmptyFunc = func (val string) string { return fmt.Sprintf("%s !== 0", val)}
	e.typeValidationFunc = func (val string) string { return fmt.Sprintf("typeof %s === 'number'", val)}
	return e, nil
}

func (m *EnumField) GenerateSerializeBlock(p Printer, val string) {
p.Printf(`if (%s){
  writer.writeEnum(%d, %s)
}
`, m.checkEmptyFunc(val), m.Number(), val)
}
func (m *EnumField) GenerateDeserializeBlock(p Printer, val string) {
	p.Printf("%s = reader.readEnum();", val)
}

func (m*EnumField) SerializeFunction() string {
	return "writeEnum"
}

func (m *EnumField) DeserializeBlock(valName string) string {
	return fmt.Sprintf("%s = reader.readEnum();", valName)
}