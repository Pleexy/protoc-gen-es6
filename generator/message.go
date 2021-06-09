package generator

import (
	"fmt"
	pgs "github.com/lyft/protoc-gen-star"
	"strings"
)

type MessageGenerator struct {
	File *FileGenerator
	Msg pgs.Message
	Fields []FieldGenerator
	Opt *Options
	Messages []*MessageGenerator
	Enums []*EnumGenerator
}

func NewMessageGenerator(msg pgs.Message, f *FileGenerator,  opt *Options, resolver FieldResolver) (*MessageGenerator, error) {
	m := &MessageGenerator{
		File: f,
		Msg: msg,
		Opt: opt,
		Fields:  make([]FieldGenerator,len(msg.Fields())),
		Messages: make([]*MessageGenerator, len(msg.Messages())),
		Enums: make([]*EnumGenerator, len(msg.Enums())),

	}
	var err error
	for i, fld := range msg.Fields() {
		m.Fields[i], err = resolver(fld, m, m.Opt)
		if err != nil {
			return nil, err
		}
		if m.Fields[i] == nil {
			return nil, fmt.Errorf("cannot find generator for field %s", fld.Name().String())
		}
	}
	for i, msg := range msg.Messages() {
		msgGen, err := NewMessageGenerator(msg, f, opt, resolver)
		if err != nil {
			return nil, err
		}
		m.Messages[i] = msgGen
	}
	for i, enum := range msg.Enums() {
		enumGen, err := NewEnumGenerator(enum)
		if err != nil {
			return nil, err
		}
		m.Enums[i] = enumGen
	}

	return m, nil
}

func (m *MessageGenerator) Generate(pr Printer )  {
	if m.Opt.Flow {
		m.GenerateObjectType(pr)
	}
	m.GenerateHeader(pr)
	m.GenerateProperties(pr.Indent())
	m.GenerateFromObject(pr.Indent())
	pr.Print("\n")
	m.GenerateToObject(pr.Indent())
	pr.Print("\n")
	m.GenerateDeserialize(pr.Indent())
	m.GenerateSerialize(pr.Indent())
	m.GenerateFooter(pr)
	for _, enum := range m.Enums {
		enum.Generate(pr)
	}
	for _, msg := range m.Messages {
		msg.Generate(pr)
	}
}


func (m *MessageGenerator) ClassName() string {

	return m.Msg.Name().UpperCamelCase().String()
}

func (m *MessageGenerator) GenerateFromObject(pr Printer) {
	if m.Opt.Flow {
		pr.Printf("static fromObject(obj:%s$Object):%s {\n", m.ClassName(), m.ClassName())
	} else {
		pr.Print("static fromObject (obj){\n")
	}
	pr.Printf( "  const newObj = new %s();\n", m.ClassName())
	pr.Print( "  if (obj) {\n")
	for _, fieldGen := range m.Fields {
		pr.Printf("    newObj.%s = %s;\n", fieldGen.GetSetName(), fieldGen.FromObjectExp("obj." + fieldGen.GetSetName()))
	}
	pr.Print( "  }\n")
	pr.Print("  return newObj;\n}\n")
}

func (m *MessageGenerator) GenerateToObject(pr Printer) {
	if m.Opt.Flow {
		pr.Printf("toObject():%s$Object {\n", m.ClassName())
	} else {
		pr.Print("toObject(){\n")
	}
	pr.Print( "  const newObj = {};\n")
	for _, fieldGen := range m.Fields {
		pr.Printf("  newObj.%s = %s;\n", fieldGen.GetSetName(), fieldGen.ToObjectExp("this." + fieldGen.GetSetName()))
	}
	pr.Print("  return newObj;\n}\n")
}


func (m *MessageGenerator) GenerateHeader(pr Printer)  {
	name := m.ClassName()
	embed := strings.Contains(name, ".")
	if embed {
		pr.Printf("%s = class {\n", name)
	} else {
		if m.Opt.ESModules {
			pr.Printf("export class %s {\n", name)
		} else {
			pr.Printf("class %s {\n", name)
		}
	}
}
func (m *MessageGenerator) GenerateFooter(pr Printer)  {
	name := m.ClassName()
	embed := strings.Contains(name, ".")
	pr.Print("}\n")
	if ! embed && ! m.Opt.ESModules {
		pr.Printf("module.exports.%s = %s;\n", name, name)
	}
}

