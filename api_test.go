package i18n

import (
	"testing"
)

func TestImport(t *testing.T) {
	err := ImportFile(AndroidImporter{}, "en-US", "android/strings_test.xml")
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
}
