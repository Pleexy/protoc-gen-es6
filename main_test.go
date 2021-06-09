package main_test

import (
	"github.com/Pleexy/protoc-gen-es6/generator"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"os"
	"testing"
)

func TestModule(t *testing.T) {
	req, err := os.Open("./code_generator_request.pb.bin")
	if err != nil {
		t.Fatal(err)
	}

//	fs := afero.NewMemMapFs()
//	res := &bytes.Buffer{}

	pgs.Init(
		pgs.ProtocInput(req),  // use the pre-generated request
		pgs.DebugEnv("DEBUG"),
	).RegisterModule(
		generator.ES6(),
	).RegisterPostProcessor(
		pgsgo.GoFmt(),
	).Render()

}