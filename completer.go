package main

import (
    "fmt"
    "reflect"
    "strings"

	"github.com/c-bata/go-prompt"
)

func parse(s string) {
/*
    defer func() {
        if r := recover(); r != nil {
            suggestions = []prompt.Suggest{}
        }
    }()
*/
    showStructArgumentPrompt = false

    var token Token
    var tokenValue string
    
    var methodToCall reflect.Value
    var methodToCallType reflect.Type

    argumentCount := 0

    currentValue = reflect.ValueOf(&m).Elem()
    getFields(currentValue)
    fields = []string{}

    lexer := NewLexer(s)

    for {
        token = lexer.NextToken()
        tokenValue = strings.TrimSpace(token.Value)

        if token.Type == TOKEN_EOF {
            break
        }

        switch token.Type {
        case TOKEN_DOT:
            currentValue = reflect.Indirect(currentValue.FieldByName(currentField))
            fields = append(fields, currentField)
            getFields(currentValue)
        case TOKEN_PATH_ELEMENT:
            currentField = tokenValue
        case TOKEN_METHOD:
            methodToCall = currentValue.MethodByName(tokenValue)
            methodToCallType = methodToCall.Type()
        case TOKEN_COMMA:
            if !showStructArgumentPrompt {
                argumentCount += 1
            }
        case TOKEN_METHOD_ARGUMENT:
        case TOKEN_FIELD:
        case TOKEN_FIELD_VALUE:
        case TOKEN_LEFT_CURLY_BRACKET:
            showStructArgumentPrompt = true
            structArgumentPrompt(methodToCallType.In(argumentCount))
        case TOKEN_RIGHT_CURLY_BRACKET:            
        }
    }
}

func completer(t prompt.Document) []prompt.Suggest {
    parse(t.TextBeforeCursor())

    if !showStructArgumentPrompt {
        p := strings.Split(t.TextBeforeCursor(), "(")
    	return prompt.FilterHasPrefix(suggestions, p[0], true)
    } else {
        return suggestions
    }
}

// TODO: encapsulate
var suggestions []prompt.Suggest
var currentValue reflect.Value
var currentField string
var fields []string
var m Modules
var showStructArgumentPrompt bool

func structArgumentPrompt(arg reflect.Type) {
    fmt.Println(arg.NumField())
    suggestions = make([]prompt.Suggest, arg.NumField())
    for i := 0; i < arg.NumField(); i++ {
        suggestions[i] = prompt.Suggest{
            Text: arg.Field(i).Name,
            Description: arg.Field(i).Type.String(),
        }
    }
}

func getFields(s reflect.Value) {
    suggestions = make([]prompt.Suggest, s.NumField() + s.NumMethod())

    typeOfT := s.Type()
    prefix := strings.Join(fields, ".")
    if len(prefix) > 0 {
        prefix += "."
    }
    for i := 0; i < s.NumField(); i++ {
        suggestions[i] = prompt.Suggest{
            Text: prefix + typeOfT.Field(i).Name,
/*            Description: typeOfT.Field(i).Name,*/
        }
    }

    for i := 0; i < s.NumMethod(); i++ {
        suggestions[s.NumField() + i] = prompt.Suggest{
            Text: prefix + typeOfT.Method(i).Name + "()",
            Description: methodSignature(typeOfT.Method(i).Type),
        }
    }
}

func methodSignature(m reflect.Type) string {
    arguments := make([]string, m.NumIn())
    returnTypes := make([]string, m.NumOut())
    for i := 0; i < m.NumIn(); i++ {
        arguments[i] = m.In(i).Name()
    }
    for i := 0; i < m.NumOut(); i++ {
        returnTypes[i] = m.Out(i).Name()
    }

    switch(len(returnTypes)) {
    case 0 :
        return fmt.Sprintf("(%s)", strings.Join(arguments[1:], ", "))
    case 1:
        return fmt.Sprintf("(%s) %s", strings.Join(arguments[1:], ", "), returnTypes[0])
    default:
        return fmt.Sprintf("(%s) (%s)", strings.Join(arguments[1:] , ", "), strings.Join(returnTypes, ", "))
    } 
}
