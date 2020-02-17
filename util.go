package i18n

import (
	"regexp"
	"strings"
)

var localeMatcher = regexp.MustCompile(`[a-z]{2}([_-])[A-Z]{2}|[a-z]{2}\.`)

// guessLocaleFromFilename tries to guess the locale from the given string. A supported filename
// looks like strings-de-DE.xml
func guessLocaleFromFilename(str string) string {
	tmp := localeMatcher.FindString(str)
	if strings.HasSuffix(tmp, ".") {
		return tmp[:len(tmp)-1]
	}
	return tmp
}
