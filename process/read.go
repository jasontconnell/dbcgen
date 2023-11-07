package process

import (
	"database/sql"
	"fmt"

	"github.com/jasontconnell/dbcgen/data"
	"github.com/jasontconnell/sqlhelp"

	_ "github.com/microsoft/go-mssqldb"
)

func ReadAll(connstr string) ([]data.Table, error) {
	conn, err := sql.Open("mssql", connstr)
	if err != nil {
		return nil, fmt.Errorf("couldn't connect. %w", err)
	}
	defer conn.Close()

	query := "select TABLE_NAME from INFORMATION_SCHEMA.TABLES"
	rows, err := sqlhelp.GetResultSet(conn, query)
	if err != nil {
		return nil, fmt.Errorf("couldn't execute query %s. %w", query, err)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("no tables found")
	}

	tables := []data.Table{}
	for _, r := range rows {
		name := r["TABLE_NAME"].(string)
		tbl := data.Table{Name: name}
		cols, err := readCols(conn, name)
		if err != nil {
			return nil, fmt.Errorf("error reading columns for %s. %w", name, err)
		}
		tbl.Columns = cols
		tables = append(tables, tbl)
	}
	return tables, nil
}

func Read(connstr, table string) (data.Table, error) {
	tbl := data.Table{}
	conn, err := sql.Open("mssql", connstr)
	if err != nil {
		return tbl, fmt.Errorf("couldn't connect. %w", err)
	}
	defer conn.Close()

	query := fmt.Sprintf(`select TABLE_NAME from INFORMATION_SCHEMA.TABLES where TABLE_NAME = '%s'`, table)

	rows, err := sqlhelp.GetResultSet(conn, query)
	if err != nil {
		return tbl, fmt.Errorf("couldn't execute query %s. %w", query, err)
	}

	if len(rows) == 0 {
		return tbl, fmt.Errorf("table not found: %s", table)
	}

	tbl.Name = table
	cols, err := readCols(conn, tbl.Name)
	if err != nil {
		return tbl, fmt.Errorf("couldn't read columns. %w", err)
	}
	tbl.Columns = cols
	return tbl, nil
}

func readCols(conn *sql.DB, tableName string) ([]data.Column, error) {
	query := fmt.Sprintf(`
	select c.COLUMN_NAME, c.ORDINAL_POSITION, c.IS_NULLABLE, c.DATA_TYPE, c.CHARACTER_MAXIMUM_LENGTH, case when tc.CONSTRAINT_TYPE = 'PRIMARY KEY' then 1 else 0 end as IS_PRIMARY_KEY
	from INFORMATION_SCHEMA.COLUMNS c 
		left join INFORMATION_SCHEMA.CONSTRAINT_COLUMN_USAGE cons
			join INFORMATION_SCHEMA.TABLE_CONSTRAINTS tc
				on tc.CONSTRAINT_NAME = cons.CONSTRAINT_NAME
					and tc.TABLE_NAME = cons.TABLE_NAME
					and tc.CONSTRAINT_TYPE = 'PRIMARY KEY'
			on c.TABLE_NAME = cons.TABLE_NAME
				and c.COLUMN_NAME = cons.COLUMN_NAME
	where
		c.TABLE_NAME = '%s' 
	order by c.ORDINAL_POSITION
		`, tableName)
	rows, err := sqlhelp.GetResultSet(conn, query)
	if err != nil {
		return nil, fmt.Errorf("error reading columns with query '%s'. %w", query, err)
	}

	cols := []data.Column{}
	for _, row := range rows {
		var clen int64
		if row["CHARACTER_MAXIMUM_LENGTH"] != nil {
			clen = row["CHARACTER_MAXIMUM_LENGTH"].(int64)
		}
		col := data.Column{
			Name:       row["COLUMN_NAME"].(string),
			Type:       row["DATA_TYPE"].(string),
			CharLen:    clen,
			Pos:        row["ORDINAL_POSITION"].(int64),
			Nullable:   row["IS_NULLABLE"].(string) != "NO",
			PrimaryKey: row["IS_PRIMARY_KEY"].(int64) == 1,
		}
		cols = append(cols, col)
	}
	return cols, nil
}
