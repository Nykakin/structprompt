package structprompt

import (
	"fmt"
	"testing"

	"github.com/c-bata/go-prompt"
)

type documentMock struct {
	s string
}

func (dm documentMock) TextBeforeCursor() string {
	return dm.s
}

type zag struct {
	S string
}

func (z zag) ZagFunc(int) int { return 0 }

type bar struct {
	Z          *string
	Zag        zag
	zagPrivate zag
	Zig        zag
	zigPrivate zag

	Soo1       string
	Soo2       string
	sooPrivate int
}

func (b bar) BarFunc1(i int) bool                { return true }
func (b bar) BarFunc2(c, d, e bool) (int, error) { return 0, nil }
func (b bar) SooFunc1() int64                    { return int64(100) }
func (b bar) SooFunc2()                          {}

type foo struct {
	Bar1       bar
	Bar2       bar
	barPrivate bar
	Foo1       bar
	Foo2       bar
	fooPrivate bar
}

func (f foo) FooFunc1()                         {}
func (f foo) FooFunc2(i int)                    {}
func (f foo) FooFunc3(b bool, s string)         {}
func (f foo) FooFunc4(i int) int                { return 0 }
func (f foo) FooFunc5(i int) (int, error)       { return 0, nil }
func (f foo) fooFuncPrivate(i int) (int, error) { return 0, nil }
func (f foo) BarFunc1()                         {}
func (f foo) BarFunc2(i int)                    {}
func (f foo) BarFunc3(b bool, s string)         {}
func (f foo) BarFunc4(i int) int                { return 0 }
func (f foo) BarFunc5(i int) (int, error)       { return 0, nil }
func (f foo) barFuncPrivate(i int) (int, error) { return 0, nil }

func TestPathCompletion(t *testing.T) {
	f := foo{}
	completer := completer{i: f}

	cases := []struct {
		str         string
		suggestions []prompt.Suggest
	}{
		{
			"nonsense",
			[]prompt.Suggest{},
		},
		{
			"B",
			[]prompt.Suggest{
				{Text: "Bar1"},
				{Text: "Bar2"},
				{Text: "BarFunc1()", Description: "()"},
				{Text: "BarFunc2()", Description: "(int)"},
				{Text: "BarFunc3()", Description: "(bool, string)"},
				{Text: "BarFunc4()", Description: "(int) int"},
				{Text: "BarFunc5()", Description: "(int) (int, error)"},
			},
		},
		{
			"Ba",
			[]prompt.Suggest{
				{Text: "Bar1"},
				{Text: "Bar2"},
				{Text: "BarFunc1()", Description: "()"},
				{Text: "BarFunc2()", Description: "(int)"},
				{Text: "BarFunc3()", Description: "(bool, string)"},
				{Text: "BarFunc4()", Description: "(int) int"},
				{Text: "BarFunc5()", Description: "(int) (int, error)"},
			},
		},
		{
			"Barf",
			[]prompt.Suggest{
				{Text: "BarFunc1()", Description: "()"},
				{Text: "BarFunc2()", Description: "(int)"},
				{Text: "BarFunc3()", Description: "(bool, string)"},
				{Text: "BarFunc4()", Description: "(int) int"},
				{Text: "BarFunc5()", Description: "(int) (int, error)"},
			},
		},
		{
			"BarFunc3",
			[]prompt.Suggest{
				{Text: "BarFunc3()", Description: "(bool, string)"},
			},
		},
		{
			"Bar1",
			[]prompt.Suggest{
				{Text: "Bar1"},
			},
		},
		{
			"Bar1.",
			[]prompt.Suggest{
				{Text: "Bar1.Z"},
				{Text: "Bar1.Zag"},
				{Text: "Bar1.Zig"},
				{Text: "Bar1.Soo1"},
				{Text: "Bar1.Soo2"},
				{Text: "Bar1.BarFunc1()", Description: "(int) bool"},
				{Text: "Bar1.BarFunc2()", Description: "(bool, bool, bool) (int, error)"},
				{Text: "Bar1.SooFunc1()", Description: "() int64"},
				{Text: "Bar1.SooFunc2()", Description: "()"},
			},
		},
		{
			"Bar1.Z",
			[]prompt.Suggest{
				{Text: "Bar1.Z"},
				{Text: "Bar1.Zag"},
				{Text: "Bar1.Zig"},
			},
		},
		{
			"Bar1.Zag",
			[]prompt.Suggest{
				{Text: "Bar1.Zag"},
			},
		},
		{
			"Bar1.Zag.",
			[]prompt.Suggest{
				{Text: "Bar1.Zag.S"},
				{Text: "Bar1.Zag.ZagFunc()", Description: "(int) int"},
			},
		},
	}

	var suggestions, expected string
	for _, c := range cases {
		suggestions = fmt.Sprintf("%v", completer.complete(documentMock{c.str}))
		expected = fmt.Sprintf("%v", c.suggestions)
		if suggestions != expected {
			t.Errorf("Expected %s, got %s", expected, suggestions)
		}
	}

}
