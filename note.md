# WRITING AN INTEPRETER IN GO

## Overview
The Monkey Programing Language & Interpreter
- C-likesyntax variablebindings
- integers and booleans
- arithmetic expressions
- built-in functions
- first-class and higher-orderfunctions
- closures
- a string data structure
- a narray data structure
- a hash data structure

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

[Pratt Parser](https://ja.wikipedia.org/wiki/Pratt%E3%83%91%E3%83%BC%E3%82%B5) 

- prefix operator
    - an operator in front of its operand
    - ex. `--5`
- postfix operator
    - an operator after its operand
    - ex `x++`
- infix operator
    - appear in binary expressions
    - ex `5 * 8`

## EVALUATION

Two apploach
- Traverse the AST, visit each node and do what the node signifies: print a string, add two numbers, execute a function's body - all on the fly.
- Compilesthe bytecode and use a virtual machine to evalueate.

Tree-walking strategy

pseudo code
```
function eval(astNode) {
  if (astNode is integerliteral) {
    return astNode.integerValue
  } else if (astNode is booleanLiteral) {
    return astNode.booleanValue
  } else if (astNode is infixExpression) {
    leftEvaluated = eval(astNode.Left)
    rightEvaluated = eval(astNode.Right)

    if astNode.Operator == "+" {
      return leftEvaluated + rightEvaluated
    } else if ast.Operator == "-" {
      return leftEvaluated - rightEvaluated
    }
  }
}
```


