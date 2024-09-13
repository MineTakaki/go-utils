package errors

import (
	goerr "errors"
	"fmt"

	"github.com/pkg/errors"
)

var fnNew func(string) error
var fnWithStack func(error) error
var fnWrap func(error, string) error

func init() {
	SetNew(errors.New)
	SetWithStack(errors.WithStack)
	SetWrap(errors.Wrap)
}

func SetNew(fn func(string) error) {
	fnNew = fn
}

func SetWithStack(fn func(error) error) {
	fnWithStack = fn
}

func SetWrap(fn func(error, string) error) {
	fnWrap = fn
}

func New(text string) error {
	if fn := fnNew; fn != nil {
		return fn(text)
	}
	return goerr.New(text)
}

func Errorf(format string, args ...interface{}) error {
	text := fmt.Sprintf(format, args...)
	if fn := fnNew; fn != nil {
		return fn(text)
	}
	return goerr.New(text)
}

func Wrap(err error, text string) error {
	if err == nil {
		return err
	}
	w := fnWrap
	if w == nil {
		return errors.Wrap(err, text)
	}
	return w(err, text)
}

func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return err
	}
	text := fmt.Sprintf(format, args...)
	w := fnWrap
	if w == nil {
		return errors.Wrap(err, text)
	}
	return w(err, text)
}

func WithStack(err error) error {
	if err == nil {
		return err
	}
	w := fnWithStack
	if w == nil {
		return err
	}
	return w(err)
}

func WithStack2[T any](t T, err error) (T, error) {
	if err == nil {
		return t, err
	}
	return t, WithStack(err)
}

func WithStack3[T, U any](t T, u U, err error) (T, U, error) {
	if err == nil {
		return t, u, err
	}
	return t, u, WithStack(err)
}

func Unwrap(err error) error {
	return goerr.Unwrap(err)
}

func As(err error, target any) bool {
	return goerr.As(err, target)
}

func Is(err, target error) bool {
	return goerr.Is(err, target)
}
