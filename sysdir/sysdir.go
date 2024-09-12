//go:build !windows

package sysdir

import (
	"os"

	"github.com/MineTakaki/go-utils/errors"
)

// Documents 既定のドキュメントフォルダーを取得します
func Documents() (string, error) {
	home, err := errors.WithStack2(os.UserHomeDir())
	if err != nil {
		return "", err
	}
	return home + string(os.PathSeparator), err
}
