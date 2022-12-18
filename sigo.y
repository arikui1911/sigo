%{
package sigo

%}

%union{
    token Token
    node Node
}

%type<node> program block compstmt stmts stmt expr primary


%token<token> TOKEN_MIN

%token<token> TOKEN_LIT_INT TOKEN_LIT_FLOAT TOKEN_LIT_STRING TOKEN_SYMBOL
%token<token> TOKEN_KW_IF TOKEN_KW_ELSIF TOKEN_KW_ELSE TOKEN_KW_WHILE
%token<token> TOKEN_LP TOKEN_RP TOKEN_LB TOKEN_RB TOKEN_LC TOKEN_RC TOKEN_DOT
              TOKEN_COMMA TOKEN_ASSIGN TOKEN_COLON TOKEN_SEMICOLON TOKEN_NOT
              TOKEN_ADD TOKEN_SUB TOKEN_MUL TOKEN_DIV TOKEN_MOD TOKEN_EQ
              TOKEN_NE TOKEN_GT TOKEN_GE TOKEN_LT TOKEN_LE TOKEN_ADD_A
              TOKEN_SUB_A TOKEN_MUL_A TOKEN_DIV_A TOKEN_MOD_A TOKEN_DAND
              TOKEN_DOR TOKEN_ARROW
%token<token> TOKEN_RP_AND_ARROW TOKEN_NL

%token<token> TOKEN_MAX



%%

program
:
{
    $$ = nil
}
| compstmt
{
    $$ = finishAST(yylex, $1)
}
;

compstmt
: stmts terms_opt
;

block
: TOKEN_LC TOKEN_RC
{
    $$ = nil
}
| TOKEN_LC compstmt TOKEN_RC
{
    $$ = finishBlock($2)
}
;

stmts
: stmt
| stmts terms stmt
{
    $$ = appendStatement($1, $3)
}
;

terms_opt
:
| terms
;

terms
: term
| terms term
;

term
: TOKEN_NL
| TOKEN_SEMICOLON
;

stmt
: expr
| block
;

expr
: primary
;

primary
: TOKEN_LIT_INT
{
    $$ = &LiteralIntNode{$1.Lineno, $1.Column, $1.Value.(int)}
}
;

%%

