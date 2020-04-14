package generator

import (
	"fmt"
	pgs "github.com/lyft/protoc-gen-star"
)

type RepeatedField struct {
	Field
	Element FieldGenerator
}

func NewRepeatedFieldWriter(pgsField pgs.Field, msg *MessageGenerator, o *Options) (FieldGenerator, error) {
	if pgsField.Type().IsRepeated() {
		s := &RepeatedField{Field: newField(pgsField, msg, o)}
		el, err := CompositeFieldResolver(NewMessageFieldGenerator,NewEnumFieldGenerator,NewPrimitiveFieldWriter)(pgsField, msg, o)
		if err != nil {
			return nil, err
		}
		if el == nil {
			return nil, fmt.Errorf("cannot find element type for repeated field %s", pgsField.Name().String())
		}
		s.Element = el
		s.es6Type = "Array<" + el.ES6Type() + ">"
		s.checkEmptyFunc = func (name string) string { return fmt.Sprintf("%s && %s.length > 0", name, name) }
		s.typeValidationFunc = func (name string) string { return fmt.Sprintf("%s == null || Array.isArray(%s)", name, name) }
		s.defaultValue = ""
		s.protoType = "repeated " + el.ProtoType()
		return s, nil
	}
	return nil, nil
}

func (rf *RepeatedField) FromObjectExp(src string) string {
	return fmt.Sprintf("%s && Array.from(%s).map(x => %s)", src, src, rf.Element.FromObjectExp("x"))
}

func (rf *RepeatedField) ToObjectExp(src string) string {
	return fmt.Sprintf("%s && %s.map(x => %s)", src, src, rf.Element.ToObjectExp("x"))
}


func (rf *RepeatedField) GenerateSerializeBlock(p Printer, val string) {
	p.Printf("if (%s) {\n", rf.checkEmptyFunc(val))
	if rf.Element.IsMessage() {
		p.Printf("  writer.writeRepeatedMessage(%d, %s, %s);\n", rf.Number(), val, rf.Element.SerializeFunction())
	} else {
		method := "Repeated"
		if rf.IsPacked() {
			method = "Packed"
		}
		p.Printf("  writer.write%s%s(%d, %s);\n", method, upperFirst(ReadWriteNames[*rf.pgsField.Descriptor().Type]), rf.Number(), val)
	}
	p.Printf("}\n")
}

func (rf *RepeatedField) GenerateDeserializeBlock(p Printer, val string) {
	if rf.IsPacked() {
		p.Printf( "%s = reader.readPacked%s();\n", val, upperFirst(ReadWriteNames[*rf.pgsField.Descriptor().Type]))
	} else {
		valName := fmt.Sprintf("rpfgVal%d", rf.Number())
		p.Printf(`if (! %s) {
  %s = []
}
let %s;
%s
%s.push(%s);
`, val, val, valName, rf.Element.DeserializeBlock(valName),val, valName)
	}
}

func (rf *RepeatedField) SerializeFunction() string {
	panic("RepeatedField:SerializeFunction is not supposed to be called")
}

func (rf *RepeatedField) DeserializeBlock(valName string) string {
	panic("RepeatedField:DeserializeBlock is not supposed to be called")
}