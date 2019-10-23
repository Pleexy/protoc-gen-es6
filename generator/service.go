package generator

import (
	pgs "github.com/lyft/protoc-gen-star"
	"github.com/pkg/errors"
	"strings"
)

type ServiceGenerator struct {
	PgsService pgs.Service
	File *FileGenerator
	ClassName string
	Opt *Options
	Methods []*MethodGenerator
}

func NewServiceGenerator(pgsService pgs.Service, f *FileGenerator, opt *Options) (*ServiceGenerator, error) {
	s := &ServiceGenerator{
		PgsService:pgsService,
		ClassName: pgsService.Name().UpperCamelCase().String()+"Service",
		Opt: opt,
		Methods: make([]*MethodGenerator, len(pgsService.Methods())),
		File: f,
	}
	var err error
	for i, pgsM := range pgsService.Methods() {
		s.Methods[i], err = NewMethodGenerator(pgsM, s, opt)
		if err != nil {
			return nil, errors.Wrapf(err, "error while processing method %s", pgsM.Name().String())
		}
	}
	return s, nil
}

func (s *ServiceGenerator) Generate (p Printer) {
	s.GenerateHeader(p)
	for _, method := range s.Methods {
		method.Generate(p)
	}
	s.GenerateFooter(p)
}

func (s *ServiceGenerator) GenerateHeader (p Printer) {
	if s.Opt.ESModules {
		p.Printf("export const %s = {\n", s.ClassName)
	} else {
		p.Printf("const %s = {\n", s.ClassName)
	}
}

func (s *ServiceGenerator) GenerateFooter (p Printer) {
	if s.Opt.ESModules {
		p.Print("}")
	} else {
		p.Printf("}\nmodule.exports.%s = %s", s.ClassName, s.ClassName)
	}
}


type MethodGenerator struct {
	PgsMethod pgs.Method
	Path string
	RequestType string
	ResponseType string
	Opt *Options
}

func NewMethodGenerator(pgsMethod pgs.Method, service *ServiceGenerator, opt *Options) (*MethodGenerator, error) {
	m := &MethodGenerator{
		PgsMethod:pgsMethod,
		Opt: opt,
	}
	m.Path = "/" + strings.TrimPrefix(service.PgsService.FullyQualifiedName(),".") + "/" + pgsMethod.Name().String()
	prefix, typeName, err := GetTypeNameAndPrefix(pgsMethod.Input().File(), pgsMethod, pgsMethod.Input().FullyQualifiedName(), service.File)
	if err != nil {
		return nil, err
	}
	m.RequestType = prefix + typeName
	prefix, typeName, err = GetTypeNameAndPrefix(pgsMethod.Output().File(), pgsMethod, pgsMethod.Output().FullyQualifiedName(), service.File)
	if err != nil {
		return nil, err
	}
	m.ResponseType = prefix + typeName
	return m, nil
}

func (m *MethodGenerator) Generate(p Printer) {
	p.PrintTpl("service_method",`{{.m.PgsMethod.Name.LowerCamelCase.String}}: {
  path: '{{.m.Path}}',
  requestStream: {{.m.PgsMethod.ClientStreaming}},
  responseStream: {{.m.PgsMethod.ServerStreaming}},
  requestType: {{.m.RequestType}},
  responseType: {{.m.ResponseType}},
  requestSerialize: (arg{{- if .flow -}}:{{.m.RequestType}}):Buffer{{- else -}}){{- end }} => {
    if (!(arg instanceof {{.m.RequestType}})) {
      throw new Error('Expected argument of type {{.m.RequestType}}');
    }
    return Buffer.from(arg.serializeBinary());
  },
  requestDeserialize: (buffer_arg{{- if .flow -}}:Buffer):{{.m.RequestType}}{{- else -}}){{- end }} =>
    {{.m.RequestType}}.deserializeBinary(new Uint8Array(buffer_arg)),
  responseSerialize: (arg{{- if .flow -}}:{{.m.ResponseType}}):Buffer{{- else -}}){{- end }} => {
    if (!(arg instanceof {{.m.ResponseType}})) {
      throw new Error('Expected argument of type {{.m.ResponseType}}');
    }
    return Buffer.from(arg.serializeBinary());
  },
  responseDeserialize: (buffer_arg{{- if .flow -}}:Buffer):{{.m.ResponseType}}{{- else -}}){{- end }} =>
    {{.m.ResponseType}}.deserializeBinary(new Uint8Array(buffer_arg)),
},
`,"m", m, "flow", m.Opt.Flow)
}
