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

package android

import (
	"reflect"
	"testing"
)

func TestParsePrintfIndexOrder(t *testing.T) {
	elems := ParsePrintf("%3$s %21$b %1$d")
	if elems[0].Verb() != 'd' || elems[0].Index != 0 {
		t.Fatalf("wrong position '%s' -> %s", string(elems[0].Verb()), elems[0].String())
	}

	if elems[1].Verb() != 's' || elems[1].Index != 1 {
		t.Fatal("unexpected position", elems[1].String())
	}

	if elems[2].Verb() != 'b' || elems[2].Index != 2 {
		t.Fatal("unexpected position", elems[2].String())
	}

	if elems[0].AsGoFormatSpecifier() != "%[1]d" {
		t.Fatalf("unexpected %s", elems[0].AsGoFormatSpecifier())
	}

	if elems[1].AsGoFormatSpecifier() != "%[3]s" {
		t.Fatalf("unexpected %s", elems[1].AsGoFormatSpecifier())
	}

	if elems[2].AsGoFormatSpecifier() != "%[21]b" {
		t.Fatalf("unexpected %s", elems[2].AsGoFormatSpecifier())
	}
}

func TestParsePrintf(t *testing.T) {
	tests := []struct {
		name string
		args string // input value
		want string // format specifier
	}{
		{"binary representation", "%%b = '%b'", "%b"},
		{"print the ascii character, same as chr() function", "%%c = '%c'", "%c"},
		{"standard integer representation", "%%d = '%d'", "%d"},
		{"scientific notation", "%%e = '%e'", "%e"},
		{"integer representation", "%%u = '%u'", "%u"},
		{"floating point representation", "%%f = '%f'", "%f"},
		{"octal representation", "%%o = '%o'", "%o"},
		{"string representation", "%%s = '%s'", "%s"},
		{"hexadecimal representation lower-case", "%%x = '%x'", "%x"},
		{"hexadecimal representation upper-case", "%%X = '%X'", "%X"},
		{"sign specifier on an integer", "%%+d = '%+d'", "%+d"},
		{"standard string output", "[%s]", "%s"},
		{"standard string output with android index", "[%1$s]", "%1$s"},
		{"standard string output with android index 2", "[%21$s]", "%21$s"},
		{"right-justification with spaces", "[%10s]", "%10s"},
		{"left-justification with spaces", "[%-10s]", "%-10s"},
		{"zero-padding works on strings too", "[%010s]", "%010s"},
		{"use the custom padding character '#'", "[%'#10s]", "%'#10s"},
		{"left-justification but with a cutoff of 10 characters", "[%10.10s]", "%10.10s"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParsePrintf(tt.args)[0].String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParsePrintf() = %v, want %v", got, tt.want)
			}
		})
	}
}
