package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/jasontconnell/dbcgen/conf"
	"github.com/jasontconnell/dbcgen/data"
	"github.com/jasontconnell/dbcgen/process"
)

func main() {
	start := time.Now()
	configFilename := flag.String("c", "config.json", "config filename")
	table := flag.String("t", "", "table name")
	objname := flag.String("o", "", "object name")
	flag.Parse()

	if *table == "" {
		flag.PrintDefaults()
		return
	}

	if *objname == "" {
		*objname = *table
	}

	cfg := conf.LoadConfig(*configFilename)
	tbl, err := process.Read(cfg.ConnectionString, *table)
	if err != nil {
		log.Fatal(err)
	}

	renames := make(map[string]string)
	for _, r := range cfg.GenerateOptions.Renames {
		renames[r.From] = r.To
	}

	code, err := process.Generate(tbl, *objname, cfg.GenerateOptions.NameStyle, cfg.GenerateOptions.AltNameStyle, cfg.TemplateLocation, getTypeMap(cfg.GenerateOptions.Types), cfg.GenerateOptions.IgnoreFields, renames)
	if err != nil {
		log.Fatal(err)
	}

	fn := fmt.Sprintf(cfg.OutputLocation, *objname)
	err = process.Write(fn, code)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("finished.", time.Since(start))
}

func getTypeMap(types []conf.Type) data.TypeMap {
	m := make(data.TypeMap)
	for _, t := range types {
		m[t.DbType] = t.CodeType
	}
	return m
}
