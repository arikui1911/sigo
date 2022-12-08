package sigo

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"unicode"
)

const TOKEN_EOF = 0

type Token struct {
	Tag    int
	Value  any
	Lineno int
	Column int
}

type lexResult struct {
	tok Token
	err error
}

type Lexer struct {
	src               *bufio.Reader
	ch                chan lexResult
	fileName          string
	lineno            int
	column            int
	lastNewlineColumn int
	isEOF             bool
	savedRune         rune
	hasSavedRune      bool
}

func NewLexer(src io.Reader, fileName string) *Lexer {
	l := &Lexer{
		src:      bufio.NewReader(src),
		ch:       make(chan lexResult, 1),
		fileName: fileName,
		lineno:   1,
	}
	go lex(l)
	return l
}

func (l *Lexer) Lex() (Token, error) {
	r := <-l.ch
	return r.tok, r.err
}

func lex(l *Lexer) {
	for {
		c, err := l.getc()
		if err != nil {
			break
		}
		switch {
		case unicode.IsSpace(c):
			// do nothing
		case c == '#':
			skipComment(l)
		case c == '"':
			lexString(l)
		case c == '0':
			lexPostZero(l)
		case unicode.IsDigit(c):
			lexNumber(l, c)
		case c == '_' || unicode.IsLetter(c):
			lexSymbol(l, c)
		case c == ')':
			lexPostRp(l)
		default:
			lexOperator(l, c)
		}
	}
	l.ch <- lexResult{tok: Token{TOKEN_EOF, nil, l.lineno, l.column}}
}

func skipComment(l *Lexer) {
	for {
		c, err := l.getc()
		if err != nil {
			break
		}
		if c == '\n' {
			l.ungetc(c)
			break
		}
	}
}

func lexString(l *Lexer) {
	lineno := l.lineno
	column := l.column
	buf := []rune{}
	for {
		c, err := l.getc()
		if err == io.EOF {
			l.ch <- lexResult{err: l.syntaxError(lineno, column, "unterminated string literal")}
			return
		}
		if err != nil {
			return
		}
		switch c {
		case '"':
			l.ch <- lexResult{tok: Token{TOKEN_LIT_STRING, string(buf), lineno, column}}
			return
		case '\\':
			buf = lexEscapeSequence(l, buf)
		default:
			buf = append(buf, c)
		}
	}
}

var escapeSequences = map[rune]rune{
	'n': '\n',
	't': '\t',
}

func lexEscapeSequence(l *Lexer, buf []rune) []rune {
	c, err := l.getc()
	if err != nil {
		return append(buf, '\\')
	}
	es, ok := escapeSequences[c]
	if ok {
		c = es
	}
	return append(buf, c)
}

func lexPostZero(l *Lexer) {
	lineno := l.lineno
	column := l.column
	c, err := l.getc()
	if err == io.EOF {
		l.ch <- lexResult{tok: Token{TOKEN_LIT_INT, 0, lineno, column}}
		return
	}
	if err != nil {
		return
	}
	if c == '.' {
		lexFloat(l, lineno, column, []rune{'0'})
		return
	}
	l.ungetc(c)
	l.ch <- lexResult{tok: Token{TOKEN_LIT_INT, 0, lineno, column}}
}

func lexNumber(l *Lexer, fc rune) {
	lineno := l.lineno
	column := l.column
	buf := []rune{fc}
	for {
		c, err := l.getc()
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}
		if c == '.' {
			lexFloat(l, lineno, column, buf)
			return
		}
		if !unicode.IsDigit(c) {
			l.ungetc(c)
			break
		}
		buf = append(buf, c)
	}
	i64, err := strconv.ParseInt(string(buf), 10, 64)
	if err != nil {
		l.ch <- lexResult{err: l.wrapError(lineno, column, err)}
		return
	}
	l.ch <- lexResult{tok: Token{TOKEN_LIT_INT, int(i64), lineno, column}}
}

func lexFloat(l *Lexer, lineno int, column int, buf []rune) {
	buf = append(buf, '.')
	for {
		c, err := l.getc()
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}
		if !unicode.IsDigit(c) {
			l.ungetc(c)
			break
		}
		buf = append(buf, c)
	}
	f64, err := strconv.ParseFloat(string(buf), 64)
	if err != nil {
		l.ch <- lexResult{err: l.wrapError(lineno, column, err)}
		return
	}
	l.ch <- lexResult{tok: Token{TOKEN_LIT_FLOAT, f64, lineno, column}}
}

var keywords = map[string]int{
	"if":    TOKEN_KW_IF,
	"elsif": TOKEN_KW_ELSIF,
	"else":  TOKEN_KW_ELSE,
	"while": TOKEN_KW_WHILE,
}

