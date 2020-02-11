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

// From returns the best matching text Resources to the given set of matching locales
func From(locales ...string) *Resources {
	return allResources.Match(locales...)
}
