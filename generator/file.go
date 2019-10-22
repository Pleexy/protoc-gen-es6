package generator

import (
	pgs "github.com/lyft/protoc-gen-star"
	"github.com/pkg/errors"
	"path/filepath"
	"strings"
)

type Import struct {
	FilePath string
	ImportPath string
	TypePrefix string
	Types map[string]struct{}
}

type FileGenerator struct {
	PgsFile pgs.File
	Messages []*MessageGenerator
	Enums []*EnumGenerator
	Imports map[string]*Import
	Opt *Options
}

func NewFileGenerator(pgsFile pgs.File, o *Options, resolver FieldResolver) (*FileGenerator, error) {
	f := &FileGenerator{
		PgsFile:pgsFile,
		Imports: make(map[string]*Import),
		Messages: make([]*MessageGenerator, len(pgsFile.Messages())),
		Enums: make([]*EnumGenerator, len(pgsFile.Enums())),
		Opt: o,
	}
	for i, msg := range pgsFile.Messages() {
		msgGen, err := NewMessageGenerator(msg, f, o, resolver)
		if err != nil {
			return nil, err
		}
		f.Messages[i] = msgGen
	}
	for i, enum := range pgsFile.Enums() {
		enumGen, err := NewEnumGenerator(enum)
		if err != nil {
			return nil, err
		}
		f.Enums[i] = enumGen
	}
	return f, nil
}

func (f *FileGenerator) Generate(pr Printer) {
	if f.Opt.Flow{
		pr.Print("// @flow\n")
	}
	f.generateImports(pr)
	pr.Print("\n\n")
	f.generateMessages(pr)
	f.generateEnums(pr)
}

func (f *FileGenerator) RegisterImport(typeName string, depFile pgs.File) (string, error) {
	depPath := depFile.InputPath().String()
	importPath, err := f.calculateDepPath(depFile)
	if err != nil {
		return "", errors.Wrapf(err, "cannot resolve import path for %s", depPath)
	}
	if _, ok := f.Imports[depPath]; !ok {
		imp := Import{
			FilePath: depPath  ,
			ImportPath: importPath,
			TypePrefix: f.depToPrefix(depFile),
			Types: make(map[string]struct{}),
		}
		f.Imports[depPath] = &imp
	}
	// add type if not added already
	if _, ok := f.Imports[depPath].Types[typeName]; !ok {
		f.Imports[depPath].Types[typeName] = struct{}{}
	}
	return f.Imports[depPath].TypePrefix, nil
}

func (f *FileGenerator) generateImports(pr Printer) {
	if f.Opt.ESModules {
		pr.Print("import * as jspb from 'google-protobuf';\n")
	} else {
		pr.Print("const jspb = require('google-protobuf');\n")
	}
	for _, imp := range f.Imports {
		if f.Opt.ESModules {
			pr.Printf("import * as %s from '%s';\n", imp.TypePrefix, imp.ImportPath)
		} else {
			pr.Printf("const %s = require('%s');\n", imp.TypePrefix, imp.ImportPath)
		}
	}
}

func (f *FileGenerator) generateMessages(pr Printer) {
	for _, msg := range f.Messages {
		msg.Generate(pr)
	}
}

func (f *FileGenerator) generateEnums(pr Printer) {
	for _, enum := range f.Enums {
		enum.Generate(pr)
	}
}

func (f *FileGenerator) calculateDepPath(dep pgs.File) (string, error) {
	if strings.Contains(dep.InputPath().String(), "google/protobuf") {
		return "google-protobuf/"+dep.InputPath().SetExt("_pb.js").String(), nil
	} else {
		rel, err := filepath.Rel(f.PgsFile.InputPath().Dir().String(), dep.InputPath().Dir().String())
		if err != nil {
			return "", err
		}
		ext := ".pb"
		if f.Opt.ESModules {
			ext = ".pb.mjs"
		}
		depPath := filepath.Join(rel, dep.InputPath().SetExt(ext).Base())
		if ! strings.Contains(depPath, "/") {
			depPath = "./"+depPath
		}
		return depPath, nil
	}
}

func (f *FileGenerator) depToPrefix(dep pgs.File) string {
	return strings.ReplaceAll(dep.InputPath().SetExt("").String(),"/", "_")
}