package generator

import (
	"fmt"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	pgs "github.com/lyft/protoc-gen-star"
)

type PrimitiveField struct {
	Field
}

func NewPrimitiveFieldWriter(pgsField pgs.Field, msg *MessageGenerator, o *Options) (FieldGenerator, error) {
	s := &PrimitiveField{Field: newField(pgsField, msg, o)}
	switch *pgsField.Descriptor().Type {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE,
		descriptor.FieldDescriptorProto_TYPE_FLOAT,
		descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_UINT64,
		descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_FIXED64,
		descriptor.FieldDescriptorProto_TYPE_FIXED32,
		descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_SFIXED64,
		descriptor.FieldDescriptorProto_TYPE_SFIXED32,
		descriptor.FieldDescriptorProto_TYPE_SINT32,
		descriptor.FieldDescriptorProto_TYPE_SINT64:
		s.es6Type = "number"
		s.checkEmptyFunc = func(name string) string { return fmt.Sprintf("%s && %s !== 0", name, name) }
		s.defaultValue = "0"
		s.typeValidationFunc = func(name string ) string {return fmt.Sprintf("%s == null || typeof %s === 'number'", name, name) }
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		s.es6Type = "Uint8Array"
		s.checkEmptyFunc = func(name string) string { return fmt.Sprintf("%s && %s.length > 0", name, name) }
		s.typeValidationFunc = func(name string) string { return fmt.Sprintf("%s == null || %s instanceof Uint8Array", name, name) }
		s.defaultValue = ""
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		s.es6Type = "string"
		s.checkEmptyFunc = func(name string) string { return fmt.Sprintf(" %s && %s.length > 0", name, name) }
		s.defaultValue = "''"
		s.typeValidationFunc = func(name string) string { return fmt.Sprintf("%s == null || typeof %s === 'string' || %s instanceof String", name, name, name) }
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		s.es6Type = "boolean"
		s.checkEmptyFunc = func(name string) string { return name }
		s.defaultValue = "false"
		s.typeValidationFunc = func(n string) string { return fmt.Sprintf("%s == null || typeof %s === 'boolean' || (typeof %s === 'object' &&  %s !== null && typeof %s.valueOf() === 'boolean')", n, n, n, n,n) }
	default:
		return nil, nil
	}
	s.protoType = ReadWriteNames[*pgsField.Descriptor().Type]
	return s, nil
}

func (pf *PrimitiveField) SerializeFunction() string {
	return "write" + upperFirst(pf.protoType)
}

func (pf *PrimitiveField) GenerateSetter(p Printer) {
	if pf.es6Type != "string" || ! pf.o.ConvertString {
		pf.Field.GenerateSetter(p)
		return
	}
	p.PrintTpl("string_set",`/** @param val string */
set {{.setName}}(val{{- if .flow -}} :string):void{{- else -}}){{- end -}}{
  if (val == null || typeof val === 'string' || val instanceof String) {
    this.{{.prop}} = val;
  } else {
     let res = String(val);
     if (res === '[object Object]') {
       res = JSON.stringify(val);
     }
     this.{{.prop}} = res;
  }
}
`,"setName",pf.GetSetName(), "flow", pf.o.Flow, "prop",  pf.PropertyName())
}

func (pf *PrimitiveField) GenerateSerializeBlock(p Printer, val string) {
	p.Printf(`if (%s){
  writer.%s(%d, %s); 
}
`, pf.CheckNotEmptyExp(val), pf.SerializeFunction() , pf.Number(), val)
}

func (pf *PrimitiveField) GenerateDeserializeBlock(p Printer, val string) {
	p.Printf( "%s = reader.%s();\n", val, "read" + upperFirst(pf.protoType))
}

func (pf *PrimitiveField) DeserializeBlock(valName string) string {
	return fmt.Sprintf("%s = reader.%s();", valName, "read" + upperFirst(pf.protoType))
}