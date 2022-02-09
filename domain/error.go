package domain

import (
	"errors"
	"fmt"
)

var (
	ErrorNotFound           = errors.New("Element not found")
	ErrorNoLoggedIn         = errors.New("User not logged in")
	ErrorLoginFailed        = errors.New("Login failed")
	ErrorChallengeFailed    = errors.New("Challenge failed")
	ErrorInvalidCredentials = errors.New("Invalid credentials")
	ErrorInvalidArgument    = errors.New("Invalid argument")
	ErrorPermissionDenied   = errors.New("Permission denied")
)

type BaseError interface {
	error
	Code() string
}

type Error struct {
	err  error
	code string
}

func (e *Error) Error() string {
	return e.err.Error()
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) Code() string {
	return e.code
}

func NewError(code string, err error) error {
	return &Error{
		code: code,
		err:  err,
	}
}

func NewErrorNotFound(msg string) error {
	return NewError(ErrorNotFound.Error(), fmt.Errorf("%w. %s", ErrorNotFound, msg))
}

func NewErrorNoLoggedIn(msg string) error {
	return NewError(ErrorNoLoggedIn.Error(), fmt.Errorf("%w. %s", ErrorNoLoggedIn, msg))
}

func NewErrorInvalidArgument(msg string) error {
	return NewError(ErrorInvalidArgument.Error(), fmt.Errorf("%w. %s", ErrorInvalidArgument, msg))
}
