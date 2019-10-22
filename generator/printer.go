package generator

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"
)

type Printer interface {
	io.Writer
	Indent() Printer
	IndentBy(indent int) Printer
	Print(s string)
	Printf(format string, params... interface{})
	PrintTpl(tmplName string, tmpl string, values... interface{})
}

type printer struct {
	indent string
	indentIncrease *string
	atTheBeginning *bool
	w io.Writer
	templates map[string]*template.Template
}

func NewPrinter(w io.Writer, indentIncrease int) Printer {
	atb := true
	ii := strings.Repeat(" ", indentIncrease)
	return &printer{
		indent: "",
		indentIncrease: &ii,
		atTheBeginning: &atb,
		w: w,
		templates: make(map[string]*template.Template),
	}
}

func (p *printer) Indent() Printer {
	return &printer{
		indent:         p.indent + *p.indentIncrease,
		indentIncrease: p.indentIncrease,
		atTheBeginning: p.atTheBeginning,
		w:              p.w,
		templates: 		p.templates,
	}
}

func (p *printer) IndentBy(indent int) Printer {
	return &printer{
		indent:         p.indent + strings.Repeat(" ", indent),
		indentIncrease: p.indentIncrease,
		atTheBeginning: p.atTheBeginning,
		w:              p.w,
		templates: 		p.templates,
	}
}

func (p *printer) Write(b []byte) (int,  error) {
	p.Print(string(b))
	return len(b), nil
}

func (p *printer) Print(s string)  {
	if *p.atTheBeginning && !strings.HasPrefix(s, "\n") {
		s = p.indent + s
	}
	lineEnds := strings.HasSuffix(s, "\n")
	s = strings.ReplaceAll(s, "\n", "\n" + p.indent)
	if lineEnds {
		// remove indent at the end
		s = strings.TrimSuffix(s, p.indent)
	}
	_, err := p.w.Write([]byte(s))
	if err != nil {
		panic(err)
	}
	*p.atTheBeginning = lineEnds
}

func (p *printer) Printf(format string, params... interface{}) {
	s := fmt.Sprintf(format, params...)
	p.Print(s)
}

func (p *printer) PrintTpl(tmplName string, tmpl string, values... interface{}) {
	if len(values) % 2 != 0 {
		panic("Invalid number of parameters for PrintTpl")
	}
	data := make(map[string]interface{})
	for i:=0; i < len(values); i += 2 {
		data[values[i].(string)] = values[i + 1]
	}
	if _, ok := p.templates[tmplName]; !ok {
		t := template.New(tmplName)

		p.templates[tmplName] = template.Must(t.Funcs(funcs).Parse(tmpl))
	}
	b := bytes.Buffer{}
	err := p.templates[tmplName].Execute(&b, data)
	if err != nil {
		panic(err)
	}
	p.Printf(b.String())
}
var funcs = template.FuncMap{"indent": indent}
func indent(s string, indent int) string {
	ind := strings.Repeat(" ", indent)
	lineEnds := strings.HasSuffix(s, "\n")
	if !strings.HasPrefix(s, "\n") {
		s = ind + s
	}
	s = strings.ReplaceAll(s, "\n", "\n" + ind)
	if lineEnds {
		// remove indent at the end
		s = strings.TrimSuffix(s, ind)
	}
	return s
}
