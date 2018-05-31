package main

import (
    "fmt"
    "reflect"
    "strconv"
    "strings"
)

func executor(s string) {
/*    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Failed to call method")
        }
    }()
*/
    currentValue := reflect.ValueOf(&m).Elem()
    currentPathField := ""
    
    var methodToCall reflect.Value
    var methodToCallType reflect.Type
    var arguments []reflect.Value
    var currentArgumentKind reflect.Kind
    
    var currentStructArgument reflect.Value
    var currentStructArgumentType reflect.Type
    var currentStructArgumentFieldKind reflect.Kind
    var currentStructArgumentField string

    var token Token
    var tokenValue string

    lexer := NewLexer(s)

    for {
        token = lexer.NextToken()
        tokenValue = strings.TrimSpace(token.Value)

        if token.Type == TOKEN_EOF {
            break
        }

        switch token.Type {
        case TOKEN_DOT:
            currentValue = reflect.Indirect(currentValue.FieldByName(currentPathField))
        case TOKEN_PATH_ELEMENT:
            currentPathField = tokenValue
        case TOKEN_METHOD:
            methodToCall = currentValue.MethodByName(tokenValue)
            methodToCallType = methodToCall.Type()
        case TOKEN_METHOD_ARGUMENT:
            currentArgumentKind = methodToCallType.In(len(arguments)).Kind()
            arguments = append(arguments, convertArgument(currentArgumentKind, tokenValue))
        case TOKEN_LEFT_CURLY_BRACKET:
            currentStructArgumentType = methodToCallType.In(len(arguments))
            currentStructArgument = reflect.New(currentStructArgumentType).Elem()
        case TOKEN_FIELD:
            currentStructArgumentField = tokenValue
            currentStructArgumentFieldKind = currentStructArgument.FieldByName(tokenValue).Kind()
        case TOKEN_FIELD_VALUE:
            currentStructArgument.FieldByName(currentStructArgumentField).Set(
                convertArgument(currentStructArgumentFieldKind, tokenValue),
            )
        case TOKEN_RIGHT_CURLY_BRACKET:
            arguments = append(arguments, currentStructArgument) 
        }
    }

    callMethod(methodToCall, arguments)
}

func convertArgument(kind reflect.Kind, s string) reflect.Value {
    switch kind {
    case reflect.Int:
        res, err := strconv.Atoi(s)
        if err != nil {
            panic("DSADS")
        }
        return reflect.ValueOf(res)
    // ...

    default:
        return reflect.ValueOf(s)
    }
}

func callMethod(methodToCall reflect.Value, arguments []reflect.Value) {
    ret := methodToCall.Call(arguments)

    for _, v := range ret {
        switch v.Kind() {
        case reflect.Slice:
            fmt.Println("\t[")
            for i := 0; i < v.Len(); i++ {
                fmt.Printf("\t\t%+v\n", v.Index(i))
            }
            fmt.Println("\t]")
        default:
            fmt.Printf("\t%+v\n", v.Interface())
        }


    }
    fmt.Println()
}
