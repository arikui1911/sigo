package sigo

import "fmt"

//go:generate goyacc sigo.y

type Node interface{}

type BlockNode struct {
	Statements []Node
}

func appendStatement(stmts Node, stmt Node) Node {
	if b, ok := stmts.(*BlockNode); ok {
		b.Statements = append(b.Statements, stmt)
		return b
	}
	return &BlockNode{Statements: []Node{stmts, stmt}}
}

func finishBlock(node Node) *BlockNode {
	if b, ok := node.(*BlockNode); ok {
		return b
	}
	return &BlockNode{Statements: []Node{node}}
}

type LiteralIntNode struct {
	Lineno int
	Column int
	Value  int
}

type adaptor struct {
	tree  Node
	tok   Token
	err   error
	lexer *Lexer
}

type bailout struct{}

func (a *adaptor) Lex(lval *yySymType) int {
	a.tok, a.err = a.lexer.Lex()
	if a.err != nil {
		panic(bailout{})
	}
	lval.token = a.tok
	return a.tok.Tag
}

func (a *adaptor) Error(e string) {
	a.err = fmt.Errorf(
		"%s:%d:%d: %s - %#v(%s)",
		a.lexer.fileName, a.tok.Lineno, a.tok.Column, e,
		a.tok.Value, getTokenName(a.tok.Tag),
	)
	panic(bailout{})
}

func getTokenName(tag int) string {
	if tag <= TOKEN_MIN || tag >= TOKEN_MAX {
		return fmt.Sprintf("token-%d", tag)
	}
	return yyToknames[tag-TOKEN_MIN+3]
}

func finishAST(yylex yyLexer, tree Node) Node {
	yylex.(*adaptor).tree = tree
	return tree
}

func Parse(l *Lexer) (tree Node, err error) {
	a := &adaptor{lexer: l}
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if _, ok := e.(bailout); !ok {
			panic(e)
		}
		err = a.err
	}()
	yyParse(a)
	tree = a.tree
	return
}
