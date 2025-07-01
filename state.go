package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Gfarf/blog_aggregator/internal/config"
	"github.com/Gfarf/blog_aggregator/internal/database"
	"github.com/google/uuid"
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

func handlerLogin(s *state, cmd command) error {
	if cmd.args == nil {
		return fmt.Errorf("expected username, found nothing")
	}
	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		fmt.Println("user not found in database")
		os.Exit(1)
	}
	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Println("User has been set.")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if cmd.args == nil {
		return fmt.Errorf("expected name, found nothing")
	}
	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err == nil {
		fmt.Println("user with given name already exists.")
		os.Exit(1)
	}

	params := database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.args[0]}
	user, err := s.db.CreateUser(context.Background(), params)

	if err != nil {
		return err
	}
	s.cfg.SetUser(user.Name)

	fmt.Println("User has been included in database.")
	fmt.Println(user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("reset failed: %s", err)
	}
	return nil
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
