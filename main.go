package main

import (
    "fmt"
	"github.com/c-bata/go-prompt"
)

type Social struct {
}

type ItemCommandHandler struct {
}

type Item struct {
    CommandHandler ItemCommandHandler
}

type MediaCommandHandler struct {
}

type Media struct {
    CommandHandler MediaCommandHandler
}

type CollectionCommandHandler struct {
}

type Collection struct {
    CommandHandler CollectionCommandHandler
}

type Assets struct {
    Item *Item
    Media Media
    Collection Collection
}

func (a Assets) Foo(i, j int) int {
    fmt.Println(i, j)
    return 100
}

func (a Assets) Bar(i int, s string) (int, error) {
    for j := 0; j < i; j++ {
        fmt.Println(s)
    }
    return -100, nil
}

func (a Assets) Zag(s1, s2 string) ([]int, error) {
    fmt.Println(s1 + " -- " + s2)
    return []int{666, 2323, 32, 32, 232213, 323232, 32323232, 323232, 323232, 323232, 3232323, 3232}, nil
}

type BooArgument struct {
    A int
    B int
    C string
}

func (a Assets) Boo(arg BooArgument) {
    fmt.Printf("%+v\n", arg)
}

type Accounts struct {
}

type Modules struct {
    Assets Assets
    Accounts Accounts
    Social Social
}

func main() {
    m = Modules{}
    m.Assets.Item = &Item{}

/*
    currentValue = reflect.ValueOf(&m).Elem()
    getFields(currentValue)

    currentValue = reflect.Indirect(currentValue.FieldByName("Assets"))
    getFields(currentValue)

    currentValue = reflect.Indirect(currentValue.FieldByName("Item"))
    getFields(currentValue)
*/

	p := prompt.New(
		executor,
		completer,
	)
	p.Run()


//    executor("Assets.Boo({B: 12})")
//    executor("Assets.Bar(12, \"Pies\")")

//    parse("Assets.Item.")










/*
    var token Token

    lexer := NewLexer("Assets.Boo({A: 12})")

    for {
        token = lexer.NextToken()

        if token.Type == TOKEN_EOF {
            break
        }

        switch token.Type {
        case TOKEN_DOT:
            fmt.Printf("TOKEN_DOT %+v\n", token.Value)
        case TOKEN_PATH_ELEMENT:
            fmt.Printf("TOKEN_PATH_ELEMENT %+v\n", token.Value)
        case TOKEN_METHOD:
            fmt.Printf("TOKEN_METHOD %+v\n", token.Value)
        case TOKEN_METHOD_ARGUMENT:
            fmt.Printf("TOKEN_METHOD_ARGUMENT %+v\n", token.Value)
        case TOKEN_LEFT_CURLY_BRACKET:
            fmt.Printf("TOKEN_LEFT_CURLY_BRACKET %+v\n", token.Value)
        case TOKEN_FIELD:
            fmt.Printf("TOKEN_FIELD %+v\n", token.Value)
        case TOKEN_FIELD_VALUE:
            fmt.Printf("TOKEN_FIELD_VALUE %+v\n", token.Value)
        case TOKEN_RIGHT_CURLY_BRACKET:
            fmt.Printf("TOKEN_RIGHT_CURLY_BRACKET %+v\n", token.Value)
        }
    }
*/

}
