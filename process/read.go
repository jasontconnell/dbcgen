package process

import (
	"database/sql"
	"fmt"

	"github.com/jasontconnell/dbcgen/data"
	"github.com/jasontconnell/sqlhelp"

	_ "github.com/microsoft/go-mssqldb"
)

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
	query := fmt.Sprintf(`select COLUMN_NAME, ORDINAL_POSITION, IS_NULLABLE, DATA_TYPE from INFORMATION_SCHEMA.COLUMNS where TABLE_NAME = '%s' order by ORDINAL_POSITION`, tableName)
	rows, err := sqlhelp.GetResultSet(conn, query)
	if err != nil {
		return nil, fmt.Errorf("error reading columns with query '%s'. %w", query, err)
	}

	cols := []data.Column{}
	for _, row := range rows {
		col := data.Column{
			Name:     row["COLUMN_NAME"].(string),
			Type:     row["DATA_TYPE"].(string),
			Pos:      row["ORDINAL_POSITION"].(int64),
			Nullable: row["IS_NULLABLE"].(string) != "NO",
		}
		cols = append(cols, col)
	}
	return cols, nil
}
