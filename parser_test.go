package sigo_test

import (
	"sigo"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name string
		src  string
		ast  sigo.Node
	}{
		{
			"int literal expr stmt",
			"123\n",
			&sigo.LiteralIntNode{1, 1, 123},
		},
		{
			"int literal expr several stmts",
			"123\n456\n",
			&sigo.BlockNode{[]sigo.Node{&sigo.LiteralIntNode{1, 1, 123}, &sigo.LiteralIntNode{2, 1, 456}}},
		},
		{
			"empty source",
			"\n",
			nil,
		},
		{
			"empty block (one line)",
			"{}\n",
			nil,
		},
		{
			"empty block",
			"{\n}\n",
			nil,
		},
		{
			"one line block",
			"{ 123 }\n",
			&sigo.BlockNode{[]sigo.Node{&sigo.LiteralIntNode{1, 3, 123}}},
		},
		{
			"int literal expr several stmts block",
			`
{
  123
  456
}
			`,
			&sigo.BlockNode{[]sigo.Node{&sigo.LiteralIntNode{3, 3, 123}, &sigo.LiteralIntNode{4, 3, 456}}},
		},
		{
			"int literal expr several stmts block (2)",
			`
{ 123
  456 }
			`,
			&sigo.BlockNode{[]sigo.Node{&sigo.LiteralIntNode{2, 3, 123}, &sigo.LiteralIntNode{3, 3, 456}}},
		},
		{
			"while stmt",
			`
while 100 {
  200
}
			`,
			&sigo.WhileNode{2, 1, &sigo.LiteralIntNode{2, 7, 100}, &sigo.BlockNode{[]sigo.Node{&sigo.LiteralIntNode{3, 3, 200}}}},
		},
	}

	for _, x := range tests {
		t.Run(x.name, func(t *testing.T) {
			l := sigo.NewLexer(strings.NewReader(x.src), "(test)")
			tree, err := sigo.Parse(l)
			if err != nil {
				t.Error(err)
				return
			}
			if d := cmp.Diff(x.ast, tree); d != "" {
				t.Errorf("want '-' got '+':\n%s", d)
			}
		})
	}
}
