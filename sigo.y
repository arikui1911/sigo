%{
package sigo

%}

%union{
    token Token
    node Node
}

%type<node> program primary

%token<token> TOKEN_LIT_INT TOKEN_LIT_FLOAT TOKEN_LIT_STRING TOKEN_SYMBOL
%token<token> TOKEN_KW_IF TOKEN_KW_ELSIF TOKEN_KW_ELSE TOKEN_KW_WHILE
%token<token> TOKEN_LP TOKEN_RP TOKEN_LB TOKEN_RB TOKEN_LC TOKEN_RC TOKEN_DOT
              TOKEN_COMMA TOKEN_ASSIGN TOKEN_COLON TOKEN_SEMICOLON TOKEN_NOT
              TOKEN_ADD TOKEN_SUB TOKEN_MUL TOKEN_DIV TOKEN_MOD TOKEN_EQ
              TOKEN_NE TOKEN_GT TOKEN_GE TOKEN_LT TOKEN_LE TOKEN_ADD_A
              TOKEN_SUB_A TOKEN_MUL_A TOKEN_DIV_A TOKEN_MOD_A TOKEN_DAND
              TOKEN_DOR TOKEN_ARROW


%%

program
: primary
{
    $$ = finishAST(yylex, $1)
}
;

primary
: TOKEN_LIT_INT
{
    $$ = &LiteralIntNode{$1.Lineno, $1.Column, $1.Value.(int)}
}
;

%%

