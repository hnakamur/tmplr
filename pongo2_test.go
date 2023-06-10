package tmplr

import (
	"os"
	"path/filepath"
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

func TestErrNotDiscardedForFileOutput(t *testing.T) {
	dir := t.TempDir()
	destFilename := filepath.Join(dir, "vanish.md")
	tmplFilename := filepath.Join(dir, "vanish.j2")
	varFilename := filepath.Join(dir, "vars.yaml")

	const varsStr = `a: b`
	const tmplStr = `{% macro show_zone(zone) %}
	zone is {{ zone }}
	{% endmacro %}

	{{ show_zone("osk", "foo") }}
`

	if err := os.WriteFile(tmplFilename, []byte(tmplStr), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(varFilename, []byte(varsStr), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := &Config{
		DestFilename:     destFilename,
		TemplateFilename: tmplFilename,
		VarFilename:      varFilename,
	}
	if err := Run(cfg); err == nil {
		t.Errorf("should got an error, but got no error")
	}
}
