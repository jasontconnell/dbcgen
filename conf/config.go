package conf

import "github.com/jasontconnell/conf"

type Config struct {
	ConnectionString string          `json:"connectionString"`
	TemplateLocation string          `json:"templateLocation"`
	OutputLocation   string          `json:"outputLocation"`
	GenerateOptions  GenerateOptions `json:"generate"`
}

type GenerateOptions struct {
	Types          []Type   `json:"types"`
	NameStyle      string   `json:"nameStyle"`
	AltNameStyle   string   `json:"altNameStyle"`
	NullableFormat string   `json:"nullableFormat"`
	Ignore         []string `json:"ignore"`
	IgnoreFields   []string `json:"ignoreFields"`
	Renames        []Rename `json:"renames"`
}

type Type struct {
	DbType           string `json:"dbType"`
	CodeType         string `json:"codeType"`
	NullableCodeType string `json:"nullableCodeType"`
	StringType       string `json:"stringType"`
}

type Rename struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func LoadConfig(filename string) Config {
	var cfg Config
	conf.LoadConfig(filename, &cfg)
	return cfg
}
