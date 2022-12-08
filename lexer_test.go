package sigo_test

import (
	"sigo"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantTag int
		wantVal any
	}{
		{"int literal", `123`, sigo.TOKEN_LIT_INT, 123},
		{"float literal", `12.3`, sigo.TOKEN_LIT_FLOAT, 12.3},
		{"string literal", `"Hello."`, sigo.TOKEN_LIT_STRING, "Hello."},
		{"identifier", `hoge`, sigo.TOKEN_SYMBOL, "hoge"},

		{"if", `if`, sigo.TOKEN_KW_IF, "if"},
		{"elsif", `elsif`, sigo.TOKEN_KW_ELSIF, "elsif"},
		{"else", `else`, sigo.TOKEN_KW_ELSE, "else"},
		{"while", `while`, sigo.TOKEN_KW_WHILE, "while"},

		{"left paren", `(`, sigo.TOKEN_LP, "("},
		{"right paren", `)`, sigo.TOKEN_RP, ")"},
		{"left bracket", `[`, sigo.TOKEN_LB, "["},
		{"right bracket", `]`, sigo.TOKEN_RB, "]"},
		{"left brace", `{`, sigo.TOKEN_LC, "{"},
		{"right brace", `}`, sigo.TOKEN_RC, "}"},
		{"dot", `.`, sigo.TOKEN_DOT, "."},
		{"comma", `,`, sigo.TOKEN_COMMA, ","},
		{"assign", `=`, sigo.TOKEN_ASSIGN, "="},
		{"colon", `:`, sigo.TOKEN_COLON, ":"},
		{"semicolon", `;`, sigo.TOKEN_SEMICOLON, ";"},
		{"not", `!`, sigo.TOKEN_NOT, "!"},
		{"addition", `+`, sigo.TOKEN_ADD, "+"},
		{"subtraction", `-`, sigo.TOKEN_SUB, "-"},
		{"multiplication", `*`, sigo.TOKEN_MUL, "*"},
		{"division", `/`, sigo.TOKEN_DIV, "/"},
		{"modulo", `%`, sigo.TOKEN_MOD, "%"},
		{"equal", `==`, sigo.TOKEN_EQ, "=="},
		{"not equal", `!=`, sigo.TOKEN_NE, "!="},
		{"greater than", `>`, sigo.TOKEN_GT, ">"},
		{"greater than equal", `>=`, sigo.TOKEN_GE, ">="},
		{"less than", `<`, sigo.TOKEN_LT, "<"},
		{"less than equal", `<=`, sigo.TOKEN_LE, "<="},
		{"add assign", `+=`, sigo.TOKEN_ADD_A, "+="},
		{"sub assign", `-=`, sigo.TOKEN_SUB_A, "-="},
		{"mul assign", `*=`, sigo.TOKEN_MUL_A, "*="},
		{"div assign", `/=`, sigo.TOKEN_DIV_A, "/="},
		{"mod assign", `%=`, sigo.TOKEN_MOD_A, "%="},
		{"2 amp", `&&`, sigo.TOKEN_DAND, "&&"},
		{"2 bar", `||`, sigo.TOKEN_DOR, "||"},
		{"arrow", `->`, sigo.TOKEN_ARROW, "->"},
		{"right paren and arrow", `) ->`, sigo.TOKEN_RP_AND_ARROW, ")->"},
	}

	for _, x := range tests {
		t.Run(x.name, func(t *testing.T) {
			l := sigo.NewLexer(strings.NewReader(x.src), "(test)")
			testLex(t, l, sigo.Token{x.wantTag, x.wantVal, 1, 1})
			testLex(t, l, sigo.Token{sigo.TOKEN_EOF, nil, 1, len(x.src)})
		})
	}
}

func TestLexRPandArrow(t *testing.T) {
	l := sigo.NewLexer(strings.NewReader(`) -`), "(test)")
	testLex(t, l, sigo.Token{sigo.TOKEN_RP, ")", 1, 1})
	testLex(t, l, sigo.Token{sigo.TOKEN_SUB, "-", 1, 3})
	testLex(t, l, sigo.Token{sigo.TOKEN_EOF, nil, 1, 3})

	l = sigo.NewLexer(strings.NewReader(`) -=`), "(test)")
	testLex(t, l, sigo.Token{sigo.TOKEN_RP, ")", 1, 1})
	testLex(t, l, sigo.Token{sigo.TOKEN_SUB_A, "-=", 1, 3})
	testLex(t, l, sigo.Token{sigo.TOKEN_EOF, nil, 1, 4})
}

func testLex(t *testing.T, l *sigo.Lexer, want sigo.Token) {
	tok, err := l.Lex()
	if err != nil {
		t.Error(err)
		return
	}
	if d := cmp.Diff(want, tok); d != "" {
		t.Errorf("want '-' got '+':\n%s", d)
	}
}
