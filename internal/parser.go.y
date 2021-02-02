%{
package internal

%}

%union {
  token     Token
  expr      interface{}
  exprList  []interface{}
  exprMap   map[string]interface{}
}


%start program

%type<expr> program
%type<expr> expr
%type<expr> literal
%type<expr> math
%type<expr> logic
%type<expr> bitManipulation
%type<expr> varAccess
%type<exprList> exprList
%type<exprMap> exprMap

%token<token> LITERAL_NIL    // nil
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
%token<token> SHL            // <<
%token<token> SHR            // >>
%token<token> BIT_NOT        // ~
%token<token> IN             // in

/* Operator precedence is taken from C/C++: http://en.cppreference.com/w/c/language/operator_precedence */

%right '?' ':'
%left  OR
%left  AND
%left  '|'
%left  '^'
%left  '&'
%left  EQL NEQ
%left  LSS LEQ GTR GEQ
%left  SHL SHR
%left  '+' '-'
%left  '*' '/' '%'
%right '!' BIT_NOT
%left  IN
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
  | bitManipulation
  | varAccess
  | expr '?' expr ':' expr { if asBool($1) { $$ = $3 } else { $$ = $5 } }
  | '(' expr ')'           { $$ = $2 }
  | IDENT '(' ')'          { $$ = callFunction(yylex.(*Lexer).functions, $1.literal, []interface{}{}) }
  | IDENT '(' exprList ')' { $$ = callFunction(yylex.(*Lexer).functions, $1.literal, $3) }
  ;

literal
  : LITERAL_NIL           { $$ = nil }
  | LITERAL_BOOL          { $$ = $1.value }
  | LITERAL_NUMBER        { $$ = $1.value }
  | LITERAL_STRING        { $$ = $1.value }
  | '[' ']'               { $$ = []interface{}{} }
  | '[' exprList ']'      { $$ = $2 }
  | '{' '}'               { $$ = map[string]interface{}{} }
  | '{' exprMap '}'       { $$ = $2 }
  ;

math
  : '-' expr %prec  '!'   { $$ = unaryMinus($2)  }  /* unary minus has higher precedence */
  | expr '+' expr         { $$ = add($1, $3) }
  | expr '-' expr         { $$ = sub($1, $3) }
  | expr '*' expr         { $$ = mul($1, $3) }
  | expr '/' expr         { $$ = div($1, $3) }
  | expr '%' expr         { $$ = mod($1, $3) }
  ;

logic
  : '!' expr              { $$ = !asBool($2) }
  | expr EQL expr         { $$ = deepEqual($1, $3) }
  | expr NEQ expr         { $$ = !deepEqual($1, $3) }
  | expr LSS expr         { $$ = compare($1, $3, "<") }
  | expr GTR expr         { $$ = compare($1, $3, ">") }
  | expr LEQ expr         { $$ = compare($1, $3, "<=") }
  | expr GEQ expr         { $$ = compare($1, $3, ">=") }
  | expr AND expr         { left := asBool($1); right := asBool($3); $$ = left && right }
  | expr OR expr          { left := asBool($1); right := asBool($3); $$ = left || right }
  ;

bitManipulation
  : expr '|' expr         { $$ = asInteger($1) | asInteger($3) }
  | expr '&' expr         { $$ = asInteger($1) & asInteger($3) }
  | expr '^' expr         { $$ = asInteger($1) ^ asInteger($3) }
  | expr SHL expr         { l := asInteger($1); r := asInteger($3); if r >= 0 { $$ = l << uint(r) } else {$$ = l >> uint(-r)} }
  | expr SHR expr         { l := asInteger($1); r := asInteger($3); if r >= 0 { $$ = l >> uint(r) } else {$$ = l << uint(-r)} }
  | BIT_NOT expr          { $$ = ^asInteger($2) }
  ;

varAccess
  : IDENT                        { $$ = accessVar(yylex.(*Lexer).variables, $1.literal) }
  | expr '.' IDENT               { $$ = accessField($1, $3.literal) }
  | expr '[' expr ']'            { $$ = accessField($1, $3) }
  | expr IN expr                 { $$ = arrayContains($3, $1) }
  | expr '[' expr ':' expr ']'   { $$ = slice($1, $3, $5) }
  | expr '['      ':' expr ']'   { $$ = slice($1, nil, $4) }
  | expr '[' expr ':'      ']'   { $$ = slice($1, $3, nil) }
  | expr '['      ':'      ']'   { $$ = slice($1, nil, nil) }
  ;

exprList
  : expr                  { $$ = []interface{}{$1} }
  | exprList ',' expr     { $$ = append($1, $3) }
  ;

exprMap
  : expr ':' expr               { $$ = make(map[string]interface{}); $$[asObjectKey($1)] = $3 }
  | exprMap ',' expr ':' expr   { $$ = addObjectMember($1, $3, $5) }
  ;

%%
