package main

import (
	"flag"
	"log"

	"github.com/hnakamur/tmplr"
)

func main() {
	varFilename := flag.String("var", "var.yml", "variable YAML file")
	destFilename := flag.String("dest", "", "destination file (\"\" or \"-\" means stdout)")
	tmplName := flag.String("name", "main", "template name to execute")
	tmplPattern := flag.String("template", "", "template filename glob pattern (\"\" or \"-\" means stdin)")
	flag.Parse()

	err := tmplr.Run(*destFilename, *varFilename, *tmplName, *tmplPattern)
	if err != nil {
		log.Fatal(err)
	}
}
