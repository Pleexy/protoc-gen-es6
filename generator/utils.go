package generator

import (
	pgs "github.com/lyft/protoc-gen-star"
	"unicode"
	"unicode/utf8"
)

func upperFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}

func getOutputPath(name pgs.FilePath, o *Options) pgs.FilePath {
	if o.ReplaceJSOut {
		name = name.SetBase(name.BaseName()+"_pb").SetExt(".js")
	} else if o.ESModules {
		name = name.SetExt(".pb.mjs")
	} else {
		name = name.SetExt(".pb.es6")
	}
	return name
}