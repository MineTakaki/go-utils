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
	fn := fnNew
	if fn == nil {
		return goerr.New(text)
	}
	return fn(text)
}

func Errorf(format string, args ...interface{}) error {
	text := fmt.Sprintf(format, args...)
	fn := fnNew
	if fn == nil {
		return goerr.New(text)
	}
	return fn(text)
}

func Wrap(err error, text string) error {
	w := fnWrap
	if w == nil {
		return errors.Wrap(err, text)
	}
	return w(err, text)
}

func Wrapf(err error, format string, args ...interface{}) error {
	text := fmt.Sprintf(format, args...)
	w := fnWrap
	if w == nil {
		return errors.Wrap(err, text)
	}
	return w(err, text)
}

func WithStack(err error) error {
	w := fnWithStack
	if w == nil {
		return err
	}
	return w(err)
}

func WithStack2[T any](t T, err error) (T, error) {
	return t, WithStack(err)
}

func WithStack3[T, U any](t T, u U, err error) (T, U, error) {
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
