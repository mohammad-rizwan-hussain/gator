package cmd

import (
	"github.com/mohammad-rizwan-hussain/gator/internal/config"
)

type Command struct {
	Name      string
	Arguments []string
}

type Commands struct {
	Handlers map[string]func(*config.State, Command) error
}

func (cmds *Commands) Run(s *config.State, cmd Command) error {
	if f, ok := cmds.Handlers[cmd.Name]; ok {
		return f(s, cmd)
	}
	return nil
}

func (cmds *Commands) Register(name string, f func(*config.State, Command) error) {
	if cmds.Handlers == nil {
		cmds.Handlers = make(map[string]func(*config.State, Command) error)
	}
	cmds.Handlers[name] = f
}
