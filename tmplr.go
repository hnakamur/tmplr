package tmplr

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/goccy/go-yaml"
)

func Run(destFilename, varFilename, tmplName, tmplPattern string) (err error) {
	var data interface{}
	data, err = ReadYAMLFile(varFilename)
	if err != nil {
		return err
	}

	var tmpl *template.Template
	if tmplPattern == "" || tmplPattern == "-" {
		var data []byte
		data, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		tmpl, err = template.New(tmplName).Parse(string(data))
		if err != nil {
			return err
		}
	} else {
		tmpl, err = template.ParseGlob(tmplPattern)
		if err != nil {
			return err
		}
	}

	var w io.Writer
	if destFilename == "" || destFilename == "-" {
		w = os.Stdout
	} else {
		var file *os.File
		file, err = os.Create(destFilename)
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

	err = tmpl.ExecuteTemplate(w, tmplName, data)
	if err != nil {
		return err
	}

	return nil
}

func ReadYAMLFile(filename string) (interface{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(bufio.NewReader(file))
	var v interface{}
	if err := d.Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}
