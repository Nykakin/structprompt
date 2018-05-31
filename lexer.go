package main

import (
    "strings"
    "unicode/utf8"
)

type TokenType int

const (
    TOKEN_ERROR TokenType = iota
    TOKEN_EOF

    TOKEN_LEFT_BRACKET
    TOKEN_LEFT_CURLY_BRACKET
    TOKEN_RIGHT_BRACKET
    TOKEN_RIGHT_CURLY_BRACKET
    TOKEN_DOT
    TOKEN_COLON
    TOKEN_COMMA
    
    TOKEN_PATH_ELEMENT
    TOKEN_METHOD
    TOKEN_METHOD_ARGUMENT
    TOKEN_FIELD
    TOKEN_FIELD_VALUE
)

const EOF rune = 0
const (
    LEFT_BRACKET = "("
    LEFT_CURLY_BRACKET = "{"
    RIGHT_BRACKET = ")"
    RIGHT_CURLY_BRACKET = "}"
    DOT = "."
    COLON = ":"
    COMMA = ","
)

type Token struct {
    Type  TokenType
    Value string
}

type Lexer struct {
    Name   string
    Input  string
    Tokens chan Token
    State  LexFn

    Start int
    Pos   int
    Width int
}

type LexFn func(*Lexer) LexFn

func (this *Lexer) Emit(tokenType TokenType) {
    this.Tokens <- Token{Type: tokenType, Value: this.Input[this.Start:this.Pos]}
    this.Start = this.Pos
}

func (this *Lexer) NextToken() Token {
	for {
		select {
		case token := <-this.Tokens:
			return token
		default:
			this.State = this.State(this)
		}
	}

	panic("NextToken reached an invalid state!!")
}

func (this *Lexer) Inc() {
    this.Pos++
    if this.Pos >= utf8.RuneCountInString(this.Input) {
        this.Emit(TOKEN_EOF)
    }
}

func (this *Lexer) InputToEnd() string {
    return this.Input[this.Pos:]
}


func (this *Lexer) IsEOF() bool {
	return this.Pos >= len(this.Input)
}

func NewLexer(input string) *Lexer {
    l := &Lexer{
        Input:  input,
        State:  LexBegin,
        Tokens: make(chan Token, 3),
    }

    return l
}

func LexBegin(lexer *Lexer) LexFn {
    return LexPathElement
}

func LexEnd(lexer *Lexer) LexFn {
    lexer.Emit(TOKEN_EOF)
    return LexEnd
}

func LexPathElement(lexer *Lexer) LexFn {
    for {
        if strings.HasPrefix(lexer.InputToEnd(), DOT) {
            lexer.Emit(TOKEN_PATH_ELEMENT)
            return LexDot
        } else if strings.HasPrefix(lexer.InputToEnd(), LEFT_BRACKET) {
            lexer.Emit(TOKEN_METHOD)            
            return LexMethod
        } 

        if !lexer.IsEOF() {
            lexer.Inc()
        } else {
            return LexEnd
        }
    }
}

func LexDot(lexer *Lexer) LexFn {
    lexer.Pos += len(DOT)
    lexer.Emit(TOKEN_DOT)
    return LexPathElement
}

func LexMethod(lexer *Lexer) LexFn {
    lexer.Pos += len(LEFT_BRACKET)
    lexer.Emit(TOKEN_LEFT_BRACKET)
    return LexMethodArguments
}

func LexMethodArguments(lexer *Lexer) LexFn {
    for {
        if strings.HasPrefix(lexer.InputToEnd(), RIGHT_BRACKET) {
            lexer.Emit(TOKEN_METHOD_ARGUMENT)
            lexer.Pos += len(RIGHT_BRACKET)
            lexer.Emit(TOKEN_RIGHT_BRACKET)
            return LexEnd
        }
        if strings.HasPrefix(lexer.InputToEnd(), LEFT_CURLY_BRACKET) {
            lexer.Pos += len(LEFT_CURLY_BRACKET)
            lexer.Emit(TOKEN_LEFT_CURLY_BRACKET)
            return LexStructFields
        }

        if strings.HasPrefix(lexer.InputToEnd(), COMMA) {
            lexer.Emit(TOKEN_METHOD_ARGUMENT)
            return LexMethodComma
        }

        if !lexer.IsEOF() {
            lexer.Inc()
        } else {
            return LexEnd
        }
    }
}

func LexMethodComma(lexer *Lexer) LexFn {
    lexer.Pos += len(COMMA)
    lexer.Emit(TOKEN_COMMA)
    return LexMethodArguments
}

func LexStructFields(lexer *Lexer) LexFn {
    for {
        if strings.HasPrefix(lexer.InputToEnd(), RIGHT_CURLY_BRACKET) {
            lexer.Emit(TOKEN_FIELD_VALUE)
            lexer.Pos += len(RIGHT_CURLY_BRACKET)
            lexer.Emit(TOKEN_RIGHT_CURLY_BRACKET)
            return LexMethodComma
        }
        if strings.HasPrefix(lexer.InputToEnd(), COMMA) {
            lexer.Emit(TOKEN_FIELD_VALUE)
            return LexStructComma
        }
        if strings.HasPrefix(lexer.InputToEnd(), COLON) {
            lexer.Emit(TOKEN_FIELD)
            return LexStructColon
        }

        if !lexer.IsEOF() {
            lexer.Inc()
        } else {
            return LexEnd
        }
    }
}

func LexStructComma(lexer *Lexer) LexFn {
    lexer.Pos += len(COMMA)
    lexer.Emit(TOKEN_COMMA)
    return LexStructFields
}

func LexStructColon(lexer *Lexer) LexFn {
    lexer.Pos += len(COLON)
    lexer.Emit(TOKEN_COLON)
    return LexStructFields
}
