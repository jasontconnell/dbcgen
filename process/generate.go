package process

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"text/template"

	"github.com/jasontconnell/dbcgen/data"
)

func Generate(tbl data.Table, objname, nameStyle, altNameStyle, tmplLoc string, typeMap data.TypeMap, ignoreFields []string, renames map[string]string) ([]byte, error) {
	obj := getObject(tbl, objname, nameStyle, altNameStyle, typeMap, ignoreFields, renames)
	tmp, err := template.New("").ParseFiles(tmplLoc)
	if err != nil {
		return nil, fmt.Errorf("can't parse template. %w", err)
	}

	buffer := new(bytes.Buffer)
	_, templateName := filepath.Split(tmplLoc)
	err = tmp.ExecuteTemplate(buffer, templateName, obj)

	if err != nil {
		return nil, fmt.Errorf("error processing template. %w", err)
	}

	return buffer.Bytes(), nil
}

func getObject(tbl data.Table, name, nameStyle, altNameStyle string, typeMap data.TypeMap, ignoreFields []string, renames map[string]string) data.Object {
	nameFunc := getCleanNameFunc(nameStyle)
	altNameFunc := getCleanNameFunc(altNameStyle)
	im := make(map[string]bool)
	for _, s := range ignoreFields {
		im[s] = true
	}

	obj := data.Object{ObjectName: name}
	for _, c := range tbl.Columns {
		if _, ok := im[c.Name]; ok {
			continue
		}
		ct, ok := typeMap[c.Type]
		if !ok {
			log.Println("no corresponding type for", c.Type)
			continue
		}

		lkey := fmt.Sprintf("%s_%d", c.Type, c.CharLen)
		special, ok := typeMap[lkey]
		if ok {
			ct = special
		}

		colname := c.Name
		name := c.Name
		if n, ok := renames[c.Name]; ok {
			name = n
		}

		p := data.Property{CodeType: ct.CodeType, CleanName: nameFunc(name), ColumnName: colname, ColumnLength: c.CharLen, AltName: altNameFunc(name), Key: c.PrimaryKey}
		obj.Properties = append(obj.Properties, p)
	}
	return obj
}
