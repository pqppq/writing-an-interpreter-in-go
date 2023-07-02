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

Lexer(Lexical Analyzer): Lexer will take source code as input. And then output the tokens that represent the source code.

```
Source -> Tokens -> AST
```

example
```
input: source code
"let x = 5 + 5;"

output: tokens
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

## PARSING

From wiki, ...
> A parser takes input data and builds a data structure - often some kind of parse tree, AST(abstruct syntax tree) or other hierarchical structure - giving a structual representation of the input, checking for correct syntax in the process. ... The parser is often preceded by a separrate lexical analyzer, which creates tokens from the sequence of input characters.

example
```
> var input = 'if ( 3 * 5 > 10 ) { return "hello"; } else { return "goodbye"; }';
> var tokens = MagicLexer.parse(input);
> MagicParser.parse(tokens);
{
    type: "if-statement",
    condition: {
        type: "operator-expression",
        operator: ">",
        left: {
            type: "operator-expression",
            operator: "*",
            left: {type: "integer-literal", value: 3},
            right: {type: "integer-literal", value: 5}
        }, 
        right: {type: "integer-literal", value: 10}
    }, 
    consequence: {
        type: "return-statement",
        returnValue: {type: "string-literal", value: "hello"}
    },
    alternative: {
        type: "return-statement",
        returnValue: {type: "string-literal", value: "goodbye"} }
}
```
