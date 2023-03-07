package template

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"
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

	tpl, err := template.New("tmp").Funcs(fmaptmp).Parse(src)
	if err != nil {
		return "", errors.WithStack(err)
	}

	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, i); err != nil {
		return "", errors.WithStack(err)
	}

	return buf.String(), nil
}

// Simple シンプルなテンプレート処理を行います
func Simple(src string, i interface{}) (string, error) {
	return WithFuncs(src, i, nil)
}
