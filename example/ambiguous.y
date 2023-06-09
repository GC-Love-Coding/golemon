%type <a> exp

%%

exp: exp '+' exp { $$ = newast('+', $1, $3); }
   | exp '-' exp { $$ = newast('-', $1, $3); }
   | exp '*' exp { $$ = newast('*', $1, $3); }
   | exp '/' exp { $$ = newast('/', $1, $3); }
   | '|' exp     { $$ = newast('|', $2, NULL); }
   | '(' exp ')' { $$ = $2; }
   | '-' exp     { $$ = newast('M', $2, NULL); }
   | NUMBER      { $$ = newnum($1); }
;