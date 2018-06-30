%{
package goval

%}

%union {
  token     tokenInf
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

/* Underlines are used, so that the generated go code stays unexported */
%token<token> _LITERAL_NIL    // nil
%token<token> _LITERAL_BOOL   // true false
%token<token> _LITERAL_NUMBER // 42 4.2 4e2 4.2e2
%token<token> _LITERAL_STRING // "text" 'text'
%token<token> _IDENT
%token<token> _AND            // &&
%token<token> _OR             // ||
%token<token> _EQL            // ==
%token<token> _NEQ            // !=
%token<token> _LSS            // <
%token<token> _GTR            // >
%token<token> _LEQ            // <=
%token<token> _GEQ            // >=
%token<token> _SHL            // <<
%token<token> _SHR            // >>
%token<token> _BIT_NOT        // ~
%token<token> _IN             // in

/* Operator precedence is taken from C/C++: http://en.cppreference.com/w/c/language/operator_precedence */

%left  _OR
%left  _AND
%left  '|'
%left  '^'
%left  '&'
%left  _EQL _NEQ
%left  _LSS _LEQ _GTR _GEQ
%left  _SHL _SHR
%left  '+' '-'
%left  '*' '/' '%'
%right '!' _BIT_NOT
%left  _IN
%left  '.' '[' ']'

%%

program
  : expr
  {
    $$ = $1
    yylex.(*lexer).result = $$
  }
  ;

expr
  : literal
  | math
  | logic
  | bitManipulation
  | varAccess
  | '(' expr ')'            { $$ = $2 }
  | _IDENT '(' ')'          { $$ = callFunction(yylex.(*lexer).functions, $1.literal, []interface{}{}) }
  | _IDENT '(' exprList ')' { $$ = callFunction(yylex.(*lexer).functions, $1.literal, $3) }
  ;

literal
  : _LITERAL_NIL          { $$ = nil }
  | _LITERAL_BOOL         { $$ = $1.value }
  | _LITERAL_NUMBER       { $$ = $1.value }
  | _LITERAL_STRING       { $$ = $1.value }
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
  | expr _EQL expr        { $$ = deepEqual($1, $3) }
  | expr _NEQ expr        { $$ = !deepEqual($1, $3) }
  | expr _LSS expr        { $$ = compare($1, $3, "<") }
  | expr _GTR expr        { $$ = compare($1, $3, ">") }
  | expr _LEQ expr        { $$ = compare($1, $3, "<=") }
  | expr _GEQ expr        { $$ = compare($1, $3, ">=") }
  | expr _AND expr        { left := asBool($1); right := asBool($3); $$ = left && right }
  | expr _OR expr         { left := asBool($1); right := asBool($3); $$ = left || right }
  ;

bitManipulation
  : expr '|' expr         { $$ = asInteger($1) | asInteger($3) }
  | expr '&' expr         { $$ = asInteger($1) & asInteger($3) }
  | expr '^' expr         { $$ = asInteger($1) ^ asInteger($3) }
  | expr _SHL expr         { l := asInteger($1); r := asInteger($3); if r >= 0 { $$ = l << uint(r) } else {$$ = l >> uint(-r)} }
  | expr _SHR expr         { l := asInteger($1); r := asInteger($3); if r >= 0 { $$ = l >> uint(r) } else {$$ = l << uint(-r)} }
  | _BIT_NOT expr          { $$ = ^asInteger($2) }
  ;

varAccess
  : _IDENT                       { $$ = accessVar(yylex.(*lexer).variables, $1.literal) }
  | expr '.' _IDENT              { $$ = accessField($1, $3.literal) }
  | expr '[' expr ']'            { $$ = accessField($1, $3) }
  | expr _IN expr                { $$ = arrayContains($3, $1) }
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

