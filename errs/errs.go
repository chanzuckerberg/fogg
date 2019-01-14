package errs

import (
	"fmt"

	"github.com/go-errors/errors"
)

type User struct {
	e *errors.Error
}

type Internal struct {
	e *errors.Error
}

func (u *User) Error() string {
	if u != nil && u.e != nil {
		return u.e.Error()
	}
	return ""
}

func (u *Internal) Error() string {
	return u.e.ErrorStack()
}

// NewUser returns a new user error
func NewUser(msg string) error {
	return &User{
		e: errors.New(msg),
	}
}

func NewUserf(msg string, a ...interface{}) error {
	return &User{
		e: errors.New(fmt.Sprintf(msg, a...)),
	}
}

func WrapUser(e error, msg string) error {
	if e == nil {
		return nil
	}
	return &User{
		e: errors.WrapPrefix(e, msg, 1),
	}
}

func WrapUserf(e error, msg string, a ...interface{}) error {
	if e == nil {
		return nil
	}
	return &User{
		e: errors.WrapPrefix(e, fmt.Sprintf(msg, a...), 1),
	}
}

func NewInternal(msg string) error {
	return &Internal{
		e: errors.New(msg),
	}
}

func WrapInternal(e error, msg string) error {
	if e == nil {
		return nil
	}
	return &User{
		e: errors.WrapPrefix(e, msg, 1),
	}
}
