package generator

import (
	"bytes"
	pgs "github.com/lyft/protoc-gen-star"
	//pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

type fieldWriterFactory func(fld pgs.Field) (FieldGenerator, error)


// JSONifyPlugin adds encoding/json Marshaler and Unmarshaler methods on PB
// messages that utilizes the more correct jsonpb package.
// See: https://godoc.org/github.com/golang/protobuf/jsonpb
type ES6Module struct {
	*pgs.ModuleBase
//	ctx pgsgo.Context
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
	}
	return &ES6Module{
		ModuleBase: &pgs.ModuleBase{},
		FieldResolver: CompositeFieldResolver( NewMapMessageField, NewRepeatedFieldWriter, NewEnumFieldGenerator,NewMessageFieldGenerator, NewPrimitiveFieldWriter),
		o: o,
	}
}

func (p *ES6Module) InitContext(c pgs.BuildContext) {
	p.ModuleBase.InitContext(c)
	//p.ctx = pgsgo.InitContext(c.Parameters())
}

// Name satisfies the generator.Plugin interface.
func (p *ES6Module) Name() string { return "jsonify" }

func (p *ES6Module) Execute(targets map[string]pgs.File, pkgs map[string]pgs.Package) []pgs.Artifact {

	for _, t := range targets {
		p.generate(t)
	}
	return p.Artifacts()
}


func (p *ES6Module) generate(f pgs.File) {
	ext := ".pb.es6"
	if p.o.ESModules {
		ext = ".pb.mjs"
	}
	name :=f.InputPath().SetExt(ext).String()
	//name := p.ctx.OutputPath(f).SetExt(".es6")
	buf := &bytes.Buffer{}
	pr := NewPrinter(buf, 2)
	fg, err := NewFileGenerator(f, p.o, p.FieldResolver)
	p.CheckErr(err)
	fg.Generate(pr)
	p.AddGeneratorFile(name, buf.String())

}



