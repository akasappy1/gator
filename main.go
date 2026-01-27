package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/akasappy1/gator/internal/config"
	"github.com/akasappy1/gator/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	var cfg config.Config
	cfg, err := cfg.Read()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	db, err := sql.Open("postgres", cfg.DbURL)
	dbQueries := database.New(db)
	var st State
	st.cfgPtr = &cfg
	st.dbPtr = dbQueries
	var cmds Commands
	cmds.commandsList = make(map[string]func(*State, Command) error)

	cmds.register("login", handlerLogins)
	inputs := os.Args
	if len(inputs) < 2 {
		fmt.Printf("Error: At least two arguments expected (ie login and username.)")
		os.Exit(1)
	}
	var cmd Command
	cmd.name = inputs[1]
	cmd.args = inputs[2:]
	stPtr := &st

	err = cmds.run(stPtr, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
