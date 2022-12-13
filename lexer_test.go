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
		newline bool
	}{
		{"int literal", `123`, sigo.TOKEN_LIT_INT, 123, true},
		{"float literal", `12.3`, sigo.TOKEN_LIT_FLOAT, 12.3, true},
		{"string literal", `"Hello."`, sigo.TOKEN_LIT_STRING, "Hello.", true},
		{"identifier", `hoge`, sigo.TOKEN_SYMBOL, "hoge", true},

		{"if", `if`, sigo.TOKEN_KW_IF, "if", false},
		{"elsif", `elsif`, sigo.TOKEN_KW_ELSIF, "elsif", false},
		{"else", `else`, sigo.TOKEN_KW_ELSE, "else", false},
		{"while", `while`, sigo.TOKEN_KW_WHILE, "while", false},

		{"left paren", `(`, sigo.TOKEN_LP, "(", false},
		{"right paren", `)`, sigo.TOKEN_RP, ")", true},
		{"left bracket", `[`, sigo.TOKEN_LB, "[", false},
		{"right bracket", `]`, sigo.TOKEN_RB, "]", true},
		{"left brace", `{`, sigo.TOKEN_LC, "{", false},
		{"right brace", `}`, sigo.TOKEN_RC, "}", true},
		{"dot", `.`, sigo.TOKEN_DOT, ".", false},
		{"comma", `,`, sigo.TOKEN_COMMA, ",", false},
		{"assign", `=`, sigo.TOKEN_ASSIGN, "=", false},
		{"colon", `:`, sigo.TOKEN_COLON, ":", false},
		{"semicolon", `;`, sigo.TOKEN_SEMICOLON, ";", false},
		{"not", `!`, sigo.TOKEN_NOT, "!", false},
		{"addition", `+`, sigo.TOKEN_ADD, "+", false},
		{"subtraction", `-`, sigo.TOKEN_SUB, "-", false},
		{"multiplication", `*`, sigo.TOKEN_MUL, "*", false},
		{"division", `/`, sigo.TOKEN_DIV, "/", false},
		{"modulo", `%`, sigo.TOKEN_MOD, "%", false},
		{"equal", `==`, sigo.TOKEN_EQ, "==", false},
		{"not equal", `!=`, sigo.TOKEN_NE, "!=", false},
		{"greater than", `>`, sigo.TOKEN_GT, ">", false},
		{"greater than equal", `>=`, sigo.TOKEN_GE, ">=", false},
		{"less than", `<`, sigo.TOKEN_LT, "<", false},
		{"less than equal", `<=`, sigo.TOKEN_LE, "<=", false},
		{"add assign", `+=`, sigo.TOKEN_ADD_A, "+=", false},
		{"sub assign", `-=`, sigo.TOKEN_SUB_A, "-=", false},
		{"mul assign", `*=`, sigo.TOKEN_MUL_A, "*=", false},
		{"div assign", `/=`, sigo.TOKEN_DIV_A, "/=", false},
		{"mod assign", `%=`, sigo.TOKEN_MOD_A, "%=", false},
		{"2 amp", `&&`, sigo.TOKEN_DAND, "&&", false},
		{"2 bar", `||`, sigo.TOKEN_DOR, "||", false},
		{"arrow", `->`, sigo.TOKEN_ARROW, "->", false},
	}

	for _, x := range tests {
		t.Run(x.name, func(t *testing.T) {
			l := sigo.NewLexer(strings.NewReader(x.src+"\n"), "(test)")
			testLex(t, l, sigo.Token{x.wantTag, x.wantVal, 1, 1})
			if x.newline {
				testLex(t, l, sigo.Token{sigo.TOKEN_NL, "\n", 1, len(x.src) + 1})
			}
			testLex(t, l, sigo.Token{sigo.TOKEN_EOF, nil, 2, 1})
		})
	}
}

func TestLexNL(t *testing.T) {
	l := sigo.NewLexer(strings.NewReader("123\n"), "(test)")
	testLex(t, l, sigo.Token{sigo.TOKEN_LIT_INT, 123, 1, 1})
	testLex(t, l, sigo.Token{sigo.TOKEN_NL, "\n", 1, 4})
	testLex(t, l, sigo.Token{sigo.TOKEN_EOF, nil, 2, 1})

	l = sigo.NewLexer(strings.NewReader("\"Hello\"\n"), "(test)")
	testLex(t, l, sigo.Token{sigo.TOKEN_LIT_STRING, "Hello", 1, 1})
	testLex(t, l, sigo.Token{sigo.TOKEN_NL, "\n", 1, 8})
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
