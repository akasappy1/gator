package main

import (
	"errors"
	"fmt"

	"github.com/akasappy1/gator/internal/config"
	"github.com/akasappy1/gator/internal/database"
)

type State struct {
	dbPtr  *database.Queries
	cfgPtr *config.Config
}

type Command struct {
	name string
	args []string
}

func handlerLogins(s *State, cmd Command) error {
	if len(cmd.args) == 0 {
		err := errors.New("No login command or username provided.")
		return err
	}
	newUser := cmd.args[0]
	if err := s.cfgPtr.SetUser(newUser); err != nil {
		return err
	}
	fmt.Printf("User has been set.")
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
	command(s, cmd)
	return nil
}

func (c *Commands) register(name string, f func(*State, Command) error) {
	c.commandsList[name] = f
}
