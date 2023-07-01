# WRITING AN INTEPRETER IN GO

## Overview
The Monkey Programing Language & Interpreter
- C-likesyntax variablebindings
- integersandbooleans
- arithmeticexpressions
- built-infunctions
- first-classandhigher-orderfunctions
- closures
- astringdatastructure
- anarraydatastructure
- ahashdatastructure

Major parts
- lexer
- parser
- AbstractSyntaxTree(AST)
- internalobjectsystem
- evaluator

## LEXING

```
Source -> Tokens -> AST
```

example.
```
souce
"let x = 5 + 5;"

lexer output
[
    LET,
    INDENTIFIER("x"),
    EQUAL_SIGN,
    INTEGER(5),
    PLUS_SIGN,
    INTEGER(5),
    SEMICOLON
]
```
