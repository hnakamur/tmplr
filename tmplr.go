package tmplr

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/goccy/go-yaml"
)

type Config struct {
	DestFilename string

	TemplateName    string
	TemplatePattern string

	VarFilename      string
	YAMLRefDirs      []string
	YAMLRefRecursive bool
}

func Run(cfg *Config) (err error) {
	var data interface{}
	data, err = readYAMLFile(cfg.VarFilename, cfg.YAMLRefDirs, cfg.YAMLRefRecursive)
	if err != nil {
		return err
	}

	base := template.New(cfg.TemplateName).Funcs(sprig.TxtFuncMap())
	var tmpl *template.Template
	if cfg.TemplatePattern == "" {
		var data []byte
		data, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		tmpl, err = base.Parse(string(data))
		if err != nil {
			return err
		}
	} else {
		tmpl, err = base.ParseGlob(cfg.TemplatePattern)
		if err != nil {
			return err
		}
	}

	var w io.Writer
	if cfg.DestFilename == "" {
		w = os.Stdout
	} else {
		var file *os.File
		file, err = os.Create(cfg.DestFilename)
		if err != nil {
			return err
		}
		bw := bufio.NewWriter(file)
		defer func() {
			err = bw.Flush()
			err2 := file.Sync()
			if err == nil {
				err = err2
			}
			file.Close()
		}()

		w = bw
	}

	err = tmpl.ExecuteTemplate(w, cfg.TemplateName, data)
	if err != nil {
		return err
	}

	return nil
}

func readYAMLFile(filename string, refDirs []string, refRecursive bool) (interface{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(bufio.NewReader(file),
		yaml.ReferenceDirs(refDirs...),
		yaml.RecursiveDir(refRecursive))
	var v interface{}
	if err := d.Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}
