package generator

import (
	"bytes"
	pgs "github.com/lyft/protoc-gen-star"
)

type fieldWriterFactory func(fld pgs.Field) (FieldGenerator, error)


// JSONifyPlugin adds encoding/json Marshaler and Unmarshaler methods on PB
// messages that utilizes the more correct jsonpb package.
// See: https://godoc.org/github.com/golang/protobuf/jsonpb
type ES6Module struct {
	*pgs.ModuleBase
	flow bool
	FieldResolver
	o *Options
}

type message struct {
	name string
}

// JSONify returns an initialized JSONifyPlugin
func ES6() *ES6Module {
	o := &Options{
		UsePrivateFields:       false,
		GenerateGettersSetters: true,
		ValidateOnSet:          true,
		Flow:                   true,
		ESModules:				false,
		ConvertString:			true,
		Grpc:					true,
		ReplaceJSOut: 			false,

	}
	return &ES6Module{
		ModuleBase: &pgs.ModuleBase{},
		FieldResolver: CompositeFieldResolver( NewMapMessageField, NewRepeatedFieldWriter, NewEnumFieldGenerator,NewMessageFieldGenerator, NewPrimitiveFieldWriter),
		o: o,
	}
}

func (p *ES6Module) InitContext(c pgs.BuildContext) {
	p.ModuleBase.InitContext(c)
	params := c.Parameters()
	if _, ok := params["noflow"]; ok {
		p.o.Flow = false
	}
	if _, ok := params["jsout"]; ok {
		p.o.ReplaceJSOut = true
	}
}

// Name satisfies the generator.Plugin interface.
func (p *ES6Module) Name() string { return "es6" }

func (p *ES6Module) Execute(targets map[string]pgs.File, pkgs map[string]pgs.Package) []pgs.Artifact {

	for _, t := range targets {
		p.generate(t)
	}
	return p.Artifacts()
}

func (p *ES6Module) generate(f pgs.File) {
	name := f.InputPath()
	if p.o.ReplaceJSOut {
		name = name.SetBase(name.BaseName()+"_pb").SetExt(".js")
	} else if p.o.ESModules {
		name = name.SetExt(".pb.mjs")
	} else {
		name = name.SetExt(".pb.es6")
	}
	buf := &bytes.Buffer{}
	pr := NewPrinter(buf, 2)
	fg, err := NewFileGenerator(f, p.o, p.FieldResolver, name)
	p.CheckErr(err)
	fg.Generate(pr)
	p.AddGeneratorFile(name.String(), buf.String())
}



