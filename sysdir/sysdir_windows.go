// +build windows

package sysdir

import (
	"os"

	"github.com/pkg/errors"
	"golang.org/x/sys/windows/registry"
)

//Documents 既定のドキュメントフォルダーを取得します
func Documents() (string, error) {
	regkey, err := registry.OpenKey(
		registry.CURRENT_USER,
		`Software\\Microsoft\\Windows\\CurrentVersion\\Explorer\\User Shell Folders`,
		registry.READ,
	)
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer regkey.Close()

	var v string
	v, _, err = regkey.GetStringValue("Personal")
	if err != nil {
		return "", errors.WithStack(err)
	}
	return v + string(os.PathSeparator), nil
}
