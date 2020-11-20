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
	"fmt"
	"github.com/golangee/i18n/internal"
	"github.com/golangee/log"
	"github.com/golangee/log/ecs"
	"io"
	"os"
	"sort"
)

// ErrTextNotFound is the sentinel error for a named string which is not available
var ErrTextNotFound = fmt.Errorf("string not found")

var allResources = newLocalizations() //nolint: gochecknoglobals

var logger = log.NewLogger(ecs.Log("i18n"))

// Import takes the importer and locale and updates the according internal localization resources.
// The order of import is relevant, because it determines the fallback matching logic. Import your default fallback
// language first.
func Import(importer Importer, locale string, src io.Reader) error {
	res := allResources.Configure(locale)
	err := importer.Import(res, src)

	if err != nil {
		return fmt.Errorf("failed to parse: %w", err)
	}

	return nil
}

// ImportFile is a convenience method for Import. It detects the locale from the file name
func ImportFile(importer Importer, fname string) error {
	file, err := os.Open(fname)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}

	defer func() {
		_ = file.Close()
	}()

	return Import(importer, guessLocaleFromFilename(fname), file)
}

// ImportValue adds or replaces any existing value
func ImportValue(value Value) {
	res := allResources.Configure(value.Locale())
	value = value.updateTag(res.tag)
	res.mutex.Lock()
	defer res.mutex.Unlock()

	if _, has := res.values[value.ID()]; has {
		logger.Print(ecs.Warn(), ecs.Msg("replacing already translated value"), log.V("key", value.ID()))
	}
	res.values[value.ID()] = value
}

// From returns the best matching text Resources to the given set of matching locales
func From(locales ...string) *Resources {
	return allResources.Match(locales...)
}

// Validates checks the current state of the global localizations to see if everything is fine. If no error is returned,
// you can be sure that at least every key is translated in every language and the printf directives are consistent
// with each other.
func Validate() error {
	allResources.translationsMutex.RLock()
	defer allResources.translationsMutex.RUnlock()

	tmp := make([]*Resources, 0, len(allResources.translations))
	for _, res := range allResources.translations {
		tmp = append(tmp, res)
	}
	return validate(tmp)
}

// TranslationPriority updates the resolution order and removes unwanted translations. "und" is the undefined default
// locale.
func TranslationPriority(locales ...string) {
	allResources.SetTranslationPriority(locales)
}

// Locales returns all translated locales
func Locales() []string {
	var res []string
	for k, _ := range allResources.translations {
		res = append(res, k.String())
	}
	sort.Strings(res)
	return res
}

// Bundle (re)generates all localizations in the current working directory.
func Bundle() error {
	dir, err := internal.ModRootDir()
	if err != nil {
		return fmt.Errorf("unable to get current working directory: %w", err)
	}
	fmt.Println(dir)
	gen := newGoGenerator(dir)
	err = gen.Scan()
	if err != nil {
		return fmt.Errorf("cannot scan module: %w", err)
	}

	err = gen.Emit()
	if err != nil {
		return fmt.Errorf("unable to generate source code: %w", err)
	}
	return nil
}
