package template

import (
	"bytes"
	"text/template"

	"github.com/MineTakaki/go-utils/errors"
)

type (
	FuncMap = template.FuncMap
)

// FuncMapを指定してテンプレート処理を行います
func WithFuncs(src string, i interface{}, fmap FuncMap) (string, error) {
	fmaptmp := Join(CompareOperators(nil))
	for k, v := range fmap {
		fmaptmp[k] = v
	}

	tpl, err := errors.WithStack2(template.New("tmp").Funcs(fmaptmp).Parse(src))
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err := errors.WithStack(tpl.Execute(buf, i)); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// Simple シンプルなテンプレート処理を行います
func Simple(src string, i interface{}) (string, error) {
	return WithFuncs(src, i, nil)
}
