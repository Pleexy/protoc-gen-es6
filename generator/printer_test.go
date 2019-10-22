package generator_test

import (
	"bytes"
	"github.com/Pleexy/protoc-gen-es6/generator"
	"log"
	"testing"
)

// todo: implement tests

func TestPrinter_Print(t *testing.T) {

}

func TestPrinter_Indent(t *testing.T) {
	b := bytes.Buffer{}
	p := generator.NewPrinter(&b,  2)
	p.Print("No indent ")
	p.Indent().Print("same line no indent\nindented\n")
	if b.String() != "No indent same line no indent\n  indented\n" {
		log.Print(b.String())
		t.Fail()
	}
}

func TestPrinter_NewLine(t *testing.T) {
	b := bytes.Buffer{}
	p := generator.NewPrinter(&b,  2)
	p.Print("No indent\n")
	p.Indent().Print("indented\n")
	if b.String() != "No indent\n  indented\n" {
		log.Print(b.String())
		t.Fail()
	}
}

func TestPrinter_PrintTpl(t *testing.T) {
	b := bytes.Buffer{}
	p := generator.NewPrinter(&b,  2)
	tmpl := "{{- if .key -}}{{.str}} {{ .int }}{{- else -}}{{.int2}}{{- end -}}"
	p.PrintTpl("testtmpl", tmpl, "key", true, "str","svalue", "int", 1, "int2", 2)
	if b.String() != "svalue 1" {
		log.Print(b.String())
		t.Fail()
	}
	b = bytes.Buffer{}
	p = generator.NewPrinter(&b,  2)
	p.PrintTpl("testtmpl", tmpl, "key", false, "str","svalue", "int", 1, "int2", 2)
	if b.String() != "2" {
		log.Print(b.String())
		t.Fail()
	}
}

