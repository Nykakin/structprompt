package structprompt

import (
	"github.com/c-bata/go-prompt"
)

type structPrompt struct {
	Prompt *prompt.Prompt

	executor  executor
	completer completer
}

func NewStructPrompt(i interface{}) structPrompt {
	executor := executor{i: i}
	completer := completer{i: i}

	return structPrompt{
		Prompt: prompt.New(executor.execute, completer.complete),
	}
}

func (sp structPrompt) Execute(str string) {
	sp.executor.execute(str)
}

func (sp structPrompt) Run() {
	sp.Prompt.Run()
}
