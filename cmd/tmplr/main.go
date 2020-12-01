package main

import (
	"flag"
	"log"
	"strings"

	"github.com/hnakamur/tmplr"
)

type stringFlags []string

func (ss *stringFlags) String() string {
	if ss == nil {
		return ""
	}

	var b strings.Builder
	for i, s := range *ss {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(s)
	}
	return b.String()
}

func (ss *stringFlags) Set(value string) error {
	*ss = append(*ss, value)
	return nil
}

func main() {
	varFilename := flag.String("var", "var.yml", "variable YAML file")
	destFilename := flag.String("dest", "", "destination file (\"\" means stdout)")
	tmplFilename := flag.String("tmpl", "", "template filename to execute")

	var yamlRefDirs stringFlags
	flag.Var(&yamlRefDirs, "ref", "variable YAML reference directory (can specify multiple times)")

	refRecursive := flag.Bool("ref-recursive", true, "searches yaml file recursively in -ref directories")
	flag.Parse()

	cfg := &tmplr.Config{
		DestFilename:     *destFilename,
		TemplateFilename: *tmplFilename,
		VarFilename:      *varFilename,
		YAMLRefDirs:      []string(yamlRefDirs),
		YAMLRefRecursive: *refRecursive,
	}
	err := tmplr.Run(cfg)
	if err != nil {
		log.Fatal(err)
	}
}
