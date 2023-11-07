package process

import (
	"regexp"
	"strings"
	"unicode"
)

var cleanreg *regexp.Regexp = regexp.MustCompile("[^a-zA-Z0-9 \\-_]+")

func clean(name string) string {
	return cleanreg.ReplaceAllString(name, "")
}

func getCleanName(name string) string {
	prefix := ""
	if strings.IndexAny(name, "0123456789") == 0 {
		prefix = "_"
	}
	return prefix + strings.Replace(strings.Replace(strings.Title(clean(name)), "-", "", -1), " ", "", -1)
}

func getExactName(name string) string {
	prefix := ""
	if strings.IndexAny(name, "0123456789") == 0 {
		prefix = "_"
	}
	return prefix + strings.Replace(strings.Replace(clean(name), "-", "", -1), " ", "", -1)
}

func getUnderscoreUppercaseName(name string) string {
	name = strings.Replace(strings.Title(name), " ", "_", -1)
	return getCleanName(name)
}

func getUnderscoreLowercaseName(name string) string {
	return strings.ToLower(getUnderscoreUppercaseName(name))
}

func getCamelCaseName(name string) string {
	asTitle := getCleanName(name)
	start := 0
	for _, x := range asTitle {
		if unicode.IsUpper(x) {
			start++
		} else {
			break
		}
	}
	if start > 1 && start != len(asTitle) {
		start--
	}
	return strings.ToLower(string(asTitle[:start])) + asTitle[start:]
}

func getCleanNameFunc(setting string) func(string) string {
	var ret func(string) string
	switch strings.ToLower(strings.ReplaceAll(setting, "_", "")) {
	case "", "pascalcase":
		ret = getCleanName
	case "camelcase":
		ret = getCamelCaseName
	case "pascalcaseunderscore":
		ret = getUnderscoreUppercaseName
	case "lowercaseunderscore":
		ret = getUnderscoreLowercaseName
	case "exact":
		ret = getExactName
	default:
		panic("Name style not recognized: " + setting)
	}
	return ret
}
