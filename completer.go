package structprompt

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/c-bata/go-prompt"
)

type completer struct {
	i interface{}

	suggestions []prompt.Suggest
	fields      []string
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

	currentValue := reflect.ValueOf(c.i)
	fmt.Printf("%v\n", c.suggestions)
	c.getFields(currentValue)
	fmt.Printf("%v\n", c.suggestions)

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
			c.fields = append(c.fields, currentField)
			c.getFields(currentValue)
		case TOKEN_PATH_ELEMENT:
			currentField = tokenValue
		}
	}

	p := strings.Split(t.TextBeforeCursor(), "(")

	return prompt.FilterHasPrefix(c.suggestions, p[0], true)
}

func (c *completer) getFields(s reflect.Value) {
	c.suggestions = make([]prompt.Suggest, s.NumField()+s.NumMethod())

	typeOfT := s.Type()
	prefix := strings.Join(c.fields, ".")
	if len(prefix) > 0 {
		prefix += "."
	}
	for i := 0; i < s.NumField(); i++ {
		c.suggestions[i] = prompt.Suggest{
			Text: prefix + typeOfT.Field(i).Name,
			/*            Description: typeOfT.Field(i).Name,*/
		}
	}

	for i := 0; i < s.NumMethod(); i++ {
		c.suggestions[s.NumField()+i] = prompt.Suggest{
			Text:        prefix + typeOfT.Method(i).Name + "()",
			Description: methodSignature(typeOfT.Method(i).Type),
		}
	}

	fmt.Printf("%v\n", c.suggestions)
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
