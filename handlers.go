package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Gfarf/blog_aggregator/internal/database"
	"github.com/google/uuid"
)

func handleFetch(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Feed: %+v\n", feed)
	return nil
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

func handlerUsers(s *state, cmd command) error {
	list, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, name := range list {
		fmt.Printf("* %s", name)
		if s.cfg.CurrentUserName == name {
			fmt.Printf(" (current)\n")
		} else {
			fmt.Println()
		}
	}
	return nil
}
