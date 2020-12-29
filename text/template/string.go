package template

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"
)

//Simple シンプルなテンプレート処理を行います
func Simple(src string, i interface{}) (txt string, err error) {
	var tpl *template.Template
	tpl, err = template.New("tmp").Funcs(Join(CompareOperators(nil))).Parse(src)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	buf := new(bytes.Buffer)

	err = tpl.Execute(buf, i)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	txt = buf.String()
	return
}
