package generator

import (
	"fmt"
	pgs "github.com/lyft/protoc-gen-star"
)

type MapMessageField struct {
	Field
	Key FieldGenerator
	Element FieldGenerator
	IsObjectMap bool
}

func NewMapMessageField(pgsField pgs.Field, msg *MessageGenerator, o *Options) (FieldGenerator, error) {
	if ! pgsField.Type().IsMap() {
		return nil, nil
	}
	m := MapMessageField{Field: newField(pgsField, msg, o)}
	var err error
	typeName := pgsField.Descriptor().GetTypeName()
	for _, mEntry := range msg.Msg.MapEntries(){
		if mEntry.FullyQualifiedName() == typeName {
			m.Key, m.Element, err = GetMapKeyEntry(&MessageGenerator{
				File:  msg.File,
				Msg:    mEntry,
			}, o)
			break
		}
	}
	if err != nil {
		return nil, err
	}
	if m.Key == nil {
		return nil, fmt.Errorf("cannot find type description for field %s", pgsField.Name().String())
	}
	if m.Key.ES6Type() == "string" {
		m.es6Type = fmt.Sprintf("{ [string]: %s }", m.Element.ES6Type())
		m.typeValidationFunc = func(val string) string { return "" }
		m.checkEmptyFunc = func(val string) string { return fmt.Sprintf("%s && Object.entries(%s).length > 0", val, val) }
		m.IsObjectMap = true
	} else {
		m.es6Type = fmt.Sprintf("Map<%s, %s>", m.Key.ES6Type(), m.Element.ES6Type())
		m.typeValidationFunc = func(val string) string { return fmt.Sprintf("%s== null || %s  instanceof Map", val, val) }
		m.checkEmptyFunc = func(val string) string { return fmt.Sprintf("%s && %s.size > 0", val, val) }
	}
	return &m, nil
}

func GetMapKeyEntry(msg *MessageGenerator, o *Options) (FieldGenerator, FieldGenerator, error) {
	resolver := CompositeFieldResolver(NewEnumFieldGenerator, NewMessageFieldGenerator, NewPrimitiveFieldWriter)
	var key, element FieldGenerator
	var err error
	for _, field := range msg.Msg.Fields() {
		if field.Descriptor().GetNumber() == 1 {
			key, err = resolver(field, msg, o)
		} else if field.Descriptor().GetNumber() == 2{
			element, err = resolver(field, msg, o)
		}
		if err != nil {
			return nil, nil, err
		}
	}
	if key == nil || element == nil {
		return nil, nil, fmt.Errorf("cannot find key or element for embedded map type %s", msg.Msg.FullyQualifiedName())
	}
	return key, element, nil
}

func (m *MapMessageField) FromObjectExp(src string) string {
	if ! m.Element.IsMessage() {
		return m.Field.FromObjectExp(src)
	}
	if m.IsObjectMap {
		return fmt.Sprintf("%s && Object.fromEntries(Object.keys(%s).map(key => [key, %s]))", src, src, m.Element.FromObjectExp(src + "[key]"))
	} else {
		return fmt.Sprintf("%s && new Map(%s.keys().map(key => [key, %s]))", src, src, m.Element.FromObjectExp(src+".get(key)"))
	}
}

func (m *MapMessageField) ToObjectExp(src string) string {
	if ! m.Element.IsMessage() {
		return m.Field.FromObjectExp(src)
	}
	if m.IsObjectMap {
		return fmt.Sprintf("%s && Object.fromEntries(Object.keys(%s).map(key => [key, %s]))", src, src, m.Element.ToObjectExp(src + "[key]"))
	} else {
		return fmt.Sprintf("%s && new Map(%s.keys().map(key => [key, %s]))", src, src, m.Element.ToObjectExp(src+".get(key)"))
	}
}


func (m *MapMessageField) SerializeFunction() string {
	panic("RepeatedMessageFieldGenerator:SerializeFunction is not supposed to be called")
}


func (m *MapMessageField) GenerateSerializeBlock(p Printer, val string) {
	p.PrintTpl("map_serialize", `if ({{.emptyCheck}}) {
  {{ if .isObjectMap -}}
  const arrKeys = Object.keys({{.val}}).sort();
  {{- else -}}
  const arrKeys = Array.from({{.val}}.keys()).sort();
  {{ end }}
  for (const key of arrKeys) {
    {{ if .isObjectMap -}}
    const value = {{.val}}[key];
    {{- else -}}
    const value = {{.val}}.get(key);
    {{ end -}} 
    writer.beginSubMessage({{.fieldNumber}});
    writer.{{.key.SerializeFunction}}(1, key);
    {{- if .el.IsMessage }}
    writer.writeMessage(2, value, {{.el.SerializeFunction}});
    {{- else }}
    writer.{{.el.SerializeFunction}}(2, value);
    {{- end }}
    writer.endSubMessage();
  }
}
`, "emptyCheck", m.CheckNotEmptyExp(val), "val", val, "fieldNumber", m.Number(), "key", m.Key, "el", m.Element, "isObjectMap", m.IsObjectMap)
}

func (m *MapMessageField) GenerateDeserializeBlock(p Printer, val string) {
	p.PrintTpl("map_deserialize", `if (! {{.val}} ){
  {{ if .isObjectMap -}}
  {{.val}} = {}
  {{- else -}}
  {{.val}} = new Map();
  {{- end }}
}
reader.readMessage({{.val}}, (map, reader) => {
  let key {{- if .key.DefaultValue -}}={{.key.DefaultValue}}{{- end -}};
  let value {{- if .el.DefaultValue -}}={{.el.DefaultValue}}{{- end -}};
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    const fieldNumber = reader.getFieldNumber();
    if (fieldNumber === 1) { // Key.
{{ indent (.key.DeserializeBlock "key") 6 }}  
    } else if (fieldNumber === 2) { // Value.
{{ indent (.el.DeserializeBlock "value") 6}}  
    }
  }
  {{ if .isObjectMap -}}
  {{.val}}[key] = value;
  {{- else -}}
  map.set(key, value);
  {{- end }}
})
`, "val", val, "key", m.Key, "el", m.Element, "isObjectMap", m.IsObjectMap)
}

func (m *MapMessageField) DeserializeBlock(valName string) string {
	panic("MapMessageField:DeserializeBlock is not supposed to be called")
}