func lexSymbol(l *Lexer, fc rune) {
	lineno := l.lineno
	column := l.column
	buf := []rune{fc}
	for {
		c, err := l.getc()
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}
		if !(c == '_' || unicode.IsLetter(c) || unicode.IsDigit(c)) {
			l.ungetc(c)
			break
		}
		buf = append(buf, c)
	}
	tag := TOKEN_SYMBOL
	val := string(buf)
	if v, ok := keywords[val]; ok {
		tag = v
	}
	l.ch <- lexResult{tok: Token{tag, val, lineno, column}}
}

func lexPostRp(l *Lexer) {
	lineno := l.lineno
	column := l.column
	for {
		c, err := l.getc()
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}
		if c == '-' {
			c, err = l.getc()
			if err == io.EOF {
				l.ch <- lexResult{tok: Token{TOKEN_RP, ")", lineno, column}}
				lexOperator(l, '-')
				return
			}
			if err != nil {
				return
			}
			if c != '>' {
				l.ungetc(c)
				l.ch <- lexResult{tok: Token{TOKEN_RP, ")", lineno, column}}
				lexOperator(l, '-')
				return
			}
			l.ch <- lexResult{tok: Token{TOKEN_RP_AND_ARROW, ")->", lineno, column}}
			return
		}
		if c == '\n' || !unicode.IsSpace(c) {
			l.ungetc(c)
			break
		}
	}
	l.ch <- lexResult{tok: Token{TOKEN_RP, ")", lineno, column}}
}

var operators = map[string]int{
	"(":  TOKEN_LP,
	")":  TOKEN_RP,
	"[":  TOKEN_LB,
	"]":  TOKEN_RB,
	"{":  TOKEN_LC,
	"}":  TOKEN_RC,
	".":  TOKEN_DOT,
	",":  TOKEN_COMMA,
	"=":  TOKEN_ASSIGN,
	":":  TOKEN_COLON,
	";":  TOKEN_SEMICOLON,
	"!":  TOKEN_NOT,
	"+":  TOKEN_ADD,
	"-":  TOKEN_SUB,
	"*":  TOKEN_MUL,
	"/":  TOKEN_DIV,
	"%":  TOKEN_MOD,
	"==": TOKEN_EQ,
	"!=": TOKEN_NE,
	">":  TOKEN_GT,
	">=": TOKEN_GE,
	"<":  TOKEN_LT,
	"<=": TOKEN_LE,
	"+=": TOKEN_ADD_A,
	"-=": TOKEN_SUB_A,
	"*=": TOKEN_MUL_A,
	"/=": TOKEN_DIV_A,
	"%=": TOKEN_MOD_A,
	"&&": TOKEN_DAND,
	"||": TOKEN_DOR,
	"->": TOKEN_ARROW,
}

func lexOperator(l *Lexer, fc rune) {
	lineno := l.lineno
	column := l.column
	buf := []rune{fc}
	for {
		c, err := l.getc()
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}
		buf = append(buf, c)
		if _, ok := operators[string(buf)]; !ok {
			l.ungetc(c)
			buf = buf[:len(buf)-1]
			break
		}
	}
	k := string(buf)
	v, ok := operators[k]
	if !ok {
		l.ch <- lexResult{err: l.syntaxError(lineno, column, "invalid character - '%c'", buf[0])}
		return
	}
	l.ch <- lexResult{tok: Token{v, k, lineno, column}}
}

/*
 * Errors
 */

type SyntaxError struct {
	Message  string
	FileName string
	Lineno   int
	Column   int
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("%s:%d:%d: %s", e.FileName, e.Lineno, e.Column, e.Message)
}

func (l *Lexer) syntaxError(lineno int, column int, format string, args ...any) error {
	return &SyntaxError{
		Message:  fmt.Sprintf(format, args...),
		FileName: l.fileName,
		Lineno:   lineno,
		Column:   column,
	}
}

func (l *Lexer) wrapError(lineno int, column int, err error) error {
	return fmt.Errorf("%s:%d:%d: %w", l.fileName, lineno, column, err)
}

/*
 * buffered src reader
 */

func (l *Lexer) getc() (c rune, err error) {
	if l.hasSavedRune {
		l.hasSavedRune = false
		c = l.savedRune
	} else {
		c, _, err = l.src.ReadRune()
	}
	if err == io.EOF {
		l.isEOF = true
		return
	}
	if err != nil {
		l.ch <- lexResult{err: l.wrapError(l.lineno, l.column, err)}
		return
	}
	l.column++
	if c == '\n' {
		l.lastNewlineColumn = l.column
		l.lineno++
		l.column = 0
	}
	return
}

func (l *Lexer) ungetc(c rune) {
	l.hasSavedRune = true
	l.savedRune = c
	l.column--
	if c == '\n' {
		l.column = l.lastNewlineColumn
	}
}
