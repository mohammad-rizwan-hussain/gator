package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/mohammad-rizwan-hussain/gator/internal/cmd"
	"github.com/mohammad-rizwan-hussain/gator/internal/config"
	"github.com/mohammad-rizwan-hussain/gator/internal/database"
	"github.com/mohammad-rizwan-hussain/gator/internal/rss"
)

func HandlerLogin(s *config.State, cmd cmd.Command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("login excepts a username")
	}

	user := cmd.Arguments[0]

	// Get the users from the database
	_, err := s.DB.GetUser(context.Background(), user)
	if err != nil {
		return fmt.Errorf("user %s does not exist\nPlease register.", user)
	}

	err = s.Config.SetUser(user)
	if err != nil {
		return err
	}
	fmt.Printf("User: %s has been logged-in\n", user)
	return nil
}

func HandlerRegister(s *config.State, cmd cmd.Command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("register excepts a username")
	}

	userName := cmd.Arguments[0]

	// check if user already exists
	_, err := s.DB.GetUser(context.Background(), userName)
	if err == nil {
		return fmt.Errorf("user %s already exists", userName)
	}

	// set the config
	err = s.Config.SetUser(userName)
	if err != nil {
		return err
	}

	user := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      userName,
	}

	_, err = s.DB.CreateUser(context.Background(), user)
	if err != nil {
		return err
	}
	fmt.Printf("User: %s has been registered\n", userName)
	return nil
}

func HandlerGetUsers(s *config.State, cmd cmd.Command) error {
	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		return err
	}

	// Read the config for current user
	cfg, err := config.Read()
	currentUser := cfg.CurrentUser

	for _, user := range users {
		if user.Name == currentUser {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("- %s\n", user.Name)
		}
	}
	return nil
}

func HandlerResetDB(s *config.State, cmd cmd.Command) error {
	err := s.DB.ResetDB(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Database has been reset")
	return nil
}

func HandlerGetFeed(s *config.State, cmd cmd.Command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("add requries time_between_reqs")
	}
	rawPeriod := cmd.Arguments[0]
	duration, err := time.ParseDuration(rawPeriod)
	if err != nil {
		return fmt.Errorf("Unable to parse the passed duration: %w", rawPeriod)
	}

	fmt.Printf("Collecing feeds every %s\n", rawPeriod)

	ticker := time.NewTicker(duration)
	for ; ; <-ticker.C {
		rss.ScrapeFeeds(s)
	}
	return nil
}

func HandlerAddFeed(s *config.State, cmd cmd.Command, user *database.User) error {
	if len(cmd.Arguments) < 2 {
		return fmt.Errorf("addfeed requries name & url")
	}

	name := cmd.Arguments[0]
	url := cmd.Arguments[1]

	feed := database.AddFeedParams{
		ID:        uuid.New(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdFeed, err := s.DB.AddFeed(context.Background(), feed)
	if err != nil {
		return err
	}

	follow_feed := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    createdFeed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = s.DB.CreateFeedFollow(context.Background(), follow_feed)
	if err != nil {
		return err
	}

	fmt.Println(createdFeed)
	return nil
}

func HandlerListFeeds(s *config.State, cmd cmd.Command) error {
	feeds, err := s.DB.ListFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf("Feed name: %s\nURL: %s\nUser Name: %s\n\n", feed.FeedName, feed.Url, feed.UserName)
	}
	return nil
}

func HandlerFollowFeed(s *config.State, cmd cmd.Command, user *database.User) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("follow requries url")
	}

	url := cmd.Arguments[0]
	feed, err := s.DB.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}

	feed_follow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feed.FeedID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	followed_feed, err := s.DB.CreateFeedFollow(context.Background(), feed_follow)
	if err != nil {
		return err
	}
	fmt.Printf("User: %s following %s feed\n", followed_feed.UserName, followed_feed.FeedName)
	return nil
}

func HandlerFollowingFeed(s *config.State, cmd cmd.Command, user *database.User) error {
	feedsFollowedBycurrentUser, err := s.DB.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return err
	}
	for _, feed := range feedsFollowedBycurrentUser {
		fmt.Printf("- %s\n", feed.FeedName)
	}
	return nil
}

func HandlerUnFollowFeed(s *config.State, cmd cmd.Command, user *database.User) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("follow requries url")
	}

	url := cmd.Arguments[0]

	feed, err := s.DB.GetFeedByURL(context.Background(), url)
	if err != nil {
		return nil
	}

	UnfollowedFeed := database.DeleteFollowFeedParams{
		FeedID: feed.FeedID,
		UserID: user.ID,
	}

	err = s.DB.DeleteFollowFeed(context.Background(), UnfollowedFeed)
	if err != nil {
		return err
	}
	return nil
}

func HandlerBrowseFeeds(s *config.State, cmd cmd.Command, user *database.User) error {
	limit := 2
	if len(cmd.Arguments) > 0 {
		parsedLimit, err := strconv.Atoi(cmd.Arguments[0])
		if err != nil {
			return err
		}
		limit = parsedLimit

	}
	query := database.GetPostsForUserParams{
		Name:  user.Name,
		Limit: int32(limit),
	}
	posts, err := s.DB.GetPostsForUser(context.Background(), query)
	if err != nil {
		return err
	}
	for _, post := range posts {
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Printf("%s\n\n", post.Description)
	}
	return nil
}
