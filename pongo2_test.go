package tmplr

import (
	"testing"

	"github.com/flosch/pongo2/v6"
)

func TestPongo2AutoEscape(t *testing.T) {
	pongo2.SetAutoescape(false)
	tpl, err := pongo2.FromString("{{ arg }}")
	if err != nil {
		t.Fatal(err)
	}
	got, err := tpl.Execute(pongo2.Context{"arg": "'--extra-vars=\"var1=value\"'"})
	if err != nil {
		t.Fatal(err)
	}
	want := "'--extra-vars=\"var1=value\"'"
	if got != want {
		t.Errorf("result unmatch, got=%s, want=%s", got, want)
	}
}
