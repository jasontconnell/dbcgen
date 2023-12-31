package data

type Table struct {
	Name    string
	Columns []Column
}

type Column struct {
	CleanName     string
	AltName       string
	Name          string
	Type          string
	CharLen       int64
	NumPrecision  int64
	NumScale      int64
	Pos           int64
	Nullable      bool
	PrimaryKey    bool
	ForeignKey    bool
	AutoIncrement bool
}
