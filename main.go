package main

import (
	"log"
	"os"
	"github.com/314159otr/gator/internal/config"
)

func main() {
	cfg, err := config.Read() 
	if err != nil {
		log.Fatalf("error reading file. Error: %s", err)
	}

	programState := &state{ cfg: &cfg }

	cmds := commands{ cmds: map[string]func(*state, command) error{} }
	cmds.register("login", handlerLogin)

	args := os.Args
	if len(args) < 2 {
		log.Fatalf("Usage: cli <command> [args...]", args)
	}

	err = cmds.run(programState, command{name:args[1], args:args[2:]})
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
}
