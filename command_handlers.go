package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/akasappy1/gator/internal/config"
	"github.com/akasappy1/gator/internal/database"
	"github.com/google/uuid"
)

type State struct {
	dbPtr  *database.Queries
	cfgPtr *config.Config
}

type Command struct {
	name string
	args []string
}

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.args) == 0 {
		err := errors.New("No login command or username provided.")
		return err
	}
	_, err := s.dbPtr.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		fmt.Printf("Cannot login as username does not exist. Try register instead.\n")
		os.Exit(1)
	}
	newUser := cmd.args[0]
	if err := s.cfgPtr.SetUser(newUser); err != nil {
		return err
	}
	fmt.Printf("User has been set.")
	return nil
}

func handlerRegister(s *State, cmd Command) error {
	if len(cmd.args) == 0 {
		err := errors.New("No username to register provided.")
		return err
	}
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}
	userCheck, err := s.dbPtr.GetUser(context.Background(), params.Name)
	if err == nil {
		if userCheck.Name == params.Name {
			fmt.Printf("Error: Cannot register username that already exists.\n")
			os.Exit(1)
		}
	}

	_, err = s.dbPtr.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}
	s.cfgPtr.CurrentUserName = params.Name

	newUser, err := s.dbPtr.GetUser(context.Background(), params.Name)
	if err != nil {
		fmt.Println(err)
		return err
	}
	s.cfgPtr.SetUser(newUser.Name)
	fmt.Printf("New User has been registered and set as current user.\n")
	fmt.Println(newUser)
	return nil
}

type Commands struct {
	commandsList map[string]func(*State, Command) error
}

func (c *Commands) run(s *State, cmd Command) error {
	command, ok := c.commandsList[cmd.name]
	if !ok {
		err := errors.New("Command not found.")
		return err
	}
	if err := command(s, cmd); err != nil {
		return err
	}
	// command(s, cmd)
	return nil
}

func (c *Commands) register(name string, f func(*State, Command) error) {
	c.commandsList[name] = f
}

func handlerReset(s *State, cmd Command) error {
	if err := s.dbPtr.DeleteAll(context.Background()); err != nil {
		fmt.Printf("Error when deleting users.\n")
		os.Exit(1)
	}
	fmt.Printf("Successfully deleted all users.\n")
	return nil
}

func handlerUsers(s *State, cmd Command) error {
	users, err := s.dbPtr.GetUsers(context.Background())
	if err != nil {
		fmt.Printf("Error generating user list")
		return err
	}
	for _, user := range users {
		if user == s.cfgPtr.CurrentUserName {
			fmt.Printf("%s (current)\n", user)
		} else {
			fmt.Println(user)
		}
	}
	return nil
}

func handlerAgg(s *State, cmd Command) error {
	ctx := context.Background()
	feedUrl := "https://www.wagslane.dev/index.xml"
	outFeed, err := fetchFeed(ctx, feedUrl)
	if err != nil {
		fmt.Printf("Error fetching RSS feed.")
		os.Exit(1)
	}
	fmt.Println(outFeed)
	return nil
}

func handlerAddfeed(s *State, cmd Command) error {
	if len(cmd.args) != 2 {
		fmt.Errorf("Error: command structure is addfeed <name> <url>.\n")
		os.Exit(1)
	}
	username := s.cfgPtr.CurrentUserName
	user, err := s.dbPtr.GetUser(context.Background(), username)
	if err != nil {
		return err
	}
	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}
	addedFeed, err := s.dbPtr.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}
	fmt.Println(addedFeed)
	return nil
}

func handlerFeeds(s *State, cmd Command) error {
	feedsList, err := s.dbPtr.GetFeeds(context.Background())
	if err != nil {
		fmt.Errorf("Error getting a list of feeds.")
		return err
	}
	for _, feed := range feedsList {
		fmt.Println(feed.Name)
		fmt.Println(feed.Url)
		fmt.Println(feed.Name_2)
	}
	return nil
}
