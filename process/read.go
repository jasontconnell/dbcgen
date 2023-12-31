package process

import (
	"database/sql"
	"fmt"

	"github.com/jasontconnell/dbcgen/data"
	"github.com/jasontconnell/sqlhelp"

	_ "github.com/microsoft/go-mssqldb"
)

func ReadAll(connstr string, ignore []string) ([]data.Table, error) {
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

	mignore := make(map[string]bool)
	for _, s := range ignore {
		mignore[s] = true
	}

	tables := []data.Table{}
	for _, r := range rows {
		name := r["TABLE_NAME"].(string)
		if _, ok := mignore[name]; ok {
			continue
		}
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
	select 
		c.COLUMN_NAME, 
		c.ORDINAL_POSITION, 
		c.IS_NULLABLE, 
		c.DATA_TYPE, 
		c.CHARACTER_MAXIMUM_LENGTH, 
		c.NUMERIC_PRECISION, 
		c.NUMERIC_SCALE, 
		case when tcpk.CONSTRAINT_TYPE = 'PRIMARY KEY' then 1 else 0 end as IS_PRIMARY_KEY,
        case when tcfk.CONSTRAINT_TYPE = 'FOREIGN KEY' then 1 else 0 end as IS_FOREIGN_KEY,
		COLUMNPROPERTY(object_id(c.TABLE_SCHEMA+'.'+c.TABLE_NAME), c.COLUMN_NAME, 'IsIdentity') as IS_IDENTITY
	from INFORMATION_SCHEMA.COLUMNS c 
		left join INFORMATION_SCHEMA.CONSTRAINT_COLUMN_USAGE cons
			join INFORMATION_SCHEMA.TABLE_CONSTRAINTS tcpk
				on tcpk.CONSTRAINT_NAME = cons.CONSTRAINT_NAME
					and tcpk.TABLE_NAME = cons.TABLE_NAME
					and tcpk.CONSTRAINT_TYPE = 'PRIMARY KEY'
			on c.TABLE_NAME = cons.TABLE_NAME
				and c.COLUMN_NAME = cons.COLUMN_NAME
		left join INFORMATION_SCHEMA.CONSTRAINT_COLUMN_USAGE consfk
            join INFORMATION_SCHEMA.TABLE_CONSTRAINTS tcfk
				on tcfk.CONSTRAINT_NAME = consfk.CONSTRAINT_NAME
					and tcfk.TABLE_NAME = consfk.TABLE_NAME
					and tcfk.CONSTRAINT_TYPE = 'FOREIGN KEY'
			on c.TABLE_NAME = consfk.TABLE_NAME
				and c.COLUMN_NAME = consfk.COLUMN_NAME
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
		col := data.Column{
			Name:          row["COLUMN_NAME"].(string),
			Type:          row["DATA_TYPE"].(string),
			CharLen:       readNullable[int64](row["CHARACTER_MAXIMUM_LENGTH"], 0),
			NumPrecision:  readNullable[int64](row["NUMERIC_PRECISION"], 0),
			NumScale:      readNullable[int64](row["NUMERIC_SCALE"], 0),
			Pos:           row["ORDINAL_POSITION"].(int64),
			Nullable:      row["IS_NULLABLE"].(string) != "NO",
			PrimaryKey:    row["IS_PRIMARY_KEY"].(int64) == 1,
			ForeignKey:    row["IS_FOREIGN_KEY"].(int64) == 1,
			AutoIncrement: row["IS_IDENTITY"].(int64) == 1,
		}
		cols = append(cols, col)
	}
	return cols, nil
}

func readNullable[T any](val interface{}, def T) T {
	if val == nil {
		return def
	}

	return val.(T)
}
