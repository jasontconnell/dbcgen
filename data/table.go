package data

type Table struct {
	Name    string
	Columns []Column
}

type Column struct {
	Name     string
	Type     string
	Pos      int64
	Nullable bool
}
