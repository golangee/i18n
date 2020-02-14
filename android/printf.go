package android

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// matcher is a monster regex to find android supported printf declarations, in contrast to go indices are %<num>$<verb>
var formatMatcher = regexp.MustCompile(`(?:\x25\x25)|(\x25(?:(?:[1-9]\d*)\$|\((?:[^\)]+)\))?(?:\+)?(?:0|'[^$])?(?:-)?(?:\d+)?(?:\.(?:\d+))?(?:[b-fiosuxX]))`)

// PrintFormatSpecifier represents a part of a printf string, like %s, %[5]d or %10.10s.
type PrintfFormatSpecifier struct {
	Src      string // Src is the entire string
	Pos      int    // Pos is the index in src where this specifier is located
	End      int    // End is the end index in src where this specifier is located
	Index    int    // Index is the argument position. Either this is the natural order or derived by the indexed position.
	PosIndex int    // PosIndex is the positional index, which is the order as parsed
}

// String returns the entire format specifier
func (f *PrintfFormatSpecifier) String() string {
	return f.Src[f.Pos:f.End]
}

// Verb returns the single character and not the entire formatting directive.
func (f *PrintfFormatSpecifier) Verb() byte {
	return f.Src[f.End-1]
}

// Indexed returns true, if a %<num>$<verb> structure is declared
func (f *PrintfFormatSpecifier) Indexed() bool {
	return strings.Index(f.String(), "$") > 0
}

// AsGoFormatSpecifier converts the $ index into a [] index
func (f *PrintfFormatSpecifier) AsGoFormatSpecifier() string {
	if f.Indexed() {
		str := f.String()
		b := strings.Index(str, "%")
		e := strings.Index(str, "$")
		return "%[" + str[b+1:e] + "]" + str[e+1:]
	}
	return f.String()
}

// index parses tries to parse the from String(). Returns -1 if not
func (f *PrintfFormatSpecifier) index() int {
	str := f.String()
	b := strings.Index(str, "%")
	e := strings.Index(str, "$")
	if e > b {
		i, err := strconv.ParseInt(str[b+1:e], 10, 32)
		if err != nil {
			return -1
		}
		return int(i)
	}
	return -1
}

// ParsePrintf returns all found format specifiers and returns them in a sorted order by index position
func ParsePrintf(str string) []PrintfFormatSpecifier {
	var specs []PrintfFormatSpecifier
	indices := formatMatcher.FindAllStringIndex(str, -1)
	idx := 0
	for _, pos := range indices {
		spec := PrintfFormatSpecifier{
			Src:      str,
			Pos:      pos[0],
			End:      pos[1],
			Index:    idx,
			PosIndex: idx,
		}
		overloadedIndex := spec.index()
		if overloadedIndex > -1 {
			spec.Index = overloadedIndex
		}
		if spec.String() == "%%" {
			continue
		}
		specs = append(specs, spec)
		idx++
	}
	sort.Sort(pfsSortByIndex(specs))
	// enumerate cleanly
	for i := range specs {
		specs[i].Index = i
	}
	return specs
}

type pfsSortByIndex []PrintfFormatSpecifier

func (p pfsSortByIndex) Len() int {
	return len(p)
}

func (p pfsSortByIndex) Less(i, j int) bool {
	return p[i].Index < p[j].Index
}

func (p pfsSortByIndex) Swap(i, j int) {
	p[j], p[i] = p[i], p[j]
}

type pfsSortByPosIndex []PrintfFormatSpecifier

func (p pfsSortByPosIndex) Len() int {
	return len(p)
}

func (p pfsSortByPosIndex) Less(i, j int) bool {
	return p[i].PosIndex < p[j].PosIndex
}

func (p pfsSortByPosIndex) Swap(i, j int) {
	p[j], p[i] = p[i], p[j]
}