func (m *MessageGenerator) GenerateProperties(pr Printer)  {
	for _, fieldGen := range m.Fields {
		fieldGen.GenerateProperty(pr)
		pr.Print("\n")
	}
	for _, fieldGen := range m.Fields {
		fieldGen.GenerateGetter(pr)
		pr.Print("\n")
		fieldGen.GenerateSetter(pr)
		pr.Print("\n")
	}
}

func (m *MessageGenerator) GenerateObjectType(pr Printer) {
	pr.Printf("export type %s$Object = {\n", m.ClassName())
	prIndent := pr.Indent()
	for _, fieldGen := range m.Fields {
		fieldGen.GenerateObjectField(prIndent)
	}
	pr.Print("}\n\n")
}

func (m *MessageGenerator) GenerateDeserialize(pr Printer) {
	dbComment := fmt.Sprintf(`
/**
* Deserializes binary data (in protobuf wire format).
* @param {Uint8Array} bytes The bytes to deserialize.
* @return {!%s}
*/`, m.ClassName())
	dbfrComment := fmt.Sprintf(`
/**
* Deserializes binary data (in protobuf wire format) from the
* given reader into the given message object.
* @param {!%s} msg The message object to deserialize into.
* @param {!jspb.BinaryReader} reader The BinaryReader to use.
* @return {!%s}
*/`, m.ClassName(), m.ClassName())

	if m.Opt.Flow {
		pr.Printf(`%s
static deserializeBinary(bytes: Uint8Array):%s {
  const reader = new jspb.BinaryReader(bytes);
  const msg = new %s();
  return %s.deserializeBinaryFromReader(msg, reader);
}
%s
static deserializeBinaryFromReader (msg: %s, reader: BinaryReader) {
`, dbComment,  m.ClassName(),  m.ClassName(),  m.ClassName(), dbfrComment,  m.ClassName())
	} else {
		pr.Printf(`%s
static deserializeBinary(bytes) {
  const reader = new jspb.BinaryReader(bytes);
  const msg = new %s();
  return %s.deserializeBinaryFromReader(msg, reader);
}
%s
  static deserializeBinaryFromReader (msg, reader) {
`, dbComment, m.ClassName(), m.ClassName(), dbfrComment)
	}
	pr.Print(`  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    let field = reader.getFieldNumber();
    switch (field) {
`)
	for _, fld := range m.Fields {
		pr.Printf("      case %d:\n", fld.Number())
		fld.GenerateDeserializeBlock(pr.IndentBy(8), "msg." + fld.PropertyName())
		pr.Print("        break;\n")
	}
	pr.Print(`      default:
        reader.skipField();
        break;
      }
    }
    return msg;
  };
`)
}
func (m *MessageGenerator) GenerateSerialize(pr Printer) {
	sbComment := `
/**
* Serializes the message to binary data (in protobuf wire format).
* @return {!Uint8Array}
*/`
	sbtwComment := fmt.Sprintf(`
/**
* Serializes the given message to binary data (in protobuf wire
* format), writing to the given BinaryWriter.
* @param {!%s} msg
* @param {!jspb.BinaryWriter} writer
* @suppress {unusedLocalVariables} f is only used for nested messages
*/`, m.ClassName())
	if m.Opt.Flow {
		pr.Printf(`%s
serializeBinary(): Uint8Array {
  const writer = new jspb.BinaryWriter();
  %s.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
}
%s
static serializeBinaryToWriter(msg: %s, writer: BinaryWriter): Uint8Array {
`, sbComment, m.ClassName(), sbtwComment, m.ClassName() )
	}  else {
		pr.Printf(`%s
serializeBinary() {
  const writer = new jspb.BinaryWriter();
  %s.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
}
%s
static serializeBinaryToWriter(msg, writer) {
`, sbComment, m.ClassName(), sbtwComment )
	}
	for _, fld := range m.Fields {
		fld.GenerateSerializeBlock(pr.Indent(), "msg." + fld.PropertyName())
	}
	pr.Print("}\n")
}


