package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Gfarf/blog_aggregator/internal/config"
	"github.com/Gfarf/blog_aggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	// "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" -- connection string

	// Obtém os argumentos da linha de comando.
	args := os.Args

	// Verifica se existem mais argumentos.
	if len(args) < 2 {
		fmt.Println("Not enough arguments passed in the command line")
		os.Exit(1)
	}

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	state1 := state{}
	state1.cfg = &cfg

	db, err := sql.Open("postgres", state1.cfg.DdbUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()
	dbQueries := database.New(db)
	state1.db = dbQueries

	cmds := commands{}
	cmds.commands = make(map[string]func(*state, command) error)

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handleFetch)
	cmds.register("addfeed", middlewareLoggedIn(handleAddFeed))
	cmds.register("feeds", handleGetFeeds)
	cmds.register("follow", middlewareLoggedIn(handleFollow))
	cmds.register("following", middlewareLoggedIn(handleFollowing))
	cmds.register("feedfollows", handleFeedFollows)
	cmds.register("unfollow", middlewareLoggedIn(handleUnfollow))
	cmds.register("brownse", middlewareLoggedIn(handleGetPosts))
	cmds.register("tenposts", allPosts)

	cmd := command{}
	cmd.name = args[1]
	if len(args) > 2 {
		cmd.args = args[2:]
	}
	ok := cmds.run(&state1, cmd)
	if ok != nil {
		fmt.Println(ok)
		os.Exit(1)
	}

}
