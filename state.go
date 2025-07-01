package main

import (
	"fmt"

	"github.com/Gfarf/blog_aggregator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	commands map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if cmd.args == nil {
		return fmt.Errorf("expected username, found nothing")
	}
	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Println("User has been set.")
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	function, ok := c.commands[cmd.name]
	if ok != true {
		return fmt.Errorf("no valid command found")
	}
	return function(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	_, ok := c.commands[name]
	if ok == true {
		return
	}
	c.commands[name] = f
}
