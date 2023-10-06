package data

type TypeMap map[string]string

type Object struct {
	ObjectName string
	Properties []Property
}

type Property struct {
	CodeType string
	Name     string
	AltName  string
}
