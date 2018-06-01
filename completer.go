package structprompt

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/c-bata/go-prompt"
)

type completer struct {
	i interface{}

	suggestions               []prompt.Suggest
	structArgumentSuggestions []prompt.Suggest

	fields []string
}

func (c *completer) complete(t prompt.Document) []prompt.Suggest {
	defer func() {
		if r := recover(); r != nil {
			c.suggestions = []prompt.Suggest{}
		}
	}()

	var token Token
	var tokenValue string

	var currentField string
	var methodToCall reflect.Value
	var methodToCallType reflect.Type

	currentValue := reflect.ValueOf(c.i)
	c.getFields(currentValue)
	c.fields = []string{}

	argumentCount := 0
	showStructArgumentPrompt := false

	lexer := NewLexer(t.TextBeforeCursor())

	for {
		token = lexer.NextToken()
		tokenValue = strings.TrimSpace(token.Value)

		if token.Type == TOKEN_EOF {
			break
		}

		switch token.Type {
		case TOKEN_DOT:
			currentValue = reflect.Indirect(currentValue.FieldByName(currentField))
			if currentValue.Kind() == reflect.Interface {
				currentValue = currentValue.Elem()
			}
			c.fields = append(c.fields, currentField)
			c.getFields(currentValue)
		case TOKEN_COMMA:
			if !showStructArgumentPrompt {
				argumentCount += 1
			}
		case TOKEN_PATH_ELEMENT:
			currentField = tokenValue
		case TOKEN_METHOD:
			methodToCall = currentValue.MethodByName(tokenValue)
			methodToCallType = methodToCall.Type()
		case TOKEN_LEFT_CURLY_BRACKET:
			c.structArgumentPrompt(methodToCallType.In(argumentCount))
			showStructArgumentPrompt = true
		case TOKEN_RIGHT_CURLY_BRACKET:
			showStructArgumentPrompt = false
			argumentCount += 1
		}
	}

	c.fields = []string{}
	p := strings.Split(t.TextBeforeCursor(), "(")

	if showStructArgumentPrompt {
		return c.structArgumentSuggestions
	}
	return prompt.FilterHasPrefix(c.suggestions, p[0], true)
}

func (c *completer) structArgumentPrompt(arg reflect.Type) {
	var name string

	c.structArgumentSuggestions = []prompt.Suggest{}
	for i := 0; i < arg.NumField(); i++ {
		name = arg.Field(i).Name
		if strings.ToUpper(name[0:1]) == name[0:1] {
			c.structArgumentSuggestions = append(c.structArgumentSuggestions, prompt.Suggest{
				Text:        arg.Field(i).Name,
				Description: arg.Field(i).Type.String(),
			})
		}
	}
}

func (c *completer) getFields(s reflect.Value) {
	c.suggestions = []prompt.Suggest{}
	var name string

	typeOfT := s.Type()
	prefix := strings.Join(c.fields, ".")
	if len(prefix) > 0 {
		prefix += "."
	}
	for i := 0; i < s.NumField(); i++ {
		name = typeOfT.Field(i).Name
		if strings.ToUpper(name[0:1]) == name[0:1] {
			c.suggestions = append(c.suggestions, prompt.Suggest{
				Text: prefix + name,
			})
		}
	}

	for i := 0; i < s.NumMethod(); i++ {
		name = typeOfT.Method(i).Name
		if strings.ToUpper(name[0:1]) == name[0:1] {
			c.suggestions = append(c.suggestions, prompt.Suggest{
				Text:        prefix + typeOfT.Method(i).Name + "()",
				Description: methodSignature(typeOfT.Method(i).Type),
			})
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

	switch len(returnTypes) {
	case 0:
		return fmt.Sprintf("(%s)", strings.Join(arguments[1:], ", "))
	case 1:
		return fmt.Sprintf("(%s) %s", strings.Join(arguments[1:], ", "), returnTypes[0])
	default:
		return fmt.Sprintf("(%s) (%s)", strings.Join(arguments[1:], ", "), strings.Join(returnTypes, ", "))
	}
}
