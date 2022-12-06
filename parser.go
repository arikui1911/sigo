package sigo

import "fmt"

type Node interface{}

type LiteralIntNode struct {
    Lineno int
    Column int
    Value int
}

type adaptor struct {
    tree Node
    tok Token
    err error
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
    a.err = fmt.Errorf("%s:%d:%d: %s", a.lexer.fileName, a.tok.Lineno, a.tok.Column, e)
    panic(bailout{})
}

func finishAST(yylex yyLexer, tree Node) Node {
    yylex.(*adaptor).tree = tree
    return tree
}

func Parse(l *Lexer) (Node, error) {
    a := &adaptor{lexer: l}
    defer func(){
        e := recover()
        if e == nil {
            return
        }
        if _, ok := e.(bailout); !ok {
            panic(e)
        }
    }()
    yyParse(a)
    return a.tree, a.err
}

