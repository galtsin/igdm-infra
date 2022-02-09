package instagram_api

import (
	"errors"
	"fmt"

	"channels-instagram-dm/domain"
)

var (
	ErrorNoLoggedIn           = errors.New("No logged in")
	ErrorLoginFailed          = errors.New("Login failed")
	ErrorChallengeFailed      = errors.New("Challenge failed")
	ErrorChallengeRequired    = errors.New("Challenge required")
	ErrorLoginInvalidUsername = errors.New("Invalid username")
	ErrorLoginBadPassword     = errors.New("Bad password")
)

var (
	errNoLoggedIn        = "User not logged in. Please call login() and then try again."
	errChallengeRequired = "challenge required"
	errBadPassword       = "bad password"
	errInvalidUser       = "invalid username"
)

func newError(text string) error {
	switch text {
	case errNoLoggedIn:
		return newErrorNoLoggedIn()
	case errChallengeRequired:
		return newErrorChallengeRequired()
	case errInvalidUser:
		return newErrorLoginInvalidUsername()
	case errBadPassword:
		return newErrorLoginBadPassword()
	default:
		return fmt.Errorf("%v", text)
	}
}

func newErrorNoLoggedIn() error {
	return domain.NewError(ErrorNoLoggedIn.Error(), fmt.Errorf("%w", domain.ErrorNoLoggedIn))
}

func newErrorChallengeFailed(text string) error {
	return domain.NewError(ErrorChallengeFailed.Error(), fmt.Errorf("%w. %s", domain.ErrorChallengeFailed, text))
}

func newErrorChallengeRequired() error {
	return domain.NewError(ErrorChallengeRequired.Error(), fmt.Errorf("%w", domain.ErrorInvalidCredentials))
}

func newErrorLoginFailed(text string) error {
	return domain.NewError(ErrorLoginFailed.Error(), fmt.Errorf("%w. %s", domain.ErrorLoginFailed, text))
}

func newErrorLoginInvalidUsername() error {
	return domain.NewError(ErrorLoginInvalidUsername.Error(), fmt.Errorf("%w", domain.ErrorInvalidCredentials))
}

func newErrorLoginBadPassword() error {
	return domain.NewError(ErrorLoginBadPassword.Error(), fmt.Errorf("%w", domain.ErrorInvalidCredentials))
}
