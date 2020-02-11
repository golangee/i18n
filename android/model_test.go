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
		{"conversion with indices", `hello %%1$s %s %2$d %3$s`, `hello %%1$s %s %[2]d %[3]s`},
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
