package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Gfarf/blog_aggregator/internal/database"
	"github.com/google/uuid"
)

func handleUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("expected url")
	}
	feed, err := s.db.GetFeedsFromURL(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}
	err = s.db.Unfollow(context.Background(), database.UnfollowParams{UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		return err
	}
	fmt.Println("Feed unfollowed")
	return nil
}

func handleFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("expected url")
	}
	feed, err := s.db.GetFeedsFromURL(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}
	ffollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}
	fmt.Println(ffollow)
	return nil
}

func handleFollowing(s *state, cmd command, user database.User) error {
	list, err := s.db.GetFeedFollowForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	for _, item := range list {
		fmt.Println(item.FeedName)
	}
	return nil
}

func handleGetFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf("Feed name: %s\n", feed.Name)
		fmt.Printf("Feed url: %s\n", feed.Url)
		userName, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("User name: %s\n", userName)
	}
	return nil
}

func handleAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("expected name and Url")
	}
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID})
	if err != nil {
		return err
	}
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}
	fmt.Println(feed)
	return nil
}

func handleFetch(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("expected time between requisitions")
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %s\n", cmd.args[0])
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("expected username")
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
	if len(cmd.args) != 1 {
		return fmt.Errorf("expected name")
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
		if s.cfg.CurrentUserName == name.Name {
			fmt.Printf(" (current)\n")
		} else {
			fmt.Println()
		}
	}
	return nil
}

func handleFeedFollows(s *state, cmd command) error {
	list, err := s.db.GetFeedFollows(context.Background())
	if err != nil {
		return err
	}
	fmt.Println(list)
	for _, name := range list {
		fmt.Println(name)
	}
	return nil
}
