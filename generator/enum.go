package generator

import (
	pgs "github.com/lyft/protoc-gen-star"
	"strings"
)

type EnumGenerator struct {
	PgsEnum pgs.Enum
	Name string
}

func NewEnumGenerator(pgsEnum pgs.Enum) (*EnumGenerator, error) {
	en := &EnumGenerator{PgsEnum:pgsEnum}
	if name, err := GetMessageName(en.PgsEnum.FullyQualifiedName(), en.PgsEnum.Package()); err != nil {
		return nil, err
	} else {
		en.Name = name
	}
	return en, nil
}

func (en *EnumGenerator) Generate(p Printer) {

	p.PrintTpl("enum",`
/** {{.name }} */
{{ if not .embed -}}const {{ end -}}{{.name}} = Object.freeze({
{{- range .values }}
  {{.Descriptor.GetName}}: {{.Value}},
{{- end }}
}){{- if not .embed }};
module.exports.{{.name}}={{.name}};
{{ end -}}
`, "embed",strings.Contains(en.Name, "."),"name", en.Name , "values", en.PgsEnum.Values())
}