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

import (
	"testing"
)

func setup() {
	allResources = newLocalizations()
}

func TestImport(t *testing.T) {
	setup()

	err := ImportFile(AndroidImporter{},  "example/strings_test.xml")
	if err != nil {
		t.Fatal(err)
	}

	res := From("de-123")
	str, err := res.Text("app_name")

	if err != nil {
		t.Fatal(err)
	}

	expected := "EasyApp"

	if str != expected {
		t.Fatalf("expected %s but got %s", expected, str)
	}

	expected = "nick has 1 cat2"
	str, err = res.QuantityText("x_has_y_cats2", 1, "nick", 1)

	if err != nil {
		t.Fatal(err)
	}

	if str != expected {
		t.Fatalf("expected '%s' but got '%s'", expected, str)
	}

	expected = "the owner of 2 cats2 is nick"
	str, err = res.QuantityText("x_has_y_cats2", 2, "nick", 2)

	if err != nil {
		t.Fatal(err)
	}

	if str != expected {
		t.Fatalf("expected '%s' but got '%s'", expected, str)
	}

	err = Validate()
	var errs []error
	if err != nil {
		errs = err.(ErrList).Errs
	}
	if len(errs) != 0 {
		t.Fatal(errs)
	}
}

func TestChecker(t *testing.T) {
	setup()

	err := ImportFile(AndroidImporter{},  "example/strings_test.xml")
	if err != nil {
		t.Fatal(err)
	}

	err = ImportFile(AndroidImporter{},  "example/ignore-strings-de-DE_broken.xml")
	if err != nil {
		t.Fatal(err)
	}

	err = Validate()
	var errs []error
	if err != nil {
		errs = err.(ErrList).Errs
	}
	if len(errs) != 5 {
		for _, err := range errs {
			t.Error(err)
		}
		t.Fatal(len(errs), "fails")
	}
}
