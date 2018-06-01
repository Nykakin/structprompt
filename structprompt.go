package structprompt

import (
	"github.com/c-bata/go-prompt"
)

type structPrompt struct {
	Prompt *prompt.Prompt
}

func NewStructPrompt(i interface{}) structPrompt {
	executor := executor{i: i}
	completer := completer{i: i}

	return structPrompt{
		Prompt: prompt.New(executor.execute, completer.complete),
	}
}

func (s structPrompt) Run() {
	s.Prompt.Run()
}
