package structprompt

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

//	"git.modulus.eu/go/common/types/uuid"
)

type executor struct {
	i interface{}
}

func (e *executor) execute(s string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Failed to call method: ", r)
		}
	}()

	currentValue := reflect.ValueOf(e.i)
	currentPathField := ""
	structArgument := false

	var methodToCall reflect.Value
	var methodToCallType reflect.Type
	var arguments []reflect.Value
	var currentArgumentType reflect.Type

	var currentStructArgument reflect.Value
	var currentStructArgumentType reflect.Type
	var currentStructArgumentFieldType reflect.Type
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
			if !structArgument {
				currentArgumentType = methodToCallType.In(len(arguments))
				arguments = append(arguments, convertArgument(currentArgumentType, tokenValue))
			} else {
				structArgument = false
			}
		case TOKEN_LEFT_CURLY_BRACKET:
			currentStructArgumentType = methodToCallType.In(len(arguments))
			currentStructArgument = reflect.New(currentStructArgumentType).Elem()
		case TOKEN_FIELD:
			currentStructArgumentField = tokenValue
			currentStructArgumentFieldType = currentStructArgument.FieldByName(tokenValue).Type()
		case TOKEN_FIELD_VALUE:
			currentStructArgument.FieldByName(currentStructArgumentField).Set(
				convertArgument(currentStructArgumentFieldType, tokenValue),
			)
		case TOKEN_RIGHT_CURLY_BRACKET:
			structArgument = true
			arguments = append(arguments, currentStructArgument)
		}
	}

	fmt.Println(arguments, len(arguments))

	callMethod(methodToCall, arguments)
}

func convertArgument(argument reflect.Type, s string) reflect.Value {
	switch argument.Kind() {
	case reflect.Int:
		res, err := strconv.Atoi(s)
		if err != nil {
			panic("DSADS")
		}
		return reflect.ValueOf(res)
	// ...
	case reflect.Struct:

		s = strings.Trim(s, "\"")

		switch argument.Name() {
/*
		case "UUID":
			return reflect.ValueOf(uuid.FromString(s))
*/
		default:
			return reflect.ValueOf(s)
		}
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
