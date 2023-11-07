package data

type Type struct {
	DbType   string
	CharLen  int64
	CodeType string
}

type TypeMap map[string]Type

type Object struct {
	ObjectName      string
	ObjectCleanName string
	Properties      []Property
}

type Property struct {
	CodeType      string
	ColumnName    string
	ColumnLength  int64
	DbTypeDef     string
	CleanName     string
	AltName       string
	Key           bool
	AutoIncrement bool
}
