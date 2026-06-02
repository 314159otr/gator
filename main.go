package main

import (
	"log"
	"os"
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/314159otr/gator/internal/config"
	"github.com/314159otr/gator/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	cfg, err := config.Read() 
	if err != nil {
		log.Fatalf("error reading file. Error: %s", err)
	}

	
	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatalf("error reading file. Error: %s", err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	programState := &state{
		cfg: &cfg,
		db:  dbQueries,
	}

	cmds := commands{ cmds: map[string]func(*state, command) error{}, }
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)

	args := os.Args
	if len(args) < 2 {
		log.Fatalf("Usage: cli <command> [args...]", args)
	}

	err = cmds.run(programState, command{name:args[1], args:args[2:]})
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
}
