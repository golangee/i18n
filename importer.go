package i18n

// nolint: goimports // the linter is broken
import (
	"fmt"
	"github.com/worldiety/i18n/android"
	"io"
	"strings"
)

const (
	zero  = "zero"
	one   = "one"
	two   = "two"
	few   = "few"
	many  = "many"
	other = "other"
)

// An Importer parses data formats and imports that into resources.
type Importer interface {
	// Import tries to parse the src bytes and imports that into the given resources.
	Import(dst *Resources, src io.Reader) error
}

// An AndroidImporter supports the android strings xml format with simple strings, interpolation, indices, plurals
// and arrays. However, references and indices with more than 9 not.
type AndroidImporter struct {
}

// Import tries to parse the src bytes and imports that into the given resources.
func (a AndroidImporter) Import(dst *Resources, src io.Reader) error {
	aRes, err := android.Read(src)
	if err != nil {
		return fmt.Errorf("failed to import android resources: %w", err)
	}

	importAndroid(dst, aRes)

	return nil
}

// importAndroid copies and converts the given android resources into our i18n resources.
func importAndroid(dst *Resources, src android.Resources) {
	dst.mutex.Lock()
	defer dst.mutex.Unlock()

	for _, str := range src.Strings {
		dst.values[str.Name] = simpleValue{
			ID:     str.Name,
			String: android.Decode(str.Text),
		}
	}

	for _, pl := range src.Plurals {
		val := pluralValue{
			ID: pl.Name,
		}

		for _, item := range pl.Items {
			switch strings.ToLower(item.Quantity) {
			case zero:
				val.Zero = android.Decode(item.Text)
			case one:
				val.One = android.Decode(item.Text)
			case two:
				val.Two = android.Decode(item.Text)
			case few:
				val.Few = android.Decode(item.Text)
			case many:
				val.Many = android.Decode(item.Text)
			case other:
				fallthrough
			default:
				val.Other = android.Decode(item.Text)
			}
		}

		dst.values[pl.Name] = val
	}

	for _, arr := range src.StringArrays {
		tmp := make([]string, 0, len(arr.Items))
		for _, s := range arr.Items {
			tmp = append(tmp, android.Decode(s))
		}

		dst.values[arr.Name] = arrayValue{
			ID:      arr.Name,
			Strings: tmp,
		}
	}
}
