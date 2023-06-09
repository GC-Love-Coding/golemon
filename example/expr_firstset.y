%{
package main
%}

%%

expression: term '+' term
;
term: factor '*' factor
;
factor: digit | '(' expression ')'
;
digit: '0'|'1'|'2'|'3'|'4'|'5'|'6'|'7'|'8'|'9'
;

%%