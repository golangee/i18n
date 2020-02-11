package i18n

// nolint: goimports // the linter is broken
import (
	"fmt"
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
)

// a value is a contract which is implemented by each kind of message value, like simple, array or plural.
type value interface {
	TextArray() ([]string, error)

	Text(args ...interface{}) (string, error)

	QuantityText(tag language.Tag, quantity int, args ...interface{}) (string, error)
}

// A pluralValue is value with CLDR plural rules, see also
// https://unicode.org/reports/tr35/tr35-numbers.html#Language_Plural_Rules
type pluralValue struct {
	ID    string
	Zero  string
	One   string
	Two   string
	Few   string
	Many  string
	Other string
}

// StringArray returns the Other value within a one element array
func (p pluralValue) TextArray() ([]string, error) {
	return []string{p.Other}, nil
}

// String returns other
func (p pluralValue) Text(args ...interface{}) (string, error) {
	return fmt.Sprintf(p.Other, args...), nil
}

// QuantityString returns the grammatical plural for the internal Plural implementation
func (p pluralValue) QuantityText(tag language.Tag, quantity int, args ...interface{}) (string, error) {
	// MatchPlural is over engineered: f and t are not used anyway and according to its test, v and w must be kept 0
	// to just get the plural for a natural number
	form := plural.Cardinal.MatchPlural(tag, quantity, 0, 0, 0, 0)
	switch form {
	case plural.Zero:
		return p.fallback(p.Zero, args...)
	case plural.One:
		return p.fallback(p.One, args...)
	case plural.Two:
		return p.fallback(p.Two, args...)
	case plural.Few:
		return p.fallback(p.Few, args...)
	case plural.Many:
		return p.fallback(p.Many, args...)
	case plural.Other:
		fallthrough
	default:
		return fmt.Sprintf(p.Other, args...), nil
	}
}

// fallback uses Other if text is empty
func (p pluralValue) fallback(text string, args ...interface{}) (string, error) {
	if len(text) == 0 {
		return fmt.Sprintf(p.Other, args...), nil
	}

	return fmt.Sprintf(text, args...), nil
}

// A simpleValue just holds a text
type simpleValue struct {
	ID     string
	String string
}

// StringArray returns the text in a single element array
func (s simpleValue) TextArray() ([]string, error) {
	return []string{s.String}, nil
}

// String interpolates and returns the text
func (s simpleValue) Text(args ...interface{}) (string, error) {
	return fmt.Sprintf(s.String, args...), nil
}

// QuantityString is equivalent to String
func (s simpleValue) QuantityText(tag language.Tag, quantity int, args ...interface{}) (string, error) {
	return s.Text(args...)
}

// An arrayValue holds an ordered bunch of strings
type arrayValue struct {
	ID      string
	Strings []string
}

// TextArray returns a defensive copy
func (a arrayValue) TextArray() ([]string, error) {
	tmp := make([]string, len(a.Strings))
	copy(tmp, a.Strings)

	return tmp, nil
}

// Text returns the first array element or the empty string
func (a arrayValue) Text(args ...interface{}) (string, error) {
	if len(a.Strings) > 0 {
		return fmt.Sprintf(a.Strings[0], args...), nil
	}

	return "", nil
}

// QuantityText just returns text
func (a arrayValue) QuantityText(tag language.Tag, quantity int, args ...interface{}) (string, error) {
	return a.Text(args...)
}
