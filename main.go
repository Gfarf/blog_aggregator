package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Gfarf/blog_aggregator/internal/config"
)

func main() {

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

	cmds := commands{}
	cmds.commands = make(map[string]func(*state, command) error)

	cmds.register("login", handlerLogin)

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
