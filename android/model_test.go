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

	bad0 := `\@ \? < & ' " \" \'`
	good0 := `@ ? < & ' " " '`

	if Decode(bad0) != good0 {
		t.Fatalf("expected %s but got %s", good0, Decode(bad0))
	}

	bad1 := `"hello '"`
	good1 := `hello '`

	if Decode(bad1) != good1 {
		t.Fatalf("expected %s but got %s", good1, Decode(bad1))
	}
}
