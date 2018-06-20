%{
package main

%}

%union {
  token     Token
  expr      interface{}
  exprList  []interface{}
}


%start program

%type<expr> program
%type<expr> expr
%type<expr> literal
%type<expr> math
%type<expr> logic
%type<expr> varAccess
%type<exprList> exprList

%token<token> LITERAL_BOOL   // true false
%token<token> LITERAL_NUMBER // 42 4.2 4e2 4.2e2
%token<token> LITERAL_STRING // "text" 'text'
%token<token> IDENT
%token<token> AND            // &&
%token<token> OR             // ||
%token<token> EQL            // ==
%token<token> NEQ            // !=
%token<token> LSS            // <
%token<token> GTR            // >
%token<token> LEQ            // <=
%token<token> GEQ            // >=

/* Operator precedence is taken from c: http://en.cppreference.com/w/c/language/operator_precedence */

%left  OR
%left  AND
%left  EQL NEQ
%left  LSS LEQ GTR GEQ
%left  '+' '-'
%left  '*' '/'
%right '!'
%left  '.' '[' ']'

%%

program
  : expr
  {
    $$ = $1
    yylex.(*Lexer).result = $$
  }
  ;

expr
  : literal
  | math
  | logic
  | varAccess
  | '(' expr ')'           { $$ = $2 }
  | IDENT '(' ')'          { $$ = callFunction(yylex.(*Lexer).functions, $1.literal, []interface{}{}) }
  | IDENT '(' exprList ')' { $$ = callFunction(yylex.(*Lexer).functions, $1.literal, $3) }
  ;

literal
  : LITERAL_BOOL          { $$ = $1.value }
  | LITERAL_NUMBER        { $$ = $1.value }
  | LITERAL_STRING        { $$ = $1.value }
  | '[' ']'               { $$ = []interface{}{} }
  | '[' exprList ']'      { $$ = $2 }
  ;

math
  : '-' expr %prec  '*'   { $$ = unaryMinus($2)  }  /* unary minus has higher precedence */
  | expr '+' expr         { $$ = add($1, $3) }
  | expr '-' expr         { $$ = sub($1, $3) }
  | expr '*' expr         { $$ = mul($1, $3) }
  | expr '/' expr         { $$ = div($1, $3) }
  ;

logic
  : '!' expr              { $$ = !asBool($2) }
  | expr EQL expr         { $$ = deepEqual($1, $3) }
  | expr NEQ expr         { $$ = !deepEqual($1, $3) }
  | expr AND expr         { left := asBool($1); right := asBool($3); $$ = left && right }
  | expr OR expr          { left := asBool($1); right := asBool($3); $$ = left || right }

varAccess
  : IDENT                 { $$ = accessVar(yylex.(*Lexer).variables, $1.literal) }
  | expr '.' IDENT        { $$ = accessField($1, $3.literal) }
  | expr '[' expr ']'     { $$ = accessField($1, $3) }
  ;

exprList
  : expr                  { $$ = []interface{}{$1} }
  | exprList ',' expr     { $$ = append($1, $3) }
  ;

%%

