package tmplr

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"

	"github.com/flosch/pongo2/v4"
	"github.com/goccy/go-yaml"
)

type Config struct {
	DestFilename     string
	TemplateFilename string

	VarFilename      string
	YAMLRefDirs      []string
	YAMLRefRecursive bool
}

func Run(cfg *Config) (err error) {
	var data map[string]interface{}
	data, err = readYAMLFile(cfg.VarFilename, cfg.YAMLRefDirs, cfg.YAMLRefRecursive)
	if err != nil {
		return err
	}

	var tmpl *pongo2.Template
	if len(cfg.TemplateFilename) == 0 {
		var data []byte
		data, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		tmpl, err = pongo2.FromString(string(data))
		if err != nil {
			return err
		}
	} else {
		tmpl, err = pongo2.FromFile(cfg.TemplateFilename)
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

	err = tmpl.ExecuteWriter(data, w)
	if err != nil {
		return err
	}

	return nil
}

func readYAMLFile(filename string, refDirs []string, refRecursive bool) (map[string]interface{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(bufio.NewReader(file),
		yaml.ReferenceDirs(refDirs...),
		yaml.RecursiveDir(refRecursive))
	var v map[string]interface{}
	if err := d.Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}
