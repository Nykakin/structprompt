package structprompt

import (
	"testing"
)

func TestLexer(t *testing.T) {
	cases := []struct {
		str    string
		tokens []Token
	}{
		{
			"Foo",
			[]Token{
				{TOKEN_EOF, ""},
			},
		},
		{
			"Foo.",
			[]Token{
				{TOKEN_PATH_ELEMENT, "Foo"},
				{TOKEN_DOT, "."},
				{TOKEN_EOF, ""},
			},
		},
		{
			"Foo.Bar",
			[]Token{
				{TOKEN_PATH_ELEMENT, "Foo"},
				{TOKEN_DOT, "."},
				{TOKEN_EOF, ""},
			},
		},
		{
			"Foo.Bar.",
			[]Token{
				{TOKEN_PATH_ELEMENT, "Foo"},
				{TOKEN_DOT, "."},
				{TOKEN_PATH_ELEMENT, "Bar"},
				{TOKEN_DOT, "."},
				{TOKEN_EOF, ""},
			},
		},
		{
			"Foo.Bar.Fun(",
			[]Token{
				{TOKEN_PATH_ELEMENT, "Foo"},
				{TOKEN_DOT, "."},
				{TOKEN_PATH_ELEMENT, "Bar"},
				{TOKEN_DOT, "."},
				{TOKEN_METHOD, "Fun"},
				{TOKEN_LEFT_BRACKET, "("},
				{TOKEN_EOF, ""},
			},
		},
		{
			"Foo.Bar.Fun()",
			[]Token{
				{TOKEN_PATH_ELEMENT, "Foo"},
				{TOKEN_DOT, "."},
				{TOKEN_PATH_ELEMENT, "Bar"},
				{TOKEN_DOT, "."},
				{TOKEN_METHOD, "Fun"},
				{TOKEN_LEFT_BRACKET, "("},
				{TOKEN_RIGHT_BRACKET, ")"},
				{TOKEN_EOF, ""},
			},
		},
		{
			"Foo.Bar.Fun(12)",
			[]Token{
				{TOKEN_PATH_ELEMENT, "Foo"},
				{TOKEN_DOT, "."},
				{TOKEN_PATH_ELEMENT, "Bar"},
				{TOKEN_DOT, "."},
				{TOKEN_METHOD, "Fun"},
				{TOKEN_LEFT_BRACKET, "("},
				{TOKEN_METHOD_ARGUMENT, "12"},
				{TOKEN_RIGHT_BRACKET, ")"},
				{TOKEN_EOF, ""},
			},
		},
		{
			"Foo.Bar.Fun(12, \"bb\", 33)",
			[]Token{
				{TOKEN_PATH_ELEMENT, "Foo"},
				{TOKEN_DOT, "."},
				{TOKEN_PATH_ELEMENT, "Bar"},
				{TOKEN_DOT, "."},
				{TOKEN_METHOD, "Fun"},
				{TOKEN_LEFT_BRACKET, "("},
				{TOKEN_METHOD_ARGUMENT, "12"},
				{TOKEN_COMMA, ","},
				{TOKEN_METHOD_ARGUMENT, "\"bb\""},
				{TOKEN_COMMA, ","},
				{TOKEN_METHOD_ARGUMENT, "33"},
				{TOKEN_RIGHT_BRACKET, ")"},
				{TOKEN_EOF, ""},
			},
		},
		{
			"Foo.Bar.Fun({Foo: 23})",
			[]Token{
				{TOKEN_PATH_ELEMENT, "Foo"},
				{TOKEN_DOT, "."},
				{TOKEN_PATH_ELEMENT, "Bar"},
				{TOKEN_DOT, "."},
				{TOKEN_METHOD, "Fun"},
				{TOKEN_LEFT_BRACKET, "("},
				{TOKEN_LEFT_CURLY_BRACKET, "{"},
				{TOKEN_FIELD, "Foo"},
				{TOKEN_COLON, ":"},
				{TOKEN_FIELD_VALUE, "23"},
				{TOKEN_RIGHT_CURLY_BRACKET, "}"},
				{TOKEN_METHOD_ARGUMENT, ""},
				{TOKEN_RIGHT_BRACKET, ")"},
				{TOKEN_EOF, ""},
			},
		},
		{
			"Foo.Bar.Fun({Foo: \"FooV\", Bar: \"BarV\"})",
			[]Token{
				{TOKEN_PATH_ELEMENT, "Foo"},
				{TOKEN_DOT, "."},
				{TOKEN_PATH_ELEMENT, "Bar"},
				{TOKEN_DOT, "."},
				{TOKEN_METHOD, "Fun"},
				{TOKEN_LEFT_BRACKET, "("},
				{TOKEN_LEFT_CURLY_BRACKET, "{"},
				{TOKEN_FIELD, "Foo"},
				{TOKEN_COLON, ":"},
				{TOKEN_FIELD_VALUE, "\"FooV\""},
				{TOKEN_COMMA, ","},
				{TOKEN_FIELD, "Bar"},
				{TOKEN_COLON, ":"},
				{TOKEN_FIELD_VALUE, "\"BarV\""},
				{TOKEN_RIGHT_CURLY_BRACKET, "}"},
				{TOKEN_METHOD_ARGUMENT, ""},
				{TOKEN_RIGHT_BRACKET, ")"},
				{TOKEN_EOF, ""},
			},
		},
		{
			"Foo.Bar.Fun('a', {B: 100}, 0.43)",
			[]Token{
				{TOKEN_PATH_ELEMENT, "Foo"},
				{TOKEN_DOT, "."},
				{TOKEN_PATH_ELEMENT, "Bar"},
				{TOKEN_DOT, "."},
				{TOKEN_METHOD, "Fun"},
				{TOKEN_LEFT_BRACKET, "("},
				{TOKEN_METHOD_ARGUMENT, "'a'"},
				{TOKEN_COMMA, ","},
				{TOKEN_LEFT_CURLY_BRACKET, "{"},
				{TOKEN_FIELD, "B"},
				{TOKEN_COLON, ":"},
				{TOKEN_FIELD_VALUE, "100"},
				{TOKEN_RIGHT_CURLY_BRACKET, "}"},
				{TOKEN_METHOD_ARGUMENT, ""},
				{TOKEN_COMMA, ","},
				{TOKEN_METHOD_ARGUMENT, "0.43"},
				{TOKEN_RIGHT_BRACKET, ")"},
				{TOKEN_EOF, ""},
			},
		},
	}

	var token Token
	for _, c := range cases {
		lexer := NewLexer(c.str)
		for _, expected := range c.tokens {
			token = lexer.NextToken()
			if token.Type != expected.Type || token.Value != expected.Value {
				t.Errorf("Expected %+v, got %+v", expected, token)
			}
		}
	}
}
