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
	"fmt"
	"testing"
)

func TestTranslate(t *testing.T) {
	res, err := ReadFile("strings_test.xml")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", res)
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{"special chars escape", `\@ \? < & ' " \" \'`, `@ ? < & ' " " '`},
		{"special chars full escape", `"hello '"`, `hello '`},
		{"conversion with indices", `hello %%1$s %s %2$d %13$s`, `hello %%1$s %s %[2]d %[13]s`},
	}
	// nolint: scopelint // tt is a value, so this is a false-positive
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Decode(tt.args); got != tt.want {
				t.Errorf("Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}
