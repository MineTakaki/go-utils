package template

import "testing"

func TestSimple(t *testing.T) {
	m := map[string]interface{}{}
	m["A"] = "A"

	txt, err := Simple(test001tmplate, m)
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	t.Logf("%s", txt)
}

var test001tmplate = `
string eq:{{if eq .A "A" }}OK{{else}}NG{{end}}
string ne:{{if ne .A "A" }}NG{{else}}OK{{end}}
string lt:{{if lt .A "B" }}OK{{else}}NG{{end}}
string lt:{{if lt .A "A" }}NG{{else}}OK{{end}}
string le:{{if le .A "B" }}OK{{else}}NG{{end}}
string le:{{if le .A "A" }}OK{{else}}NG{{end}}
string gt:{{if gt .A "!" }}OK{{else}}NG{{end}}
string gt:{{if gt .A "A" }}NG{{else}}OK{{end}}
string ge:{{if ge .A "!" }}OK{{else}}NG{{end}}
string gt:{{if ge .A "A" }}OK{{else}}NG{{end}}
`
