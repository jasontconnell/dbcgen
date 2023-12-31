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

	obj := data.Object{ObjectName: name, ObjectCleanName: nameFunc(name)}
	for _, c := range tbl.Columns {
		if _, ok := im[c.Name]; ok {
			continue
		}
		ct, ok := typeMap[c.Type]
		if !ok {
			log.Println("no corresponding type for", c.Type)
			continue
		}

		colname := c.Name
		name := c.Name
		if n, ok := renames[c.Name]; ok {
			name = n
		}

		typedef := c.Type
		if c.CharLen > 0 {
			typedef += fmt.Sprintf("(%d)", c.CharLen)
		} else if c.NumPrecision > 0 && c.Type != "int" {
			typedef += fmt.Sprintf("(%d, %d)", c.NumPrecision, c.NumScale)
		}

		codeType := ct.CodeType
		if c.Nullable && ct.NullableCodeType != "" {
			codeType = ct.NullableCodeType
		}

		if c.CharLen > 1 && ct.StringType != "" {
			codeType = ct.StringType
		}

		p := data.Property{
			CodeType:      codeType,
			CleanName:     nameFunc(name),
			ColumnName:    colname,
			DbTypeDef:     typedef,
			ColumnLength:  c.CharLen,
			AltName:       altNameFunc(name),
			Key:           c.PrimaryKey,
			ForeignKey:    c.ForeignKey,
			AutoIncrement: c.AutoIncrement,
			Nullable:      c.Nullable,
		}
		obj.Properties = append(obj.Properties, p)
	}
	return obj
}
