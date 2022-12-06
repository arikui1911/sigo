package sigo_test

import (
	"sigo"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParser(t *testing.T){
    l := sigo.NewLexer(strings.NewReader(`123`), "(test)")
    tree, err := sigo.Parse(l)
    if err != nil {
        t.Error(err)
        return
    }
    if d := cmp.Diff(&sigo.LiteralIntNode{1, 1, 123}, tree); d != "" {
        t.Errorf("want '-' got '+':\n%s", d)
    }
}
