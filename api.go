package i18n

import (
	"fmt"
	"io"
	"os"
)

// ErrTextNotFound is the sentinel error for a named string which is not available
var ErrTextNotFound = fmt.Errorf("string not found")

var allResources = newLocalizations() //nolint: gochecknoglobals

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

// ImportFile is a convenience method for Import.
func ImportFile(importer Importer, locale string, fname string) error {
	file, err := os.Open(fname)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}

	defer func() {
		_ = file.Close()
	}()

	return Import(importer, locale, file)
}

// ImportValues adds or replaces any existing value
func ImportValue(value Value) {
	res := allResources.Configure(value.Locale())
	value.updateTag(res.tag)
	res.mutex.Lock()
	defer res.mutex.Unlock()

	res.values[value.ID()] = value
}

// From returns the best matching text Resources to the given set of matching locales
func From(locales ...string) *Resources {
	return allResources.Match(locales...)
}

// Validates checks the current state of the global localizations to see if everything is fine. If no error is returned,
// you can be sure that at least every key is translated in every language and the printf directives are consistent
// with each other.
func Validate() []error {
	allResources.translationsMutex.RLock()
	defer allResources.translationsMutex.RUnlock()

	tmp := make([]*Resources, 0, len(allResources.translations))
	for _, res := range allResources.translations {
		tmp = append(tmp, res)
	}
	return validate(tmp)
}
