package generator

import (
	"fmt"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	pgs "github.com/lyft/protoc-gen-star"
	"strings"
)

type MessageField struct {
	Field
}

func NewMessageFieldGenerator(pgsField pgs.Field, msg *MessageGenerator, o *Options) (FieldGenerator, error) {
	if  *pgsField.Descriptor().Type != descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		return nil, nil
	}
	f, err :=  NewWellKnownField(pgsField, msg, o)
	if err != nil || f != nil {
		return f, err
	}
	prefix, typeName, err := GetTypeNameAndPrefixForField(pgsField, msg.File)
	if err != nil {
		return nil, err
	}
	mf := &MessageField{
		Field: newField(pgsField, msg, o),
	}
	mf.es6Type = prefix + typeName
	mf.protoType = strings.TrimPrefix(pgsField.Descriptor().GetTypeName(), ".")
	mf.checkEmptyFunc = func (name string ) string { return name }
	mf.typeValidationFunc = func (name string) string { return fmt.Sprintf("%s == null || %s instanceof %s", name, name, mf.es6Type) }
	return mf, nil
}

func (m *MessageField) GenerateFromObject(p Printer, destField, srcField string) {
	p.Printf("%s = %s.fromObject(%s);\n", destField, m.es6Type, srcField)
}

func (m *MessageField) FromObjectExp(src string) string {
	return fmt.Sprintf("%s.fromObject(%s)", m.es6Type, src)
}

func (m *MessageField) ToObjectExp(src string) string {
	return fmt.Sprintf("%s.toObject()", src)
}

func (m *MessageField) GenerateSerializeBlock(p Printer, val string) {
	p.Printf(`if (%s){
  writer.writeMessage(%d, %s, %s.serializeBinaryToWriter)
}
`, m.checkEmptyFunc(val), m.Number(), val, m.es6Type)
}
func (m *MessageField) GenerateDeserializeBlock(p Printer, val string) {
	p.Printf(`%s = new %s();
reader.readMessage(%s, %s.deserializeBinaryFromReader);
`, val, m.es6Type, val, m.es6Type)
}

func (m*MessageField) SerializeFunction() string {
	return fmt.Sprintf(" %s.serializeBinaryToWriter", m.es6Type)
}

func (m *MessageField) DeserializeBlock(valName string) string {
	return fmt.Sprintf(`%s = new %s();
reader.readMessage(%s, %s.deserializeBinaryFromReader);`, valName, m.es6Type, valName, m.es6Type)
}
