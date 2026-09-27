package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"github.com/mohammad-rizwan-hussain/gator/internal/cmd"
	"github.com/mohammad-rizwan-hussain/gator/internal/config"
	"github.com/mohammad-rizwan-hussain/gator/internal/database"
	"github.com/mohammad-rizwan-hussain/gator/internal/handlers"
	"github.com/mohammad-rizwan-hussain/gator/internal/middleware"
)

func main() {
	// Read the command aurguments
	args := os.Args

	commandName := args[1]
	commandArgs := args[2:]

	// Read the config file
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("error reading config: %v\n", err)
		panic("error reading config file")
	}

	// Initialize the database connection
	db, err := sql.Open("postgres", cfg.DB_URL)
	if err != nil {
		fmt.Printf("error connecting to database: %v\n", err)
		panic("error connecting to database")
	}
	defer db.Close()

	dbQueries := database.New(db)

	// Initialize the state
	state := &config.State{
		DB:     dbQueries,
		Config: cfg,
	}

	// Register the commands
	UserCmd := cmd.Command{
		Name:      commandName,
		Arguments: commandArgs,
	}

	// Initialize the command handler
	cmds := cmd.Commands{}
	cmds.Register("login", handlers.HandlerLogin)
	cmds.Register("register", handlers.HandlerRegister)
	cmds.Register("reset", handlers.HandlerResetDB)
	cmds.Register("users", handlers.HandlerGetUsers)
	cmds.Register("agg", handlers.HandlerGetFeed)
	cmds.Register("addfeed", middleware.MiddlewareLoggedIn(handlers.HandlerAddFeed))
	cmds.Register("feeds", handlers.HandlerListFeeds)
	cmds.Register("follow", middleware.MiddlewareLoggedIn(handlers.HandlerFollowFeed))
	cmds.Register("following", middleware.MiddlewareLoggedIn(handlers.HandlerFollowingFeed))
	cmds.Register("unfollow", middleware.MiddlewareLoggedIn(handlers.HandlerUnFollowFeed))
	cmds.Register("browse", middleware.MiddlewareLoggedIn(handlers.HandlerBrowseFeeds))

	// Run the command
	err = cmds.Run(state, UserCmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
