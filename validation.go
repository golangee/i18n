package i18n

import (
	"fmt"
	"reflect"
	"strings"
)

// ErrMissingValue contains an example value and the resources which is missing it
type ErrMissingValue struct {
	Value Value
}

func (e ErrMissingValue) Error() string {
	return e.Value.Locale() + " is missing" + e.Value.ID()
}

// ErrTypeMismatch contains two Values of two different values which have different types, which is not allowed.
type ErrTypeMismatch struct {
	Value0 Value
	Value1 Value
}

func (e ErrTypeMismatch) Error() string {
	return e.Value0.ID() + " is a " + reflect.TypeOf(e.Value0).String() + " in " + e.Value0.Locale() +
		" but in " + e.Value1.Locale() + " a " + reflect.TypeOf(e.Value0).String()
}

// ErrFormatSpecifierCountMismatch is used to indicates that the number printf formatting directives is different
// but they must be equal.
type ErrFormatSpecifierCountMismatch struct {
	Value0 Value
	Specs0 []PrintfFormatSpecifier
	Value1 Value
	Specs1 []PrintfFormatSpecifier
}

func (e *ErrFormatSpecifierCountMismatch) Error() string {
	return fmt.Sprintf("printf argument count mismatch: %s.%s has %d arguments but %s.%s has %d",
		e.Value0.Locale(), e.Value0.ID(), len(e.Specs0), e.Value1.Locale(), e.Value1.ID(), len(e.Specs1))
}

// ErrArrayCountMismatch indicates that two arrays must have the same amount of entries
type ErrArrayCountMismatch struct {
	Value0 Value
	Count0 int
	Value1 Value
	Count1 int
}

func (e ErrArrayCountMismatch) Error() string {
	return fmt.Sprintf("The array count in %s.%s is %d but %s.%s has %d",
		e.Value0.Locale(), e.Value0.ID(), e.Count0, e.Value1.Locale(), e.Value1.ID(), e.Count1)
}

// ErrUnexpectedAmountOfFormatSpecifiers indicates that a value has an unexpected amount of specifiers.
// E.g. arrays must not contain any specifiers.
type ErrUnexpectedAmountOfFormatSpecifiers struct {
	Value    Value
	Found    int
	Expected int
	Text     string
}

func (e *ErrUnexpectedAmountOfFormatSpecifiers) Error() string {
	return fmt.Sprintf("The value %s.%s has %d format specifiers but expected are %d (%s)",
		e.Value.Locale(), e.Value.ID(), e.Found, e.Expected, e.Text)
}

// ErrVerbConflict is returned, if two strings have different verb specifiers for the same position
type ErrVerbConflict struct {
	Value0 Value
	Verb0  PrintfFormatSpecifier
	Value1 Value
	Verb1  PrintfFormatSpecifier
}

func (e *ErrVerbConflict) Error() string {
	return fmt.Sprintf("The value %s.%s has at index %d the verb '%s' but"+
		" %s.%s has the verb '%s'",
		e.Value0.Locale(), e.Value0.ID(), e.Verb0.Index, string(e.Verb0.Verb()),
		e.Value1.Locale(), e.Value1.ID(), string(e.Verb1.Verb()))
}

// ErrOtherMissing indicates a missing "other" value for a plural. You may omit everything else but
// other is the fallback and must not be empty at least.
type ErrOtherMissing struct {
	Value Value
}

func (e ErrOtherMissing) Error() string {
	return "the plural 'other' must not be empty of " + e.Value.Locale() + "." + e.Value.ID()
}

// ErrList is a list of errors
type ErrList struct {
	Errs []error
}

func (e ErrList) Error() string {
	sb := &strings.Builder{}
	for _, err := range e.Errs {
		sb.WriteString(err.Error())
		sb.WriteByte('\n')
	}
	return sb.String()
}

// validate checks the consistency of the given resources. The following checks are made
//  * each resources have the same keys
//  * each resources have the same type
//  * the order and type of verbs are equal
func validate(resources []*Resources) error {
	var errs []error
	for i0, r0 := range resources {
		for i1 := i0 + 1; i1 < len(resources); i1++ {
			r1 := resources[i1]
			if r0 == r1 {
				continue
			}
			r0.mutex.RLock()
			r1.mutex.RLock()

			for k0, v0 := range r0.values {
				v1, exists := r1.values[k0]
				if !exists {
					errs = append(errs, ErrMissingValue{
						Value: v0,
					})
				} else {
					if reflect.TypeOf(v0) != reflect.TypeOf(v1) {
						errs = append(errs, ErrTypeMismatch{
							Value0: v0,
							Value1: v1,
						})
					} else {
						var strErr error
						switch t0 := (v0).(type) {
						case simpleValue:
							t1 := v1.(simpleValue)
							strErr = validatePrintf(t0.String, t1.String, -1)
						case pluralValue:
							t1 := v1.(pluralValue)
							if len(t0.other) == 0 {
								strErr = ErrOtherMissing{Value: v0}
								break
							}
							if len(t1.other) == 0 {
								strErr = ErrOtherMissing{Value: v1}
								break
							}

						anyPlural:
							for _, p0 := range t0.mapOf() {
								for _, p1 := range t1.mapOf() {
									strErr = validatePrintf(p0, p1, -1)
									if strErr != nil {
										break anyPlural
									}
								}
							}

						case arrayValue:
							t1 := v1.(arrayValue)
							if len(t0.Strings) != len(t1.Strings) {
								errs = append(errs, ErrArrayCountMismatch{
									Value0: v0,
									Count0: len(t0.Strings),
									Value1: v1,
									Count1: len(t1.Strings),
								})
							} else {
								for i := range t0.Strings {
									strErr = validatePrintf(t0.Strings[i], t1.Strings[i], 0)
									if strErr != nil {
										break
									}
								}
							}

						}
						// we enrich the errors here instead of wrapping over, which is unnecessary complex
						if strErr != nil {
							switch e := (strErr).(type) {
							case *ErrFormatSpecifierCountMismatch:
								e.Value0 = v0
								e.Value1 = v1
							case *ErrUnexpectedAmountOfFormatSpecifiers:
								e.Value = v0
							case *ErrVerbConflict:
								e.Value0 = v0
								e.Value1 = v1
							}
							errs = append(errs, strErr)
						}
					}

				}
			}

			r0.mutex.RUnlock()
			r1.mutex.RUnlock()
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return ErrList{errs}
}

// validatePrintf validates str0 and str1 to be of equal golang printf format directives. If expected is not -1
// an error is returned, if the amount of directives does not match the expected number.
func validatePrintf(str0, str1 string, expected int) error {
	specs0 := ParsePrintf(str0)
	specs1 := ParsePrintf(str1)

	if len(specs0) != len(specs1) {
		return &ErrFormatSpecifierCountMismatch{
			Specs0: specs0,
			Specs1: specs1,
		}
	}

	if expected >= 0 && len(specs0) != expected {
		return &ErrUnexpectedAmountOfFormatSpecifiers{
			Found:    len(specs0),
			Expected: expected,
			Text:     str0,
		}
	}

	for i := range specs0 {
		spec0 := specs0[i]
		spec1 := specs1[i]

		if spec0.Verb() != spec1.Verb() {
			return &ErrVerbConflict{
				Verb0: spec0,
				Verb1: spec1,
			}
		}
	}
	return nil
}
