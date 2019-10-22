package generator

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	pgs "github.com/lyft/protoc-gen-star"
)

type FieldGenerator interface {
	GenerateProperty(p Printer)
	GenerateObjectField (p Printer)
	GenerateGetter(p Printer)
	GenerateSetter(p Printer)
	GenerateSerializeBlock(p Printer, val string)
	GenerateDeserializeBlock(p Printer, val string)
	TypeValidationExp(name string) string
	CheckNotEmptyExp(name string) string
	PropertyName() string
	Number() int32
	GetSetName() string
	ProtoName() string
	DefaultValue() string
	ES6Type() string
	ProtoType() string
	SerializeFunction() string
	IsMessage() bool
	IsPacked() bool
	FromObjectExp(src string) string
	ToObjectExp(src string) string
	DeserializeBlock(valName string) string
}

type Field struct {
	Message *MessageGenerator
	pgsField pgs.Field
	es6Type string
	defaultValue string
	typeValidationFunc func(name string) string
	checkEmptyFunc func(name string) string
	protoType string
	o *Options
}

func newField(pgsField pgs.Field, msg *MessageGenerator, o *Options) Field {
	return Field{
		Message: msg,
		pgsField: pgsField,
		o: o,
	}
}

func (f *Field) IsPacked() bool  {
	if f.IsMessage() {
		return false
	}
	return IsPacked(f.pgsField.Descriptor())
}

func (f *Field) IsMessage() bool {
	return *f.pgsField.Descriptor().Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE
}

func (f *Field) PropertyName() string {
	if f.o.UsePrivateFields {
		return "#" + *f.pgsField.Descriptor().JsonName
	} else {
		if f.o.GenerateGettersSetters {
			return "_" + *f.pgsField.Descriptor().JsonName	// rename fields
		} else {
			return *f.pgsField.Descriptor().JsonName
		}
	}
}

func (f* Field) ES6Type() string {
	return f.es6Type
}

func (f* Field) ProtoType() string {
	return f.protoType
}


func (f *Field) GetSetName() string {
	return *f.pgsField.Descriptor().JsonName
}

func (f *Field) Number() int32 {
	return *f.pgsField.Descriptor().Number
}

func (f *Field) ProtoName() string {
	return *f.pgsField.Descriptor().Name
}

func (f *Field) DefaultValue() string {
	return f.defaultValue
}

func (f *Field) TypeValidationExp(name string) string {
	if f.typeValidationFunc != nil {
		return f.typeValidationFunc(name)
	} else {
		return ""
	}
}

func (f *Field) CheckNotEmptyExp(name string) string {
	if f.checkEmptyFunc != nil {
		return f.checkEmptyFunc(name)
	} else {
		panic("Check empty is not defined for the field")
	}
}

func (f *Field) GenerateProperty(p Printer) {
	defaultValue := ""
	if f.defaultValue != "" {
		defaultValue = " = " + f.DefaultValue()
	}
	if f.o.Flow {
		p.Printf("%s: %s%s;  // %s %s = %d;\n", f.PropertyName(), f.es6Type, defaultValue, f.protoType, f.ProtoName(), f.Number())
	} else {
		p.Printf("%s%s;  // %s %s = %d;\n", f.PropertyName(), defaultValue, f.protoType, f.ProtoName(), f.Number())
	}
}
func (f *Field) GenerateObjectField(p Printer) {
	if f.o.Flow {
		p.Printf("%s: %s,  // %s %s = %d;\n", f.GetSetName(), f.es6Type, f.protoType, f.ProtoName(), f.Number())
	} else {
		p.Printf("%s,  // %s %s = %d;\n", f.GetSetName(),  f.protoType, f.ProtoName(), f.Number())
	}
}

func (f *Field) GenerateFromObject(p Printer, destField, srcField string) {
	p.Printf("%s = %s;\n", destField, srcField)
}

func (f *Field) FromObjectExp(src string) string {
	return src
}

func (f *Field) ToObjectExp(src string) string {
	return src
}

func (f *Field) GenerateGetter(p Printer) {
	p.Printf(`/**
* optional %s %s = %d;
* @return {%s}
*/
`, f.protoType, f.ProtoName(), f.Number(), f.es6Type)
	if f.o.Flow {
		p.Printf("get %s():%s {\n", f.GetSetName(), f.es6Type)
	} else {
		p.Printf("get %s(){\n", f.GetSetName())
	}
	p.Printf("  return this.%s;\n}\n\n", f.PropertyName())
}

func (f *Field) GenerateSetter(p Printer) {
	p.Printf("/** @param {%s} val */\n", f.es6Type)
	if f.o.Flow {
		p.Printf("set %s(val: %s):void {\n", f.GetSetName(), f.es6Type)
	} else {
		p.Printf("set %s(val) {\n", f.GetSetName())
	}
	typeValidation := f.TypeValidationExp("val")
	if f.o.ValidateOnSet && typeValidation != "" {
		p.Printf(`  if (%s) {
    this.%s = val;
  } else {
    throw new Error('Expected type %s for field %s');
  }
}
`, typeValidation, f.PropertyName(), f.es6Type, f.PropertyName())
	} else {
		p.Printf("  this.%s = val;\n}\n\n\n", f.PropertyName())
	}
}

type FieldResolver func (pgsField pgs.Field, msg *MessageGenerator, o *Options) (FieldGenerator, error)

func CompositeFieldResolver(resolvers... FieldResolver) FieldResolver {
	return func(pgsField pgs.Field, msg *MessageGenerator, o *Options) (FieldGenerator, error) {
		for _, resolver := range resolvers {
			fr, err := resolver(pgsField, msg, o)
			if err != nil {
				return nil, err
			}
			if fr != nil {
				return fr, nil
			}
		}
		return nil, nil
	}
}