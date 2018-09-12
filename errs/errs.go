package errs

import (
	"github.com/go-errors/errors"
)

type User struct {
	e *errors.Error
}

type Internal struct {
	e *errors.Error
}

func (u *User) Error() string {
	return u.e.Error()
}

func (u *Internal) Error() string {
	return u.e.ErrorStack()
}

func NewUser(msg string) *User {
	return &User{
		e: errors.New(msg),
	}
}

func WrapUser(e error, msg string) *User {
	return &User{
		e: errors.WrapPrefix(e, msg, 1),
	}
}

func WrapInternal(e error, msg string) *User {
	return &User{
		e: errors.WrapPrefix(e, msg, 1),
	}
}
