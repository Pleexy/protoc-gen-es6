package generator

import (
	"log"
	"strings"
	"testing"
)

func TestIndent(t *testing.T) {
	a := strings.Split("","\n")
	b := strings.Split("abc", "\n")
	c := strings.Split("abc\n", "\n")
	log.Print(a,b,c)
}