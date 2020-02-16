package i18n

// nolint: goimports // the linter is broken
import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
)

// A Value is a contract which is implemented by each kind of message Value, like simple, array or plural.
type Value interface {
	ID() string

	// TextArray returns the backing strings array
	TextArray() ([]string, error)

	// Text formats the value with the given arguments
	Text(args ...interface{}) (string, error)

	// QuantityText formats the value with the given arguments
	QuantityText(quantity int, args ...interface{}) (string, error)

	// Locale returns the CLDR language tag
	Locale() string
	goEmitImportValue(group *jen.Group)
	goEmitGetter() *jen.Statement
	exampleText() string

	// implementation detail
	updateTag(tag language.Tag)
}

type PluralBuilder interface {
	Value
	Zero(text string) PluralBuilder
	One(text string) PluralBuilder
	Two(text string) PluralBuilder
	Few(text string) PluralBuilder
	Many(text string) PluralBuilder
	Other(text string) PluralBuilder
}

// A pluralValue is Value with CLDR plural rules, see also
// https://unicode.org/reports/tr35/tr35-numbers.html#Language_Plural_Rules
type pluralValue struct {
	locale string
	Id     string
	zero   string
	one    string
	two    string
	few    string
	many   string
	other  string
	tag    language.Tag
}

func NewQuantityText(locale string, id string) PluralBuilder {
	return pluralValue{
		locale: locale,
		Id:     id,
	}
}

// mapOf returns a map of zero,one,two,few,many,other strings with their according value, if the value is not empty
func (p pluralValue) mapOf() map[string]string {
	m := make(map[string]string)
	if len(p.zero) > 0 {
		m[zero] = p.zero
	}
	if len(p.one) > 0 {
		m[one] = p.zero
	}
	if len(p.two) > 0 {
		m[two] = p.zero
	}
	if len(p.few) > 0 {
		m[few] = p.zero
	}
	if len(p.many) > 0 {
		m[many] = p.zero
	}
	if len(p.other) > 0 {
		m[other] = p.zero
	}
	return m
}

func (p pluralValue) Zero(text string) PluralBuilder {
	p.zero = text
	return p
}

func (p pluralValue) One(text string) PluralBuilder {
	p.one = text
	return p
}

func (p pluralValue) Two(text string) PluralBuilder {
	p.two = text
	return p
}

func (p pluralValue) Few(text string) PluralBuilder {
	p.few = text
	return p
}

func (p pluralValue) Many(text string) PluralBuilder {
	p.many = text
	return p
}

func (p pluralValue) Other(text string) PluralBuilder {
	p.other = text
	return p
}

func (p pluralValue) Locale() string {
	return p.locale
}

func (p pluralValue) updateTag(tag language.Tag) {
	p.tag = tag
}

func (p pluralValue) ID() string {
	return p.Id
}

// StringArray returns the Other Value within a one element array
func (p pluralValue) TextArray() ([]string, error) {
	return []string{p.other}, nil
}

// String returns other
func (p pluralValue) Text(args ...interface{}) (string, error) {
	return fmt.Sprintf(p.other, args...), nil
}

// QuantityString returns the grammatical plural for the internal Plural implementation
func (p pluralValue) QuantityText(quantity int, args ...interface{}) (string, error) {
	// MatchPlural is over engineered: f and t are not used anyway and according to its test, v and w must be kept 0
	// to just get the plural for a natural number
	form := plural.Cardinal.MatchPlural(p.tag, quantity, 0, 0, 0, 0)
	switch form {
	case plural.Zero:
		return p.fallback(p.zero, args...)
	case plural.One:
		return p.fallback(p.one, args...)
	case plural.Two:
		return p.fallback(p.two, args...)
	case plural.Few:
		return p.fallback(p.few, args...)
	case plural.Many:
		return p.fallback(p.many, args...)
	case plural.Other:
		fallthrough
	default:
		return fmt.Sprintf(p.other, args...), nil
	}
}

// fallback uses Other if text is empty
func (p pluralValue) fallback(text string, args ...interface{}) (string, error) {
	if len(text) == 0 {
		return fmt.Sprintf(p.other, args...), nil
	}

	return fmt.Sprintf(text, args...), nil
}

// A simpleValue just holds a text
type simpleValue struct {
	locale string
	Id     string
	String string
}

// NewText returns a
func NewText(locale string, id string, text string) Value {
	return simpleValue{
		locale: locale,
		Id:     id,
		String: text,
	}
}

func (s simpleValue) ID() string {
	return s.Id
}

func (s simpleValue) Locale() string {
	return s.locale
}

func (s simpleValue) updateTag(tag language.Tag) {

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
func (s simpleValue) QuantityText(quantity int, args ...interface{}) (string, error) {
	return s.Text(args...)
}

// An arrayValue holds an ordered bunch of strings
type arrayValue struct {
	locale  string
	Id      string
	Strings []string
}

// NewTextArray creates a new translated array value
func NewTextArray(locale string, id string, items ...string) Value {
	return arrayValue{
		locale:  locale,
		Id:      id,
		Strings: items,
	}
}

func (a arrayValue) updateTag(tag language.Tag) {

}

func (a arrayValue) Locale() string {
	return a.locale
}

func (a arrayValue) ID() string {
	return a.Id
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
func (a arrayValue) QuantityText(quantity int, args ...interface{}) (string, error) {
	return a.Text(args...)
}
