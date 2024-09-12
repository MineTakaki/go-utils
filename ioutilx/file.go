package ioutilx

import (
	"os"

	"github.com/MineTakaki/go-utils/errors"
)

// Open os.Openのラッパー関数
func Open(name string) (*os.File, error) {
	return errors.WithStack2(os.Open(name))
}

// OpenFn
func OpenFn(name string, fn func(f *os.File) error) error {
	f, err := errors.WithStack2(os.Open(name))
	if err != nil {
		return err
	}
	if fn != nil {
		if err := fn(f); err != nil {
			f.Close()
			return err
		}
	}
	return errors.WithStack(f.Close())
}

// Create os.Createのラッパー関数
func Create(name string) (*os.File, error) {
	return errors.WithStack2(os.Create(name))
}

// CreateFn
func CreateFn(name string, fn func(f *os.File) error) error {
	f, err := errors.WithStack2(os.Create(name))
	if err != nil {
		return err
	}
	if fn != nil {
		if err := fn(f); err != nil {
			f.Close()
			return err
		}
	}
	return errors.WithStack(f.Close())
}

// CreateTemp os.CreateTempのラッパー関数
func CreateTemp(dir, pattern string) (*os.File, error) {
	return errors.WithStack2(os.CreateTemp(dir, pattern))
}

// CreateTempFn
func CreateTempFn(dir, pattern string, fn func(f *os.File) error) error {
	f, err := errors.WithStack2(os.CreateTemp(dir, pattern))
	if err != nil {
		return err
	}
	if fn != nil {
		if err := fn(f); err != nil {
			f.Close()
			return err
		}
	}
	return errors.WithStack(f.Close())
}
