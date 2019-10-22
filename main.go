package main

import (
	"github.com/Pleexy/protoc-gen-es6/generator"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"

)

func main() {
	pgs.Init(
		pgs.DebugEnv("DEBUG"),
	).RegisterModule(
		generator.ES6(),
	).RegisterPostProcessor(
		pgsgo.GoFmt(),
	).Render()
}