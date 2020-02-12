package i18n

import (
	"testing"
)

func Test_goGenerator_Scan(t *testing.T) {
	gen := NewGoGenerator(".")
	err := gen.Scan()
	if err != nil {
		t.Fatal(err)
	}

	if len(gen.translations) != 1 {
		t.Fatal("expected 1 translation")
	}

	gen.Emit()
}
