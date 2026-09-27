package middleware

import (
	"context"
	"fmt"

	"github.com/mohammad-rizwan-hussain/gator/internal/cmd"
	"github.com/mohammad-rizwan-hussain/gator/internal/config"
	"github.com/mohammad-rizwan-hussain/gator/internal/database"
)

func MiddlewareLoggedIn(handler func(s *config.State, cmd cmd.Command, user *database.User) error) func(*config.State, cmd.Command) error {
	return func(s *config.State, cmd cmd.Command) error {
		user, err := s.DB.GetUser(
			context.Background(),
			s.Config.CurrentUser,
		)
		if err != nil {
			return fmt.Errorf("user not logged in")
		}

		return handler(s, cmd, &user)
	}
}
