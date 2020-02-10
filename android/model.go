/*
 * Copyright 2020 Torben Schinke
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package android

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type special struct {
	c       byte
	escaped string
}

// nolint: gochecknoglobals
var specials = []special{
	{'@', `\@`},
	{'?', `\?`},
	{'\'', `\'`},
	{'"', `\"`},
}

// Resources is the root element of android resources in general
type Resources struct {
	XMLName      xml.Name      `xml:"resources"`
	Strings      []String      `xml:"string"`
	StringArrays []StringArray `xml:"string-array"`
	Plurals      []Plurals     `xml:"plurals"`
}

// String is the simple text element
type String struct {
	XMLName      xml.Name `xml:"string"`
	Name         string   `xml:"name,attr"`
	Translatable *bool    `xml:"translatable,attr"`
	Text         string   `xml:",chardata"`
}

// StringArray cannot contain placeholders or plurals
type StringArray struct {
	XMLName      xml.Name `xml:"string-array"`
	Name         string   `xml:"name,attr"`
	Translatable *bool    `xml:"translatable,attr"`
	Items        []string `xml:"item"`
}

// Plurals contains the CLDR classified translations for one, other, many etc
type Plurals struct {
	XMLName xml.Name     `xml:"plurals"`
	Name    string       `xml:"name,attr"`
	Items   []PluralItem `xml:"item"`
}

// PluralItem is the quantified message
type PluralItem struct {
	XMLName  xml.Name `xml:"item"`
	Quantity string   `xml:"quantity,attr"`
	Text     string   `xml:",chardata"`
}

// Read parses an android strings.xml document
func Read(reader io.Reader) (Resources, error) {
	res := Resources{}
	tmp, err := ioutil.ReadAll(reader)

	if err != nil {
		return res, fmt.Errorf("failed to read entire xml: %w", err)
	}

	err = xml.Unmarshal(tmp, &res)
	if err != nil {
		return res, fmt.Errorf("failed to parse xml: %w", err)
	}

	return res, nil
}

// ReadFile parses an android strings.xml file from the file system
func ReadFile(fname string) (Resources, error) {
	file, err := os.Open(fname)
	if err != nil {
		return Resources{}, fmt.Errorf("cannot open '%s':%w", fname, err)
	}

	defer func() {
		_ = file.Close()
	}()

	return Read(file)
}

// Decodes unescapes the android string and also replaces the indexed arguments with the notation understood by go
func Decode(androidStr string) string {
	//nolint: gomnd // cannot be escaped with 1 or less chars
	if len(androidStr) <= 1 {
		return androidStr
	}

	// trim left and right, see https://developer.android.com/guide/topics/resources/string-resource.html#escaping_quotes
	if androidStr[0] == '"' && androidStr[len(androidStr)-1] == '"' {
		androidStr = androidStr[1 : len(androidStr)-1]
	}

	// decode special chars, see https://developer.android.com/guide/topics/resources/string-resource.html#escaping_quotes
	for _, special := range specials {
		// this not optimal, a string builder with a custom loop would be way more efficient
		androidStr = strings.ReplaceAll(androidStr, special.escaped, string(special.c))
	}

	// detect %1$s and %2$d types of indices
	// there are a lot of them https://developer.android.com/reference/java/util/Formatter
	//TODO
	return androidStr
}
