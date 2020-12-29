// +build !windows

package sysdir

import (
	"os"

	"github.com/pkg/errors"
)

//Documents 既定のドキュメントフォルダーを取得します
func Documents() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.WithStack(err)
	}
	return home + string(os.PathSeparator), err
}
