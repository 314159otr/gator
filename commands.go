package main

import (
	"fmt"
	"errors"
	"github.com/314159otr/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("username is required")
	}
	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldnt set user: %w", err)
	}
	fmt.Printf("user %s has been set\n", cmd.args[0])
	return nil
}

func (c * commands) run(s *state, cmd command) error {
	f, ok := c.cmds[cmd.name]
	if !ok {
		return errors.New("command not found")
	}
	return f(s, cmd)
}

func (c * commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}
