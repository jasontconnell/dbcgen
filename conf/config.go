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
	IgnoreFields   []string `json:"ignoreFields"`
	Renames        []Rename `json:"renames"`
}

type Type struct {
	DbType   string `json:"dbType"`
	CharLen  int64  `json:"charLen"`
	CodeType string `json:"codeType"`
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
