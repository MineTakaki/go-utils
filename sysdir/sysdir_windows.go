//go:build windows

package sysdir

import (
	"os"

	"github.com/MineTakaki/go-utils/errors"
	"golang.org/x/sys/windows/registry"
)

// Documents 既定のドキュメントフォルダーを取得します
func Documents() (string, error) {
	regkey, err := errors.WithStack2(registry.OpenKey(
		registry.CURRENT_USER,
		`Software\\Microsoft\\Windows\\CurrentVersion\\Explorer\\User Shell Folders`,
		registry.READ,
	))
	if err != nil {
		return "", err
	}
	defer regkey.Close()

	var v string
	v, _, err = errors.WithStack3(regkey.GetStringValue("Personal"))
	if err != nil {
		return "", err
	}
	return v + string(os.PathSeparator), nil
}
