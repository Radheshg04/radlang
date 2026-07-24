## Functions
- RL001(err): Cant redeclare builtin functions
- RL002(err): Cant declare multiple funcs with same name

## Statements
- RL401(err): Args passed must match declared params
- RL402(err): "if" expression must be boolean
- RL403(err): no unused declared variables
- RL404(err): redeclared variables
- RL405(err): idents must be declared before assignment (walrus exception as it declares on assignment)
- RL406(err): there must be atleast 1 new variable on left side of walrus
- RL407(err): too many return vals
- RL408(err): not enough return vals
- RL409(err): can only return stmt in blocks
- RL410(err): Bad Return Type
- RL411(err): can use jump stmt only in "for"
- RL412(err): exxpressions in expr_stmt must be fn_call or postfix
- RL413(warn): warn for unreachable code
- RL414(err): performing postfix op on non-numeric ident

## Expressions
- RL501(err): cmpr op, arithmetic ops can only be used with declared idents, or lit_vals
<!-- - RL018(err): error literal can only have string as args -> RL014 -->

## Misc
<!-- - define err(error string) and print(...any) be treated as builtins -->
