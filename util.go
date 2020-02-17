package i18n

import (
	"path/filepath"
	"regexp"
	"strings"
)

var localeMatcher = regexp.MustCompile(`[a-z]{2}([_-])[A-Z]{2}|[a-z]{2}\.`)

// guessLocaleFromFilename tries to guess the locale from the given string. A supported filename
// looks like strings-de-DE.xml. If the file name is just strings.xml
func guessLocaleFromFilename(str string) string {
	if strings.Contains(str, string(filepath.Separator)) {
		str = filepath.Base(str)
	}
	if strings.ToLower(str) == "strings.xml" {
		return "und"
	}
	tmp := localeMatcher.FindString(str)
	if strings.HasSuffix(tmp, ".") {
		return tmp[:len(tmp)-1]
	}
	return tmp
}
