package errs

import (
	"fmt"

	"github.com/go-errors/errors"
	log "github.com/sirupsen/logrus"
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

func NewUser(msg string) *User {
	return &User{
		e: errors.New(msg),
	}
}

func NewUserf(msg string, a ...interface{}) *User {
	return &User{
		e: errors.New(fmt.Sprintf(msg, a...)),
	}
}

func WrapUser(e error, msg string) *User {
	if e == nil {
		return nil
	}
	return &User{
		e: errors.WrapPrefix(e, msg, 1),
	}
}

func WrapUserf(e error, msg string, a ...interface{}) *User {
	log.Debugf("wrapuser e: %#v", e)
	if e == nil {
		return nil
	}
	return &User{
		e: errors.WrapPrefix(e, fmt.Sprintf(msg, a...), 1),
	}
}

func NewInternal(msg string) *Internal {
	return &Internal{
		e: errors.New(msg),
	}
}

func WrapInternal(e error, msg string) *User {
	if e == nil {
		return nil
	}
	return &User{
		e: errors.WrapPrefix(e, msg, 1),
	}
}
