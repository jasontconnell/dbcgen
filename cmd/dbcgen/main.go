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
	alltables := flag.Bool("all", false, "generate all tables with default options")
	flag.Parse()

	if *table == "" && !*alltables {
		flag.PrintDefaults()
		return
	}

	if *objname == "" {
		*objname = *table
	}

	cfg := conf.LoadConfig(*configFilename)
	var list []data.Table
	if *table != "" && !*alltables {
		tbl, err := process.Read(cfg.ConnectionString, *table)
		if err != nil {
			log.Fatal(err)
		}
		list = append(list, tbl)
	} else {
		tbls, err := process.ReadAll(cfg.ConnectionString, cfg.GenerateOptions.Ignore)
		if err != nil {
			log.Fatal(err)
		}
		list = append(list, tbls...)
	}

	renames := make(map[string]string)
	for _, r := range cfg.GenerateOptions.Renames {
		renames[r.From] = r.To
	}

	for _, tbl := range list {
		cname := *objname
		if *alltables {
			cname = tbl.Name
		}
		code, err := process.Generate(tbl, cname, cfg.GenerateOptions.NameStyle, cfg.GenerateOptions.AltNameStyle, cfg.TemplateLocation, getTypeMap(cfg.GenerateOptions.Types), cfg.GenerateOptions.IgnoreFields, renames)
		if err != nil {
			log.Fatal(err)
		}

		fn := fmt.Sprintf(cfg.OutputLocation, cname)
		err = process.Write(fn, code)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("finished.", time.Since(start))
}

func getTypeMap(types []conf.Type) data.TypeMap {
	m := make(data.TypeMap)
	for _, t := range types {

		key := fmt.Sprintf("%s", t.DbType)
		if t.CharLen > 0 {
			key += fmt.Sprintf("_%d", t.CharLen)
		}
		m[key] = data.Type{DbType: t.DbType, CharLen: t.CharLen, CodeType: t.CodeType}
	}
	return m
}
