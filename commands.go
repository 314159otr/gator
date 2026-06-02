package main

import (
	"fmt"
	"errors"
	"time"
	"context"
	"database/sql"

	"github.com/google/uuid"

	"github.com/314159otr/gator/internal/database"
)

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
	username := cmd.args[0]

	ctx := context.Background()
	_, err := s.db.GetUser(ctx, username)	
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("user \"%s\" doesnt exist", username)
	}
	if err != nil {
		return fmt.Errorf("error getting the user: %w", err)
	}

	if err := s.cfg.SetUser(username); err != nil {
		return fmt.Errorf("couldnt set user: %w", err)
	}
	fmt.Printf("user %s has been set\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("username is required")
	}
	username := cmd.args[0]
	ctx := context.Background()
	userParams := database.CreateUserParams {
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}
	user, err := s.db.CreateUser(ctx, userParams)
	if err != nil {
		return fmt.Errorf("couldnt create user %s: %w", username, err)
	}
	if err := s.cfg.SetUser(user.Name); err != nil {
		return fmt.Errorf("couldnt set user: %w", err)
	}
	fmt.Println("user was created:")
	printUser(user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()
	err := s.db.DeleteUsers(ctx)
	if err != nil {
		return fmt.Errorf("couldnt delete users. Error: %w", err)
	}
	fmt.Println("all users deleted")
	return nil
}

func printUser(user database.User) {
	fmt.Printf("ID:        %v\n", user.ID)
	fmt.Printf("Name:      %v\n", user.Name)
	fmt.Printf("CreatedAt: %v\n", user.CreatedAt)
	fmt.Printf("UpdatedAt: %v\n", user.UpdatedAt)
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
