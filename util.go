// Copyright 2020 Torben Schinke
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
