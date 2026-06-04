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

func handlerUsers(s *state, cmd command) error {
	ctx := context.Background()
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("couldnt get users. Error: %w", err)
	}
	fmt.Println("all users:")
	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
		fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	url :=  "https://www.wagslane.dev/index.xml"
	ctx := context.Background()
	rssFeed, err := fetchFeed(ctx, url)
	if err != nil {
		return fmt.Errorf("couldnt fetch URL: %v. Error: %w", url,  err)
	}
	fmt.Printf("feed: %+v\n", rssFeed)
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("Usage: addfeed <name> <url>")
	}

	ctx := context.Background()

	name := cmd.args[0]
	url := cmd.args[1]
	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
	}
	feed, err := s.db.CreateFeed(ctx, feedParams)
	if err != nil {
		return fmt.Errorf("error creating feed [%s - %s]: %w", name, url, err)
	}
	fmt.Println("feed was created:")
	printFeed(feed)

	feedFollowRow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("error creating the feed_follow: %w", err)
	}

	fmt.Println("feed follow was created:")
	printFeedFollow(feedFollowRow.UserName, feedFollowRow.FeedName)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	ctx := context.Background()
	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("couldnt get feeds. Error: %w", err)
	}
	fmt.Println("all feeds:")
	for _, feed := range feeds {
		user, err := s.db.GetUserById(ctx, feed.UserID)
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user id: \"%s\" doesnt exist", feed.UserID)
		}
		if err != nil {
			return fmt.Errorf("error getting the user: %w", err)
		}
		fmt.Println("===========")
		printFeed(feed)
		fmt.Println("User:")
		printUser(user)
		fmt.Println("===========")
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("Usage: follow <url>")
	}
	url := cmd.args[0]

	ctx := context.Background()

	printUser(user)

	feed, err := s.db.GetFeedByURL(ctx, url)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("feed \"%s\" doesnt exist", url)
	}
	if err != nil {
		return fmt.Errorf("error getting the feed: %w", err)
	}
	printFeed(feed)

	feedFollowRow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("error creating the feed_follow: %w", err)
	}
	fmt.Println("feed follow was created:")
	printFeedFollow(feedFollowRow.UserName, feedFollowRow.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	ctx := context.Background()
	feedFollowRows, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("error getting the feeds for %v: %w", user.Name, err)
	}
	fmt.Printf("feed follow for user %s:\n", user.Name)
	for _, feedFollowRow := range feedFollowRows {
		fmt.Printf("- %s\n", feedFollowRow.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("Usage: unfollow <url>")
	}
	url := cmd.args[0]

	ctx := context.Background()

	feed, err := s.db.GetFeedByURL(ctx, url)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("feed \"%s\" doesnt exist", url)
	}
	if err != nil {
		return fmt.Errorf("error getting the feed: %w", err)
	}

	deleteFeedFollowParams := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	if err := s.db.DeleteFeedFollow(ctx, deleteFeedFollowParams); err != nil {
		return fmt.Errorf("error deleting the follow feed: %w", err)
	}
	return nil

}

func middlewareLoggedIn(handler func(*state, command, database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		ctx := context.Background()
		user, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user \"%s\" doesnt exist", s.cfg.CurrentUserName)
		}
		if err != nil {
			return fmt.Errorf("error getting the user: %w", err)
		}
		return handler(s, cmd, user)
	}
}

func printUser(user database.User) {
	fmt.Printf("ID:        %v\n", user.ID)
	fmt.Printf("Name:      %v\n", user.Name)
	fmt.Printf("CreatedAt: %v\n", user.CreatedAt)
	fmt.Printf("UpdatedAt: %v\n", user.UpdatedAt)
}

func printFeed(feed database.Feed, ) {
	fmt.Printf("ID:        %v\n", feed.ID)
	fmt.Printf("UserID:    %v\n", feed.UserID)
	fmt.Printf("Name:      %v\n", feed.Name)
	fmt.Printf("Url:       %v\n", feed.Url)
	fmt.Printf("CreatedAt: %v\n", feed.CreatedAt)
	fmt.Printf("UpdatedAt: %v\n", feed.UpdatedAt)
}

func printFeedFollow(username string, feedname string) {
	fmt.Printf("User: %v\n", username)
	fmt.Printf("Feed: %v\n", feedname)
}
