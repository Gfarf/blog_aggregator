package main

import (
	"fmt"

	"github.com/Gfarf/blog_aggregator/internal/config"
	"github.com/Gfarf/blog_aggregator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	commands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	function, ok := c.commands[cmd.name]
	if !ok {
		return fmt.Errorf("no valid command found")
	}
	return function(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	_, ok := c.commands[name]
	if ok {
		return
	}
	c.commands[name] = f
}
