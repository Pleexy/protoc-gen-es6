package generator_test

import (
	"bytes"
	pgs "github.com/lyft/protoc-gen-star"
	"github.com/spf13/afero"
	"log"
	"os"
	"testing"
	"github.com/Pleexy/protoc-gen-es6/generator"
)


func TestRS6(t *testing.T) {
	req, err := os.Open("../code_generator_request.pb.bin")
	if err != nil {
		t.Fatal(err)
	}

	fs := afero.NewOsFs()
	res := &bytes.Buffer{}

	pgs.Init(
		pgs.ProtocInput(req),  // use the pre-generated request
		pgs.ProtocOutput(res), // capture CodeGeneratorResponse
		pgs.FileSystem(fs),    // capture any custom files written directly to disk
	).RegisterModule(generator.ES6()).Render()

	// info, err := fs.Stat("protobuf/core/item.json.go")
	// log.Print(info.Name())
	log.Print(res.String())
	// check res and the fs for output
}