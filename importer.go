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

// nolint: goimports // the linter is broken
import (
	"fmt"
	"github.com/golangee/i18n/android"
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

	locale := dst.tag.String()
	for _, str := range src.Strings {
		dst.values[str.Name] = simpleValue{
			Id:     str.Name,
			locale: locale,
			String: android.Decode(str.Text),
		}
	}

	for _, pl := range src.Plurals {
		val := pluralValue{
			Id:     pl.Name,
			tag:    dst.tag,
			locale: locale,
		}

		for _, item := range pl.Items {
			switch strings.ToLower(item.Quantity) {
			case zero:
				val.zero = android.Decode(item.Text)
			case one:
				val.one = android.Decode(item.Text)
			case two:
				val.two = android.Decode(item.Text)
			case few:
				val.few = android.Decode(item.Text)
			case many:
				val.many = android.Decode(item.Text)
			case other:
				fallthrough
			default:
				val.other = android.Decode(item.Text)
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
			Id:      arr.Name,
			locale:  locale,
			Strings: tmp,
		}
	}
}
